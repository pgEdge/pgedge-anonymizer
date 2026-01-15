# Installation

pgEdge Anonymizer modifies data in a PostgreSQL database, replacing sensitive
data with realistic but fake values while maintaining referential integrity,
providing your development team with data for experimentation and testing.

!!! warning

    pgEdge Anonymizer is a development tool intended to obscure PII data, and
    should not be applied to data in a production environment. Anonymizing is
    not reversible.

You can install pgEdge Anonymizer with
[pgEdge Enterprise Postgres](https://docs.pgedge.com/enterprise/) packages
or build Anonymizer from source code from the
[pgEdge repository](https://github.com/pgEdge/pgedge-anonymizer).

pgEdge Anonymizer is open-source and licensed with the
[PostgreSQL license](LICENCE.md).

## Prerequisites

Before building Anonymizer, install:

- Go 1.24 or later.
- PostgreSQL 12 or later.
- Make (optional, for using Makefile targets to build from source code).

## Building Anonymizer from Source

To build Anonymizer from source, clone the `pgedge-anonymizer`
repository:

```bash
git clone https://github.com/pgedge/pgedge-anonymizer.git
cd pgedge-anonymizer
```

Then, use `make build` to build anonymizer:

```bash
make build
```

The `make build` command installs dependencies (if needed) and creates the
`pgedge-anonymizer` binary in the `/bin` directory in your current directory.

When the installation completes, you can use the following command to see
pgEdge Anonymizer help information:

```bash
pgedge-anonymizer help
```
