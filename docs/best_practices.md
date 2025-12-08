# Best Practices

Before anonymizing, you should:

1. **Back up your database** - Anonymization is irreversible
2. **Test on a copy** - Validate on a non-production database first
3. **Review columns** - Ensure all PII columns are included
4. **Check foreign keys** - Understand CASCADE relationships


## Performance

For large databases, pgEdge Anonymizer uses:

- **Server-side cursors**: Rows are fetched in configurable batches
  (default 10,000)
- **Batch updates**: Multiple rows updated in single statements using
  CTID-based unnest operations
- **Tiered caching**: LRU in-memory cache with SQLite spillover for
  value dictionaries

For the best performance:

- Process during low-traffic periods
- Ensure adequate disk space for transaction logs
- Consider table-level locks for very large updates
- Monitor PostgreSQL logs for any issues

## Security Considerations

- Run with minimal required privileges
- Use SSL for database connections in production
- Secure your configuration file (contains connection details)
- Consider running on the database server to avoid network transfer of
  sensitive data


