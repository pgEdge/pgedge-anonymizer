# Installation

pgEdge Anonymizer is open-source and licensed with the [PostgreSQL license](LICENCE.md).  You can download pgEdge Anonymizer source code from the [pgEdge repository](https://github.com/pgEdge/pgedge-anonymizer).

## Prerequisites

Before building Anonymizer, install:

- Go 1.21 or later
- PostgreSQL 12 or later
- Make (optional, for using Makefile targets)

## Building Anonymizer from Source

To build Anonymizer from source, clone the `pgedge-anonymizer` repository:

```bash
git clone https://github.com/pgedge/pgedge-anonymizer.git
cd pgedge-anonymizer
```

Then, use `make build` to build anonymizer:

```bash
make build
```

The `make build` command installs dependencies (if needed) and creates the `pgedge-anonymizer` binary in the `/bin` directory in your current directory.  

When the installation completes, you can use the following command to see pgEdge Anonymizer help information:

```bash
pgedge-anonymizer help
```
