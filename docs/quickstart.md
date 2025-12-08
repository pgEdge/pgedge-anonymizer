# pgEdge Anonymizer Tutorial

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