# Using Custom Patterns

When you build your [configuration file](configuration.md), you will specify the patterns used to replace sensitive data with similar but meaningless values.  pgEdge Anonymizer includes a number of built-in anonymization patterns and can use a combination of [pre-defined patterns](patterns.md) and user-defined patterns when processing a file.  Custom patterns can either reference built-in generators or define format-based patterns that generate data based on format strings.

## Using Custom Patterns

After creating and saving custom patterns in a file, provide a reference to that file in your configuration file (with the `patterns` property):

```yaml
patterns:
  user_path: ./custom-patterns.yaml
```

Or, on the command line when you invoke `pgedge-anonymizer`::

```bash
pgedge-anonymizer run --config config.yaml --patterns custom-patterns.yaml
```

Then, use your custom patterns in the `columns` property of your configuration file:

```yaml
columns:
  - column: public.employees.employee_id
    pattern: EMPLOYEE_ID
```

!!! warning
    Custom pattern names must not conflict with built-in patterns unless you set `disable_defaults: true` in your configuration.


## Pattern File Format

Within your custom pattern file, you need to define and describe any custom patterns that Anonymizer will use when replacing the data.  For example, you might use the following pattern to define the shape of your account identifier column:

```yaml
  - name: ACCOUNT_ID
    format: "%010d"
    type: number
    min: 1000000000
    max: 9999999999
    note: 10-digit account ID
```

The `format: "%010d"` property instructs Anonymizer to replace the account number with a 10 digit integer value.  Descriptions of the formats and examples follow below, and in the `examples` folder of the [GitHub repo](https://github.com/pgEdge/pgedge-anonymizer).

### Format-Based Patterns

Format-based patterns allow you to define custom data formats using format codes. This is useful when you need to generate data in a specific format that isn't covered by the built-in patterns.

There are three types of format patterns:

- **date** - Uses strftime format codes for date/time generation
- **number** - Uses printf format codes for number generation
- **mask** - Uses character placeholders for pattern-based generation

#### Date Format Patterns

Date patterns use strftime-style format codes:

| Code | Description | Example |
|------|-------------|---------|
| `%Y` | 4-digit year | 2024 |
| `%y` | 2-digit year | 24 |
| `%m` | 2-digit month | 03 |
| `%d` | 2-digit day | 15 |
| `%H` | 24-hour hour | 14 |
| `%I` | 12-hour hour | 02 |
| `%M` | Minute | 30 |
| `%S` | Second | 45 |
| `%B` | Full month name | January |
| `%b` | Abbreviated month | Jan |
| `%A` | Full weekday | Monday |
| `%a` | Abbreviated weekday | Mon |
| `%p` | AM/PM | PM |
| `%P` | am/pm | pm |

**Example date format patterns:**

```yaml
patterns:
  - name: DATE_US_FORMAT
    format: "%m/%d/%Y"
    type: date
    note: US-style date (MM/DD/YYYY)

  - name: DATE_ISO
    format: "%Y-%m-%d"
    type: date
    min_year: 1990
    max_year: 2024
    note: ISO date with year range

  - name: DATETIME_FULL
    format: "%Y-%m-%d %H:%M:%S"
    type: date
    note: Full datetime

  - name: DATE_LONG
    format: "%B %d, %Y"
    type: date
    note: Long format (January 15, 2024)
```

#### Number Format Patterns

Number patterns use printf-style format codes:

| Code | Description | Example |
|------|-------------|---------|
| `%d` | Decimal integer | 123 |
| `%08d` | Zero-padded to 8 digits | 00000123 |
| `%5d` | Space-padded to 5 digits | 123 |

**Example number format patterns:**

```yaml
patterns:
  - name: ORDER_NUMBER
    format: "ORD-%08d"
    type: number
    min: 1
    max: 99999999
    note: Order number (ORD-00000001)

  - name: INVOICE_NUMBER
    format: "INV-%06d"
    type: number
    min: 100000
    max: 999999
    note: Invoice number (INV-100000)

  - name: ACCOUNT_ID
    format: "%010d"
    type: number
    min: 1000000000
    max: 9999999999
    note: 10-digit account ID
```

#### Mask Format Patterns

Mask patterns use character placeholders:

| Placeholder | Description | Example Output |
|-------------|-------------|----------------|
| `#` or `9` | Random digit (0-9) | 5 |
| `A` | Uppercase letter (A-Z) | K |
| `a` | Lowercase letter (a-z) | k |
| `X` | Uppercase alphanumeric | K or 5 |
| `x` | Lowercase alphanumeric | k or 5 |
| `*` | Any character | K, k, or 5 |
| `\` | Escape next character | (literal) |

**Example mask format patterns:**

```yaml
patterns:
  - name: PRODUCT_SKU
    format: "SKU-AA-####"
    type: mask
    note: Product SKU (SKU-AB-1234)

  - name: LICENSE_PLATE
    format: "AAA-####"
    type: mask
    note: License plate (ABC-1234)

  - name: EMPLOYEE_ID
    format: "EMP-AAA-###"
    type: mask
    note: Employee ID (EMP-ABC-123)

  - name: SERIAL_NUMBER
    format: "SN\\-AAAA\\-####\\-XX"
    type: mask
    note: Serial number with escaped hyphens
```

### Auto-Detecting the Type Field

The `type` field is optional. If omitted, the type is auto-detected based on the format string:

- Formats containing strftime codes (`%Y`, `%m`, etc.) are detected as
  `date`
- Formats containing printf codes (`%d`, `%08d`, etc.) are detected as
  `number`
- All other formats are treated as `mask`

```yaml
patterns:
  - name: AUTO_DATE
    format: "%Y-%m-%d"
    note: Auto-detected as date type

  - name: AUTO_NUMBER
    format: "%05d"
    min: 1
    max: 99999
    note: Auto-detected as number type

  - name: AUTO_MASK
    format: "##-AAA-##"
    note: Auto-detected as mask type
```