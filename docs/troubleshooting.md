# Reviewing Online Help and Troubleshooting

To review online help on all command-line options, use the command:

```bash
pgedge-anonymizer help
```

To see pgEdge Anonymizer version information:

```bash
pgedge-anonymizer version
```

## Troubleshooting

### Common Issues

**"column not found in database"**

- Verify the schema.table.column path is correct
- Check that the user has SELECT permission on the table

**"unknown pattern"**

- Ensure the pattern name matches exactly (case-sensitive)
- Check that custom patterns are properly defined

**"failed to connect"**

- Verify database credentials
- Check network connectivity
- Ensure PostgreSQL is accepting connections


## Getting Help

- [GitHub Issues](https://github.com/pgEdge/pgedge-anonymizer/issues)
- [Documentation](https://pgedge.github.io/pgedge-anonymizer/)
