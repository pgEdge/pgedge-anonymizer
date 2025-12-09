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

Include the following [command line options](configuration.md#using-command-line-options) as needed:

| Flag            | Description                                                    |
|-----------------|----------------------------------------------------------------|
| `--config, -c`  | Path to Configuration File (default: `pgedge-anonymizer.yaml`) |
| `--quiet, -q`   | Suppress progress output                                       |
| `--host`        | Database host (overrides value in configuration file)          |
| `--port`        | Database port (overrides value in configuration file)          |
| `--database, -d`| Database name (overrides value in configuration file)          |
| `--user, -U`    | Database user (overrides value in configuration file)          |
| `--password`    | Database password (overrides value in configuration file)      |
| `--sslmode`     | SSL mode (overrides value in configuration file)               |


To review online help, use the command:

```bash
pgedge-anonymizer help
``
