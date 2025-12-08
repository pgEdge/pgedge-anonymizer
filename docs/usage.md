# Using pgEdge Anonymizer

When you invoke Anonymizer, you can specify execution preferences in a configuration file, on the command line, or by setting environment variables.  Configuration options are used in the following order:

1. Command-line flags take precedence over configuration file values and environment variables.
2. Configuration file values take precedence over environment variable settings.
3. Environment variable settings are used for details that pgEdge Anonymizer cannot locate on the command-line or in a configuration file.

Before running Anonymizer, validate your configuration details:

```bash
pgedge-anonymizer validate [flags]
```

Include the `validate` keyword and the following configuration options:

- Use `--config, -c` to specify the path to the configuration file.
- Include any database connection flags that you'll be using.

The validation run checks your:

- configuration file syntax.
- database connectivity.
- that the specified columns exist in the database.
- pattern validity.

When you've successfully validated the deployment options, you're ready to run Anonymizer.


## Running pgEdge Anonymizer

Include the run keyword when you invoke `pgedge-anonymizer` to start anonymization:

```bash
pgedge-anonymizer run [flags]
```

Include the following [command line options](configuration.md#command-line-overrides) as needed:

- `--config, -c`: Configuration file path (default: pgedge-anonymizer.yaml)
- `--quiet, -q`: Suppress progress output
- `--host`: Database host (overrides config)
- `--port`: Database port (overrides config)
- `--database, -d`: Database name (overrides config)
- `--user, -U`: Database user (overrides config)
- `--password`: Database password (overrides config)
- `--sslmode`: SSL mode (overrides config)

To review online help, use the command:

```bash
pgedge-anonymizer help
``
