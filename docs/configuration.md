# Configuration Reference

This page provides complete documentation for pgEdge Anonymizer
configuration options.

## Configuration File

By default, pgEdge Anonymizer looks for `pgedge-anonymizer.yaml` in the
current directory. Use `--config` to specify an alternative path:

```bash
pgedge-anonymizer run --config /path/to/config.yaml
```

## Database Section

Configure PostgreSQL connection settings:

```yaml
database:
  host: localhost
  port: 5432
  database: myapp
  user: anonymizer
  password: ""
  sslmode: prefer
```

### Connection Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `host` | string | localhost | Database server hostname or IP |
| `port` | integer | 5432 | Database server port |
| `database` | string | (required) | Database name to connect to |
| `user` | string | (required) | Database user for authentication |
| `password` | string | "" | Database password (prompts if empty) |
| `sslmode` | string | prefer | SSL connection mode |

### SSL Modes

| Mode | Description |
|------|-------------|
| `disable` | No SSL |
| `prefer` | Use SSL if available (default) |
| `require` | Require SSL, don't verify certificate |
| `verify-ca` | Require SSL, verify CA signature |
| `verify-full` | Require SSL, verify CA and hostname |

### SSL Certificate Options

For SSL modes that require certificates:

```yaml
database:
  sslmode: verify-full
  sslcert: /path/to/client-cert.pem
  sslkey: /path/to/client-key.pem
  sslrootcert: /path/to/ca-cert.pem
```

| Option | Description |
|--------|-------------|
| `sslcert` | Path to client certificate file |
| `sslkey` | Path to client private key file |
| `sslrootcert` | Path to CA certificate file |

### Environment Variables

Database options can also be set via standard PostgreSQL environment
variables:

| Config Option | Environment Variable |
|---------------|---------------------|
| `host` | `PGHOST` |
| `port` | `PGPORT` |
| `database` | `PGDATABASE` |
| `user` | `PGUSER` |
| `password` | `PGPASSWORD` |
| `sslmode` | `PGSSLMODE` |

Priority order (highest to lowest):

1. Command-line flags
2. Configuration file values
3. Environment variables

## Patterns Section

Configure pattern file loading:

```yaml
patterns:
  default_path: /etc/pgedge/pgedge-anonymizer-patterns.yaml
  user_path: ./my-patterns.yaml
  disable_defaults: false
```

### Pattern Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `default_path` | string | (auto) | Path to default patterns file |
| `user_path` | string | "" | Path to user-defined patterns |
| `disable_defaults` | boolean | false | Skip loading built-in patterns |

### Default Pattern Search Locations

If `default_path` is not specified, the tool searches for
`pgedge-anonymizer-patterns.yaml` in:

1. Current working directory
2. `/etc/pgedge/`
3. `~/.config/pgedge/`

### User-Defined Patterns

Create custom patterns in a separate YAML file:

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

Reference in configuration:

```yaml
patterns:
  user_path: ./my-patterns.yaml
```

!!! note
    User-defined pattern names must not conflict with built-in patterns
    unless `disable_defaults: true` is set.

## Columns Section

Specify columns to anonymize:

```yaml
columns:
  - column: schema.table.column
    pattern: PATTERN_NAME
```

### Column Configuration

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `column` | string | Yes | Fully-qualified column name |
| `pattern` | string | Yes | Pattern name to apply |

### Column Name Format

Columns must be specified as `schema.table.column`:

- **schema**: PostgreSQL schema (usually `public`)
- **table**: Table name
- **column**: Column name

Examples:

```yaml
columns:
  # Public schema (most common)
  - column: public.users.email
    pattern: EMAIL

  # Custom schema
  - column: sales.customers.phone
    pattern: US_PHONE

  # Schema with special characters (use quotes)
  - column: "my-schema.my-table.my-column"
    pattern: PERSON_NAME
```

### Multiple Columns Example

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

## Complete Example

A full configuration file combining all sections:

```yaml
# pgEdge Anonymizer Configuration

# Database connection
database:
  host: db.example.com
  port: 5432
  database: production_copy
  user: anonymizer
  password: ""  # Will prompt or use PGPASSWORD
  sslmode: require

# Pattern configuration
patterns:
  user_path: ./custom-patterns.yaml
  disable_defaults: false

# Columns to anonymize
columns:
  # Customer PII
  - column: public.customers.name
    pattern: PERSON_NAME

  - column: public.customers.email
    pattern: EMAIL

  - column: public.customers.phone
    pattern: US_PHONE

  - column: public.customers.dob
    pattern: DOB_OVER_18

  # Employee data
  - column: hr.employees.ssn
    pattern: US_SSN

  - column: hr.employees.passport_number
    pattern: PASSPORT

  # UK-specific data
  - column: uk.customers.ni_number
    pattern: UK_NI

  - column: uk.customers.nhs_number
    pattern: UK_NHS

  # Notes and free text
  - column: public.support_tickets.description
    pattern: LOREMIPSUM
```

## Command-Line Overrides

All database settings can be overridden via command-line flags:

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

## Validation

Always validate your configuration before running:

```bash
pgedge-anonymizer validate --config config.yaml
```

This verifies:

- YAML syntax is valid
- Required fields are present
- Database connection succeeds
- All specified columns exist
- All pattern names are recognized
