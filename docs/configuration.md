# Configuration Reference

By default, pgEdge Anonymizer looks for a configuration file named `pgedge-anonymizer.yaml` in the current directory. When invoking `pgedge-anonymizer`, include the `--config` option to specify an alternative path to the configuration file:

```bash
pgedge-anonymizer run --config /path/to/config.yaml
```

The configuration file is organized in three major sections:

* [Database Properties](#specifying-properties-in-the-database-section)
* [Pattern Properties](#specifying-properties-in-the-pattern-section)
* [Column Properties](#specifying-properties-in-the-columns-section)

When invoking `pgedge-anonymizer`, you can specify database connection settings with command-line flags or in a configuration file; command-line options for database settings will override values set elsewhere:

```bash
pgedge-anonymizer run \
  --config config.yaml \              # Uses config.yaml instead of the default config file
  --host production-db.example.com \  # Overrides database host from config file
  --port 5433 \                       # Overrides database port from config file
  --database myapp_staging \          # Overrides database name from config file
  --user admin \                      # Overrides database user from config file
  --password secret \                 # Overrides database password from config file
  --sslmode require \                 # Overrides SSL mode from config file
  --quiet                             # Suppresses progress output
```

!!! hint

    Anonymizer also supports the use of standard PostgreSQL environment variables for database connection options; options specified on the command line and in the configuration file take precedence over environment variable values.


## Specifying Properties in the Database Section

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
    * `sslcert` is the path to the client certificate file.
    * `sslkey` is the path to the client private key file.
    * `sslrootcert` is the path to the CA certificate file.

**Using Environment Variables for Options**

Database options can also be set via standard PostgreSQL environment variables.  If pgEdge Anonymizer does not locate database connection information on the command line or in the configuration file, it will then check the values specified in the following environment variables:

| Config Option | Environment Variable |
|---------------|---------------------|
| `host` | `PGHOST` |
| `port` | `PGPORT` |
| `database` | `PGDATABASE` |
| `user` | `PGUSER` |
| `password` | `PGPASSWORD` |
| `sslmode` | `PGSSLMODE` |


## Specifying Properties in the Pattern Section

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

You can also [create custom patterns](custom_pattern.md) for your data in a separate .yaml file; for example:

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
    User-defined pattern names must not conflict with built-in patterns unless `disable_defaults: true` is set.


## Specifying Properties in the Columns Section

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
