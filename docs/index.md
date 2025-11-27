# pgEdge Anonymizer

pgEdge Anonymizer is a command-line tool for anonymizing personally
identifiable information (PII) in PostgreSQL databases. It replaces
sensitive data with realistic but fake values while maintaining data
consistency and referential integrity.

## Features

- **Pattern-based anonymization**: 100+ built-in patterns for common PII types
- **Consistent replacement**: Same input values produce the same anonymized
  output within a run
- **Foreign key awareness**: Automatically handles CASCADE relationships
- **Large database support**: Efficient batch processing with server-side
  cursors
- **Format preservation**: Maintains original data formatting where possible
- **Single transaction**: All changes committed atomically or rolled back
- **Extensible**: Define custom patterns for your specific needs

## Quick Start

### Installation

Build from source:

```bash
git clone https://github.com/pgEdge/pgedge-anonymizer.git
cd pgedge-anonymizer
make build
```

### Basic Usage

1. Create a configuration file `pgedge-anonymizer.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  database: myapp
  user: anonymizer

columns:
  - column: public.users.email
    pattern: EMAIL

  - column: public.users.phone
    pattern: US_PHONE

  - column: public.users.ssn
    pattern: US_SSN
```

2. Run the anonymizer:

```bash
pgedge-anonymizer run
```

3. View the results:

```
Processing public.users.email (est. 50000 rows)...
  10000 rows processed
  20000 rows processed
  30000 rows processed
  40000 rows processed
  50000 rows processed
  Completed: 50000 rows, 48234 values anonymized

=== Anonymization Statistics ===
Total columns processed: 1
Total rows processed:    50000
Total values anonymized: 48234
Total duration:          2.34s
Throughput:              21367 rows/sec
```

## Configuration

### Database Connection

Configure database connection in your YAML file or via environment variables:

```yaml
database:
  host: localhost       # or PGHOST
  port: 5432            # or PGPORT
  database: myapp       # or PGDATABASE
  user: anonymizer      # or PGUSER
  password: ""          # or PGPASSWORD (leave empty to prompt)
  sslmode: prefer       # or PGSSLMODE
```

Command-line flags override configuration file settings:

```bash
pgedge-anonymizer run --host myserver --database production --user admin
```

### Specifying Columns

List columns to anonymize using fully-qualified names:

```yaml
columns:
  - column: schema.table.column
    pattern: PATTERN_NAME
```

Example:

```yaml
columns:
  - column: public.users.first_name
    pattern: PERSON_FIRST_NAME

  - column: public.users.last_name
    pattern: PERSON_LAST_NAME

  - column: public.users.email
    pattern: EMAIL

  - column: public.customers.phone_number
    pattern: US_PHONE
```

## Built-in Patterns

pgEdge Anonymizer includes patterns for common PII types:

### Personal Information

| Pattern | Description | Example Output |
|---------|-------------|----------------|
| `PERSON_NAME` | Full name | John Smith |
| `PERSON_FIRST_NAME` | First name only | Jennifer |
| `PERSON_LAST_NAME` | Last name only | Williams |
| `DOB` | Date of birth (any age) | 1985-03-15 |
| `DOB_OVER_13` | DOB for 13+ years old | 2008-07-22 |
| `DOB_OVER_16` | DOB for 16+ years old | 2005-11-08 |
| `DOB_OVER_18` | DOB for 18+ years old | 2003-04-30 |
| `DOB_OVER_21` | DOB for 21+ years old | 2000-09-12 |

### Contact Information

| Pattern | Description | Example Output |
|---------|-------------|----------------|
| `EMAIL` | Email addresses | john.smith@example.com |
| `US_PHONE` | US phone numbers | 555-234-5678 |
| `UK_PHONE` | UK phone numbers | +44 20 7946 0958 |
| `INTERNATIONAL_PHONE` | International format | +33 1234 5678901 |
| `WORLDWIDE_PHONE` | Generic phone | 5552345678 |
| `ADDRESS` | Street addresses | 742 Oak Avenue, Springfield, CA 90210 |

### Financial Information

| Pattern | Description | Example Output |
|---------|-------------|----------------|
| `CREDIT_CARD` | Credit card numbers | 4532-1234-5678-9012 |
| `CREDIT_CARD_EXPIRY` | Expiry dates | 08/27 |
| `CREDIT_CARD_CVV` | CVV codes | 847 |

### Government IDs

| Pattern | Description | Example Output |
|---------|-------------|----------------|
| `US_SSN` | US Social Security | 234-56-7890 |
| `UK_NI` | UK National Insurance | AB123456C |
| `UK_NHS` | UK NHS numbers | 485 777 3456 |
| `PASSPORT` | Passport numbers | A12345678 |

### Text

| Pattern | Description | Example Output |
|---------|-------------|----------------|
| `LOREMIPSUM` | Lorem ipsum text | Lorem ipsum dolor sit... |

## Validation

Before running anonymization, validate your configuration:

```bash
pgedge-anonymizer validate
```

This checks:

- Configuration file syntax
- Database connectivity
- Column existence in the database
- Pattern validity

## Foreign Key Handling

pgEdge Anonymizer automatically analyzes foreign key relationships:

- **CASCADE updates**: If a column has referencing foreign keys with
  `ON UPDATE CASCADE`, the tool updates the source column and PostgreSQL
  propagates changes automatically.

- **Processing order**: Columns are processed in dependency order to
  maintain referential integrity.

- **Skip targets**: Columns that are CASCADE targets of other configured
  columns are automatically skipped to avoid duplicate processing.

## Performance

For large databases, pgEdge Anonymizer uses:

- **Server-side cursors**: Rows are fetched in configurable batches
  (default 10,000)
- **Batch updates**: Multiple rows updated in single statements using
  CTID-based unnest operations
- **Tiered caching**: LRU in-memory cache with SQLite spillover for
  value dictionaries

## Command Reference

### run

Execute anonymization:

```bash
pgedge-anonymizer run [flags]
```

Flags:

- `--config, -c`: Configuration file path (default: pgedge-anonymizer.yaml)
- `--quiet, -q`: Suppress progress output
- `--host`: Database host (overrides config)
- `--port`: Database port (overrides config)
- `--database, -d`: Database name (overrides config)
- `--user, -U`: Database user (overrides config)
- `--password`: Database password (overrides config)
- `--sslmode`: SSL mode (overrides config)

### validate

Validate configuration without running:

```bash
pgedge-anonymizer validate [flags]
```

Flags:

- `--config, -c`: Configuration file path
- Database connection flags (same as run)

### version

Display version information:

```bash
pgedge-anonymizer version
```

## Environment Variables

Standard PostgreSQL environment variables are supported:

- `PGHOST` - Database host
- `PGPORT` - Database port
- `PGDATABASE` - Database name
- `PGUSER` - Database user
- `PGPASSWORD` - Database password
- `PGSSLMODE` - SSL mode

Priority order (highest to lowest):

1. Command-line flags
2. Configuration file
3. Environment variables

## Best Practices

### Before Anonymizing

1. **Back up your database** - Anonymization is irreversible
2. **Test on a copy** - Validate on a non-production database first
3. **Review columns** - Ensure all PII columns are included
4. **Check foreign keys** - Understand CASCADE relationships

### Security Considerations

- Run with minimal required privileges
- Use SSL for database connections in production
- Secure your configuration file (contains connection details)
- Consider running on the database server to avoid network transfer of
  sensitive data

### Performance Tips

- Process during low-traffic periods
- Ensure adequate disk space for transaction logs
- Consider table-level locks for very large updates
- Monitor PostgreSQL logs for any issues

## Troubleshooting

### Common Issues

**"column not found in database"**

- Verify the schema.table.column path is correct
- Check that the user has SELECT permission on the table

**"unknown pattern"**

- Ensure the pattern name matches exactly (case-sensitive)
- Check that custom patterns are properly defined

**"failed to connect"**

- Verify database credentials
- Check network connectivity
- Ensure PostgreSQL is accepting connections

### Getting Help

- [GitHub Issues](https://github.com/pgEdge/pgedge-anonymizer/issues)
- [Documentation](https://pgedge.github.io/pgedge-anonymizer/)
