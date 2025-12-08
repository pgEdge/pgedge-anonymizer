# Usage

When you invoke Anonymizer, you can specify the location of a configuration file and/or options on the command line.  Configuration options are invoked in the following order (highest to lowest):

1. Command-line flags
2. Configuration file values
3. Environment variable settings

To review online help, use the command:

```bash
pgedge-anonymizer help
``

Before running Anonymizer, validate the configuration setup:

```bash
pgedge-anonymizer validate [flags]
```

Flags:

- `--config, -c`: Configuration file path
- Database connection flags (same as run)

This validation checks:

- Configuration file syntax
- Database connectivity
- Column existence in the database
- Pattern validity


## Running pgEdge Anonymizer

The following command executes anonymization:

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



