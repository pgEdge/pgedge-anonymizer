# pgEdge Anonymizer

[![CI](https://github.com/pgEdge/pgedge-anonymizer/actions/workflows/ci.yml/badge.svg)](https://github.com/pgEdge/pgedge-anonymizer/actions/workflows/ci.yml)

A command-line tool for anonymizing personally identifiable information (PII)
in PostgreSQL databases. Replace sensitive data with realistic fake values
while maintaining data consistency and referential integrity.

## Features

- **100+ built-in patterns** for common PII types across 19 countries
- **Consistent replacement** - same input produces same output within a run
- **Foreign key awareness** - automatically handles CASCADE relationships
- **Large database support** - efficient batch processing with server-side
  cursors
- **Format preservation** - maintains original data formatting where possible
- **Single transaction** - all changes committed atomically or rolled back
- **Extensible** - define custom patterns using date, number, or mask formats

## Quick Start

### Installation

Download the latest release from the
[releases page](https://github.com/pgEdge/pgedge-anonymizer/releases), or
build from source:

```bash
git clone https://github.com/pgEdge/pgedge-anonymizer.git
cd pgedge-anonymizer
make build
```

### Usage

1. Create a configuration file `pgedge-anonymizer.yaml`:

```yaml
database:
  host: localhost
  port: 5432
  database: myapp

columns:
  - column: public.users.email
    pattern: EMAIL

  - column: public.users.phone
    pattern: US_PHONE

  - column: public.users.ssn
    pattern: US_SSN

  - column: public.users.first_name
    pattern: PERSON_FIRST_NAME

  - column: public.users.last_name
    pattern: PERSON_LAST_NAME
```

2. Run the anonymizer:

```bash
pgedge-anonymizer run --user myuser --password mypassword
```

3. Validate configuration without making changes:

```bash
pgedge-anonymizer validate
```

## Built-in Patterns

### Worldwide Patterns

| Pattern | Description | Example |
|---------|-------------|---------|
| `PERSON_NAME` | Full name | Michael Williams |
| `PERSON_FIRST_NAME` | First name | Jennifer |
| `PERSON_LAST_NAME` | Last name | Anderson |
| `EMAIL` | Email address | john.smith.a1b2c3@example.com |
| `ADDRESS` | Street address | 456 Oak Avenue, London, Greater London |
| `CITY` | City name | Manchester, Greater Manchester |
| `DOB` | Date of birth | 1985-03-15 |
| `CREDIT_CARD` | Credit card number | 4532-1234-5678-9012 |
| `PASSPORT` | Passport number | A12345678 |

### Country-Specific Patterns

Patterns are available for Australia, Canada, Finland, France, Germany,
India, Ireland, Italy, Japan, Mexico, New Zealand, Norway, Pakistan,
Singapore, South Korea, Spain, Sweden, United Kingdom, and United States.

Examples:

| Pattern | Description | Example |
|---------|-------------|---------|
| `US_SSN` | US Social Security Number | 234-56-7890 |
| `US_PHONE` | US phone number | (555) 234-5678 |
| `UK_NI` | UK National Insurance | AB123456C |
| `UK_POSTCODE` | UK postcode | SW1A 1AA |
| `DE_STEUERID` | German tax ID | 12345678901 |
| `FR_NIR` | French social security | 1 85 01 75 123 456 00 |
| `IN_AADHAAR` | Indian Aadhaar | 1234 5678 9012 |
| `JP_MYNUMBER` | Japanese My Number | 123456789012 |

See the [patterns documentation](https://pgedge.github.io/pgedge-anonymizer/patterns/)
for the complete list.

## Configuration

Database connection can be configured via:

- Configuration file (`database` section)
- Environment variables (`PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`,
  `PGPASSWORD`, `PGSSLMODE`)
- Command-line flags (`--host`, `--port`, `--database`, `--user`,
  `--password`, `--sslmode`)

Priority: command-line flags > config file > environment variables.

See the [configuration reference](https://pgedge.github.io/pgedge-anonymizer/configuration/)
for complete details.

## Custom Patterns

Define custom patterns using format strings:

```yaml
patterns:
  - name: ORDER_NUMBER
    format: "ORD-%08d"
    type: number
    min: 1
    max: 99999999

  - name: PRODUCT_SKU
    format: "SKU-AA-####"
    type: mask

  - name: HIRE_DATE
    format: "%Y-%m-%d"
    type: date
    min_year: 2010
    max_year: 2024
```

## Documentation

Full documentation is available at
[pgedge.github.io/pgedge-anonymizer](https://pgedge.github.io/pgedge-anonymizer/).

## Development

### Prerequisites

- Go 1.24 or later
- PostgreSQL (for integration tests)
- Python 3.12+ (for documentation)

### Building

```bash
make build        # Build binary
make test         # Run tests
make lint         # Run linter
make fmt          # Format code
```

### Running Tests

```bash
make test
```

## License

Copyright 2025 pgEdge, Inc. All rights reserved.

## Support

- [GitHub Issues](https://github.com/pgEdge/pgedge-anonymizer/issues)
- [Documentation](https://pgedge.github.io/pgedge-anonymizer/)
