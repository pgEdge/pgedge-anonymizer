# pgEdge Anonymizer 

pgEdge Anonymizer is a command-line tool for anonymizing personally
identifiable information (PII) in PostgreSQL databases. It replaces
sensitive data with realistic but fake values while maintaining data
consistency and referential integrity.

pgEdge Anonymizer features:

- **Pattern-based anonymization**: 100+ built-in patterns for common PII types
- **Consistent replacement**: Same input values produce the same anonymized
  output within a run
- **Foreign key awareness**: Automatically handles CASCADE relationships
- **Large database support**: Efficient batch processing with server-side
  cursors
- **Format preservation**: Maintains original data formatting where possible
- **Single transaction**: All changes committed atomically or rolled back
- **Extensible**: Define custom patterns for your specific needs

pgEdge Anonymizer automatically analyzes foreign key relationships:

- **CASCADE updates**: If a column has referencing foreign keys with
  `ON UPDATE CASCADE`, the tool updates the source column and PostgreSQL
  propagates changes automatically.

- **Processing order**: Columns are processed in dependency order to
  maintain referential integrity.

- **Skip targets**: Columns that are CASCADE targets of other configured
  columns are automatically skipped to avoid duplicate processing.

**Getting Help**

- [GitHub Issues](https://github.com/pgEdge/pgedge-anonymizer/issues)
- [Documentation](https://pgedge.github.io/pgedge-anonymizer/)

