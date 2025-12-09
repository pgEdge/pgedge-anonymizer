# pgEdge Anonymizer

[![CI](https://github.com/pgEdge/pgedge-anonymizer/actions/workflows/ci.yml/badge.svg)](https://github.com/pgEdge/pgedge-anonymizer/actions/workflows/ci.yml)

- [Introduction](docs/index.md)
      - [Best Practices](docs/best_practices.md)
  - Installing pgEdge Document Loader
      - [Installing Document Loader](docs/installation.md)
      - [Creating a Configuration File](docs/configuration.md)
      - [pgEdge Anonymizer Quickstart](docs/quickstart.md)
  - [Using pgEdge Anonymizer](docs/usage.md)
  - Creating and Using Patterns
      - [Creating a User-Defined Pattern](docs/custom_pattern.md)
      - [Using Pre-defined Patterns](docs/pattern.md)
      - [Example - Configuration File](docs/sample_config.md)
  - [Troubleshooting](docs/troubleshooting.md)
  - [Release Notes](docs/changelog.md)
  - [Licence](docs/LICENCE.md)

pgEdge Anonymizer is a  command-line tool for anonymizing personally identifiable information (PII) in PostgreSQL databases. The tool replaces sensitive data with realistic fake values that you can use for development and testing, while maintaining data consistency and referential integrity.

## Features

- **100+ built-in patterns** for common PII types across 19 countries
- **Consistent replacement** - same input produces same output within a run
- **Foreign key awareness** - automatically handles CASCADE relationships
- **Large database support** - efficient batch processing with server-side cursors
- **Format preservation** - maintains original data formatting where possible
- **Single transaction** - all changes committed atomically or rolled back
- **Extensible** - define custom patterns using date, number, or mask formats


## Quick Start

Anonymizer lets you create an experimental data set that preserves the shape and integrity of a Postgres database in just three steps:

1. Create a configuration file that specifies the replacement patterns for your columns.
2. Build and run the `pgedge-anonymizer` to convert your columns.
3. Review the results.

Before running `pgedge-anonymizer`, you need to create a [configuration file](docs/configuration.md) named `pgedge-anonymizer.yaml`; the file should contain:

   * a [`database` section](docs/configuration.md#specifying-properties-in-the-database-section), with connection details for your database.
   * a [`columns` section](docs/configuration.md#specifying-properties-in-the-columns-section), listing the fully-qualified columns that you wish to anonymize (in `schema_name.table_name.column name` format).
   * [`patterns` properties](docs/configuration.md#specifying-properties-in-the-pattern-section) for each column that specifies the form that replacement content will take.

For example:

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

After creating a configuration file, [run the anonymizer](docs/usage.md):

```bash
pgedge-anonymizer run
```

Review the list of changes as `pgedge-anonymizer` runs, displaying statistics:

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


## Developer Notes

**Prerequisites**

- Go 1.24 or later
- PostgreSQL (for integration tests)
- Python 3.12+ (for documentation)

Use the following command to build pgedge-anonymizer:

```bash
make build        # Build binary
```

Use the following command to run the Anonymizer test suite:

```bash
make test
```

Use the following command to run the Go Linter:

```bash
make lint
```

Use the following command to format the code:

```bash
make format
```

## License

Copyright 2025 pgEdge, Inc. All rights reserved.

## Support

- [GitHub Issues](https://github.com/pgEdge/pgedge-anonymizer/issues)
- Full documentation is available at [the pgEdge website](https://docs.pgedge.com/pgedge-anonymizer/).
