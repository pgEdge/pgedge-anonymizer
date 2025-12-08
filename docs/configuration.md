# Configuration Reference

By default, pgEdge Anonymizer looks for a configuration file named `pgedge-anonymizer.yaml` in the current directory. Include the `--config` option to specify an alternative path to the configuration file when invoking `pgedge-anonymizer`:

```bash
pgedge-anonymizer run --config /path/to/config.yaml
```

Anonymizer also supports the use of standard [PostgreSQL environment variables](#using-environment-variables-for-options); options specified on the command line and in the configuration file take precedence over environment variable values.


## Using Command-Line Options

When invoking `pgedge-anonymizer`, you can specify database settings with command-line flags or in a configuration file; command-line options will override any preferences specified in the configuration file or as an environment variable:

```bash
pgedge-anonymizer run \
  --config config.yaml \
  --host production-db.example.com \
  --port 5433 \
  --database myapp_staging \
  --user admin \
  --password secret \
  --sslmode require \
  --quiet
```

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Configuration file path |
| `--quiet` | `-q` | Suppress progress output |
| `--host` | | Database host |
| `--port` | | Database port |
| `--database` | `-d` | Database name |
| `--user` | `-U` | Database user |
| `--password` | | Database password |
| `--sslmode` | | SSL mode |


## Specifying Properties in a Configuration File

The configuration file contains content in three major sections:

* Database Properties
* Pattern Properties
* Column Properties

### Database Properties

Include a `database` properties section in your configuration file to specify PostgreSQL connection settings:

```yaml
database:
  host: localhost
  port: 5432
  database: myapp
  user: anonymizer
  password: ""
  sslmode: prefer
```

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `host` | string | localhost | Database server hostname or IP |
| `port` | integer | 5432 | Database server port |
| `database` | string | (required) | Database name to connect to |
| `user` | string | (required) | Database user for authentication |
| `password` | string | "" | Database password (prompts if empty) |
| `sslmode` | string | prefer | SSL connection mode; specify `disable`, `prefer` (use if available), `require` (require SSL, but don't verify the certificate), `verify-ca` (require SSL, and verify the CA signature), or `verify-full` (require SSL and verify the CA and hostname). |

If you specify an SSL mode that requires certificates, include supporting properties in the `database` section:

```yaml
database:
  sslmode: verify-full
  sslcert: /path/to/client-cert.pem
  sslkey: /path/to/client-key.pem
  sslrootcert: /path/to/ca-cert.pem
```

Where:
    *  `sslcert` is the path to the client certificate file.
    * `sslkey` is the path to the client private key file.
    * `sslrootcert` is the path to the CA certificate file.


### Pattern Options

Patterns specify the form that replacement content will take when anonymizing your columns.  Patterns can be either user-defined, or a pre-defined pattern.  Patterns are stored in a .yaml file identified in the configuration file by the following properties:

```yaml
patterns:
  default_path: /etc/pgedge/pgedge-anonymizer-patterns.yaml
  user_path: ./my-patterns.yaml
  disable_defaults: false
```

Pattern properties specify:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `default_path` | string | (auto) | Path to the [default patterns file](patterns.md). |
| `user_path` | string | "" | Path to [user-defined patterns](#creating-user-defined-patterns). |
| `disable_defaults` | boolean | false | Skip loading built-in patterns. |

If a `default_path` is not specified, the tool searches for `pgedge-anonymizer-patterns.yaml` in the following locations:

1. The current working directory.
2. `/etc/pgedge/`
3. `~/.config/pgedge/`

#### Creating User-Defined Patterns

You can create custom patterns for your data in a separate .yaml file; for example:

```yaml
# my-patterns.yaml
patterns:
  - name: EMPLOYEE_ID
    replacement: "EMP-XXXXXX"
    note: "Internal employee identifiers"

  - name: CUSTOMER_REF
    replacement: "CUST-XXXXXXXX"
    note: "Customer reference numbers"
```

Then, reference your user-defined pattern file in the configuration file with the following properties:

```yaml
patterns:
  user_path: ./my-patterns.yaml
```

!!! note
    User-defined pattern names must not conflict with built-in patterns
    unless `disable_defaults: true` is set.


## Columns Section

Use the configuration file to specify the columns to anonymize with fully-qualified names that include the `schema_name`, `table_name`, and `column_name` information, and the pattern_name that will apply to the data stored in that column:

```yaml
columns:
  - column: schema_name.table_name.column_name
    pattern: pattern_name
```

Where:

* `column` is a fully-qualified string that specifies the schema_name, table_name, and column name of the column you are anonymizing.
* `pattern` specifies the name of the pattern that you wish to apply to the column.

For example:

```yaml
columns:
  # User personal data
  - column: public.users.first_name
    pattern: PERSON_FIRST_NAME

  - column: public.users.last_name
    pattern: PERSON_LAST_NAME

  - column: public.users.email
    pattern: EMAIL

  - column: public.users.phone
    pattern: US_PHONE

  - column: public.users.date_of_birth
    pattern: DOB_OVER_18

  - column: public.users.ssn
    pattern: US_SSN

  # Address data
  - column: public.addresses.street_address
    pattern: ADDRESS

  # Payment data
  - column: public.payments.card_number
    pattern: CREDIT_CARD

  - column: public.payments.card_expiry
    pattern: CREDIT_CARD_EXPIRY

  - column: public.payments.card_cvv
    pattern: CREDIT_CARD_CVV
```

## Using Environment Variables for Options

Database options can also be set via standard PostgreSQL environment variables.  If pgEdge Anonymizer does not locate information on the command line or in the configuration file, it will then check the values specified in the following environment variables:

| Config Option | Environment Variable |
|---------------|---------------------|
| `host` | `PGHOST` |
| `port` | `PGPORT` |
| `database` | `PGDATABASE` |
| `user` | `PGUSER` |
| `password` | `PGPASSWORD` |
| `sslmode` | `PGSSLMODE` |


