# Best Practices

Anonymizer replaces sensitive PII data with fake values while maintaining data consistency and referential integrity, to create a working copy for experimentation and testing that does not violate PII laws.
    
**Before anonymizing:**

1. **Back up your database** - Anonymization is irreversible.
2. **Test on a copy** - Validate your configuration to ensure that Anonymizer is applied to a non-production database.
3. **Review columns** - Ensure all PII columns are included when obscuring test data.
4. **Check foreign keys** - Understand `CASCADE` relationships within your tables.

To maintain a secure environment while using Anonymizer, you should:

- run Anonymizer with minimal required privileges.
- use SSL for production database connections.
- secure your configuration file (the file contains connection details).
- consider running Anonymizer on the database server to avoid network transfer of
  sensitive data.

When anonymizing large databases, Anonymizer improves performance by using:

- **Server-side cursors**: Rows are fetched in configurable batches
  (default 10,000).
- **Batch updates**: Multiple rows updated in single statements using
  CTID-based unnest operations.
- **Tiered caching**: LRU in-memory cache with SQLite spillover for
  value dictionaries.

To ensure you're getting the best performance, you should:

- anonymize content during low-traffic periods.
- ensure that you have adequate disk space for transaction logs.
- consider using table-level locks for very large updates.
- monitor the PostgreSQL logs for any issues.
