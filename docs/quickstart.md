# pgEdge Anonymizer Tutorial

Anonymizer lets you create an experimental data set that preserves the shape and integrity of a Postgres database in just three steps:

1. Create a configuration file.
2. Build and run the Anonymizer.
3. View the results.

## Creating a Configuration File

Before 
1. Create a [configuration file](configuration.md) named `pgedge-anonymizer.yaml`; the file should contain: 
   * a [`database` section](configuration.md#database-properties), with connection details for your database.
   * a [`columns` section](configuration.md#columns-section), listing the fully-qualified columns that you wish to anonymize (in `schema_name.table_name.column name` format).

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

1. [Run the anonymizer](usage.md), specifying the keyword `run`:

```bash
pgedge-anonymizer run
```

1. Review the results as the anonymizer runs, displaying statistics:

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