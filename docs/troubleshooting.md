# Reviewing Online Help and Troubleshooting

To review online help on all command-line options, use the command:

```bash
pgedge-anonymizer help
```

## Troubleshooting

### Column Not Found in Database

- Verify the schema.table.column path is correct
- Check that the user has SELECT permission on the table

### Unknown Pattern

- Ensure the pattern name matches exactly (case-sensitive)
- Check that custom patterns are properly defined

### Failed to Connect

- Verify database credentials
- Check network connectivity
- Ensure PostgreSQL is accepting connections


## Getting Help

- [GitHub Issues](https://github.com/pgEdge/pgedge-anonymizer/issues)
- [Documentation](https://pgedge.github.io/pgedge-anonymizer/)
