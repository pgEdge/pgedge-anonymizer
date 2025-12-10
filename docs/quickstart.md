# pgEdge Anonymizer Tutorial

Anonymizer lets you create an experimental data set that preserves the shape and integrity of a Postgres database in just three steps:

1. Create a configuration file that specifies the replacement patterns for your columns.
2. Build and run the `pgedge-anonymizer` to convert your columns.
3. Review the results.

Before running `pgedge-anonymizer`, you need to create a [configuration file](configuration.md) named `pgedge-anonymizer.yaml`; the file should contain:

   * a [`database` section](configuration.md#specifying-properties-in-the-database-section), with connection details for your database.
   * a [`columns` section](configuration.md#specifying-properties-in-the-columns-section), listing the fully-qualified columns that you wish to anonymize (in `schema_name.table_name.column_name` format).
   * [`patterns` properties](configuration.md#specifying-properties-in-the-pattern-section) for each column that specifies the form that replacement content will take.

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

After creating a configuration file, [run the anonymizer](usage.md):

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
