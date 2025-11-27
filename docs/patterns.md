# Patterns Reference

This page provides detailed documentation for all built-in anonymization
patterns and guidance on creating custom patterns.

## Built-in Patterns

pgEdge Anonymizer includes over 180 built-in patterns covering common PII
types, with extensive support for country-specific data formats across 19
countries: Australia, Canada, Finland, France, Germany, India, Ireland,
Italy, Japan, Mexico, New Zealand, Norway, Pakistan, Singapore, South Korea,
Spain, Sweden, United Kingdom, and United States.

### Person Names

#### PERSON_NAME

Generates realistic full names (first and last name combined).

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| John Smith | Jennifer Williams |
| JANE DOE | MICHAEL JOHNSON |
| alice jones | robert davis |
| Smith, John | Williams, Jennifer |

**Format Preservation:**

- Detects comma-separated format (Last, First)
- Preserves all-uppercase input
- Preserves all-lowercase input

---

#### PERSON_FIRST_NAME

Generates first names only.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| John | Michael |
| SARAH | JENNIFER |
| alice | robert |

**Format Preservation:**

- Preserves all-uppercase input
- Preserves all-lowercase input

---

#### PERSON_LAST_NAME

Generates last names only.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| Smith | Williams |
| JOHNSON | DAVIS |
| williams | anderson |

**Format Preservation:**

- Preserves all-uppercase input
- Preserves all-lowercase input

---

### Date of Birth

#### DOB

Generates dates of birth for any age (0-100 years).

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 1985-03-15 | 1967-08-22 |
| 03/15/1985 | 08/22/1967 |
| March 15, 1985 | August 22, 1967 |

**Format Detection:**

- ISO format: YYYY-MM-DD
- US format: MM/DD/YYYY
- US short: MM/DD/YY
- Long format: Month DD, YYYY

---

#### DOB_OVER_13

Generates dates of birth ensuring the person is at least 13 years old.

---

#### DOB_OVER_16

Generates dates of birth ensuring the person is at least 16 years old.

---

#### DOB_OVER_18

Generates dates of birth ensuring the person is at least 18 years old.
Commonly used for adult-only services.

---

#### DOB_OVER_21

Generates dates of birth ensuring the person is at least 21 years old.
Commonly used for alcohol-related services in the US.

---

### Contact Information

#### EMAIL

Generates realistic email addresses using combinations of first names,
last names, and common domain names. Each generated email includes a unique
hash-based suffix derived from the input to ensure uniqueness even with
large datasets and unique constraints.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| john.smith@company.com | michael.williams.e3b0c4@example.com |
| jane@work.org | jennifer.davis.7d865e@example.net |
| user123@mail.co.uk | robert.jones.2c26b4@example.org |

**Features:**

- Unique suffix ensures no collisions with unique database constraints
- Same input always produces same output (deterministic)
- Multiple format variations (first.last, flast, firstl, first_last)

---

#### US_PHONE

Generates US-format phone numbers.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 555-123-4567 | 555-234-5678 |
| (555) 123-4567 | (555) 234-5678 |
| 5551234567 | 5552345678 |
| 555.123.4567 | 555.234.5678 |

**Format Preservation:**

- Detects and preserves separator style (dashes, dots, spaces)
- Detects and preserves parentheses around area code
- Preserves digit-only format

---

#### UK_PHONE

Generates UK-format phone numbers with country code.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| +44 20 7946 0958 | +44 20 8234 5678 |
| 020 7946 0958 | 020 8234 5678 |

---

#### INTERNATIONAL_PHONE

Generates international phone numbers with country codes.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| +1 555 123 4567 | +33 1234 5678901 |
| +44 20 7946 0958 | +49 30 12345678 |

---

#### WORLDWIDE_PHONE

Generates generic phone numbers (digits only, most permissive).

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| Any phone format | 5552345678 |

---

#### ADDRESS

Generates street addresses from diverse worldwide data. This pattern randomly
selects a country format (US, UK, German, Japanese, etc.) and generates an
address in that country's style.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 123 Main St, City 12345 | 742 Oak Ave, Springfield 90210 |
| 456 Oak Avenue | Hauptstraße 891, 10115 Berlin |
| Simple address | 〒100-0001 東京都 1-2-3 |

**Generated Format:**

- Randomly selects from 19 country-specific formats
- Street number, street name, city, and postal code
- Country-appropriate address structure

For country-specific addresses, use patterns like `US_ADDRESS`, `UK_ADDRESS`,
`DE_ADDRESS`, etc.

---

#### WORLDWIDE_ADDRESS

Alias for ADDRESS. Generates street addresses from diverse worldwide data.

---

#### CITY

Generates city names from diverse worldwide data. Combines city names from
all 19 supported countries.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| New York | Tokyo |
| CHICAGO | MÜNCHEN |
| boston | sydney |

**Format Preservation:**

- Preserves all-uppercase input
- Preserves all-lowercase input

For country-specific cities, use patterns like `US_CITY`, `UK_CITY`,
`DE_CITY`, etc.

---

#### WORLDWIDE_CITY

Alias for CITY. Generates city names from diverse worldwide data.

---

### Postal Codes

#### US_ZIP

Generates US ZIP codes in 5-digit or ZIP+4 format.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 12345 | 90210 |
| 12345-6789 | 90210-1234 |

**Format Preservation:**

- Detects 5-digit vs ZIP+4 format
- Preserves hyphen separator for ZIP+4

---

#### UK_POSTCODE

Generates UK postcodes in valid format.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| SW1A 1AA | EC2A 3BT |
| M1 1AE | B33 8TH |
| SW1A1AA | EC2A3BT |

**Features:**

- Generates valid UK postcode patterns (A9 9AA, A99 9AA, AA9 9AA,
  AA99 9AA)
- Preserves space separator
- Uses valid postcode letter combinations

---

#### CA_POSTCODE

Generates Canadian postcodes in valid format.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| K1A 0B1 | M5V 3L9 |
| K1A0B1 | M5V3L9 |

**Features:**

- Generates valid Canadian format (A9A 9A9)
- Preserves space separator
- Uses valid postal code letters (excludes D, F, I, O, Q, U in certain
  positions)

---

#### WORLDWIDE_POSTCODE

Generates postcodes in various international formats. Auto-detects the
input format and generates matching output.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 12345 | 90210 (US format) |
| SW1A 1AA | EC2A 3BT (UK format) |
| K1A 0B1 | M5V 3L9 (Canadian format) |

**Features:**

- Auto-detects US, UK, or Canadian format from input
- Falls back to random format selection for ambiguous input
- Useful when postal code format varies within a column

---

### Financial Information

#### CREDIT_CARD

Generates valid credit card numbers that pass Luhn checksum validation.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 4532-1234-5678-9012 | 4532-8765-4321-0987 |
| 4532123456789012 | 4532876543210987 |
| 4532 1234 5678 9012 | 4532 8765 4321 0987 |

**Features:**

- Generates valid Luhn checksum
- Preserves card type prefix (Visa starts with 4, etc.)
- Preserves separator style (dashes, spaces, none)

---

#### CREDIT_CARD_EXPIRY

Generates credit card expiration dates (always in the future).

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 12/25 | 08/27 |
| 12/2025 | 08/2027 |

---

#### CREDIT_CARD_CVV

Generates 3 or 4 digit CVV codes.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 123 | 456 |
| 1234 | 7890 |

**Format Preservation:**

- 3-digit input produces 3-digit output
- 4-digit input produces 4-digit output (Amex style)

---

### Government Identifiers

#### US_SSN

Generates US Social Security Numbers in valid format.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 123-45-6789 | 234-56-7890 |
| 123456789 | 234567890 |

**Features:**

- Avoids invalid SSN ranges (000, 666, 900-999 area numbers)
- Preserves separator style

---

#### UK_NI

Generates UK National Insurance numbers in valid format.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| AB123456C | CD789012E |
| AB 12 34 56 C | EF 78 90 12 G |

**Format:**

- Two prefix letters (excluding D, F, I, Q, U, V)
- Six digits
- One suffix letter (A, B, C, or D)

---

#### UK_NHS

Generates UK NHS numbers with valid check digit.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 485 777 3456 | 234 567 8901 |
| 4857773456 | 2345678901 |

**Features:**

- 10 digits total
- Valid modulus 11 check digit
- Preserves space formatting

---

#### PASSPORT

Generates passport-style document numbers.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| A12345678 | B98765432 |
| 123456789 | A23456789 |

**Format:**

- 9 alphanumeric characters
- Mix of letters and numbers

---

### Text Content

#### LOREMIPSUM

Generates lorem ipsum placeholder text matching the approximate length
of the original content.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| Short text | Lorem ipsum |
| A longer piece of text | Lorem ipsum dolor sit amet |

**Features:**

- Matches approximate word count of input
- Uses standard lorem ipsum vocabulary
- Suitable for notes, comments, and free-text fields

---

### Network Identifiers

#### IPV4_ADDRESS

Generates IPv4 addresses.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 192.168.1.100 | 10.45.23.187 |
| 10.0.0.1 | 172.16.55.12 |

**Features:**

- Generates valid IPv4 format (x.x.x.x)
- Mixes private and public-looking ranges
- Avoids reserved ranges (0.x.x.x, 127.x.x.x, multicast)

---

#### IPV6_ADDRESS

Generates IPv6 addresses.

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| 2001:0db8:85a3:0000:0000:8a2e:0370:7334 | 4f2c:8a91:3d7e:1b4a:9c2f:5d8e:7a3b:2c1d |
| 2001:db8::1 | 2001:db8:a1b2:c3d4::e5f6 |

**Features:**

- Supports full and compressed (::) formats
- Preserves uppercase/lowercase preference
- Generates valid hex groups

---

#### HOSTNAME

Generates hostnames and fully qualified domain names (FQDNs).

**Input/Output Examples:**

| Input | Output |
|-------|--------|
| webserver | proxy |
| server01.example.com | node42.internal |
| db | api |

**Features:**

- Detects FQDN vs simple hostname
- Preserves numeric suffixes (e.g., web01 → node42)
- Uses realistic server naming conventions

---

### Country-Specific Patterns

pgEdge Anonymizer provides extensive country-specific patterns for names,
phone numbers, addresses, postal codes, and ID numbers. This allows you to
generate data that matches the format and style expected for a specific
country.

#### Country-Specific Phone Numbers

| Pattern | Description | Example Output |
|---------|-------------|----------------|
| `AU_PHONE` | Australian | 0412 345 678 |
| `CA_PHONE` | Canadian | (555) 555-0123 |
| `DE_PHONE` | German | +49 30 12345678 |
| `ES_PHONE` | Spanish | +34 612 345 678 |
| `FI_PHONE` | Finnish | +358 40 123 4567 |
| `FR_PHONE` | French | +33 1 23 45 67 89 |
| `IE_PHONE` | Irish | +353 87 123 4567 |
| `IN_PHONE` | Indian | +91 98765 43210 |
| `IT_PHONE` | Italian | +39 333 123 4567 |
| `JP_PHONE` | Japanese | +81 3-1234-5678 |
| `KR_PHONE` | South Korean | +82 10-1234-5678 |
| `MX_PHONE` | Mexican | +52 55 1234 5678 |
| `NO_PHONE` | Norwegian | +47 912 34 567 |
| `NZ_PHONE` | New Zealand | +64 21 123 4567 |
| `PK_PHONE` | Pakistani | +92 300 1234567 |
| `SE_PHONE` | Swedish | +46 70 123 45 67 |
| `SG_PHONE` | Singaporean | +65 9123 4567 |

---

#### Country-Specific Postcodes

| Pattern | Description | Format |
|---------|-------------|--------|
| `AU_POSTCODE` | Australian | 4 digits (2000-7999) |
| `DE_POSTCODE` | German PLZ | 5 digits |
| `ES_POSTCODE` | Spanish | 5 digits |
| `FI_POSTCODE` | Finnish | 5 digits |
| `FR_POSTCODE` | French | 5 digits (department + commune) |
| `IE_POSTCODE` | Irish Eircode | A9A A9A9 format |
| `IN_POSTCODE` | Indian PIN | 6 digits |
| `IT_POSTCODE` | Italian CAP | 5 digits |
| `JP_POSTCODE` | Japanese | XXX-XXXX format |
| `KR_POSTCODE` | South Korean | 5 digits |
| `MX_POSTCODE` | Mexican | 5 digits |
| `NO_POSTCODE` | Norwegian | 4 digits |
| `NZ_POSTCODE` | New Zealand | 4 digits |
| `PK_POSTCODE` | Pakistani | 5 digits |
| `SE_POSTCODE` | Swedish | XXX XX format |
| `SG_POSTCODE` | Singaporean | 6 digits |

---

#### Country-Specific Addresses

| Pattern | Description | Example Output |
|---------|-------------|----------------|
| `AU_ADDRESS` | Australian | 123 Main Street, Sydney 2000 |
| `CA_ADDRESS` | Canadian | 123 Main St, Toronto M5V 3K9 |
| `DE_ADDRESS` | German | Hauptstraße 123, 10115 Berlin |
| `ES_ADDRESS` | Spanish | Calle Mayor 123, 28001 Madrid |
| `FI_ADDRESS` | Finnish | Mannerheimintie 123, 00100 Helsinki |
| `FR_ADDRESS` | French | 123 Rue de Paris, 75001 Paris |
| `IE_ADDRESS` | Irish | 123 Main Street, Dublin, D02 X285 |
| `IN_ADDRESS` | Indian | 123 MG Road, Mumbai - 400001 |
| `IT_ADDRESS` | Italian | Via Roma 123, 00100 Roma |
| `JP_ADDRESS` | Japanese | 〒100-0001 Tokyo 1-2-3 |
| `KR_ADDRESS` | South Korean | Seoul Gangnam-ro 123 (06000) |
| `MX_ADDRESS` | Mexican | Calle Principal #123, 01000 México |
| `NO_ADDRESS` | Norwegian | Storgata 123, 0001 Oslo |
| `NZ_ADDRESS` | New Zealand | 123 Queen Street, Auckland 1010 |
| `PK_ADDRESS` | Pakistani | 123 Main Road, Islamabad - 44000 |
| `SE_ADDRESS` | Swedish | Kungsgatan 123, 100 00 Stockholm |
| `SG_ADDRESS` | Singaporean | Blk 123 Orchard Road, Singapore 018956 |
| `UK_ADDRESS` | UK | 123 High Street, London, SW1A 1AA |
| `US_ADDRESS` | US | 123 Main St, New York 10001 |

---

#### Worldwide Name Generators

For data that may come from any country, use the worldwide generators:

| Pattern | Description |
|---------|-------------|
| `WORLDWIDE_CITY` | City names from any supported country |
| `WORLDWIDE_FIRST_NAME` | First names from any supported country |
| `WORLDWIDE_LAST_NAME` | Last names from any supported country |
| `WORLDWIDE_NAME` | Full names from any supported country |

---

#### Country-Specific Names

For each supported country, generators are available for first names, last
names, full names, and cities. The pattern names follow the format
`{COUNTRY_CODE}_{TYPE}`:

| Pattern Type | Countries Available |
|--------------|---------------------|
| `{CC}_CITY` | AU, CA, DE, ES, FI, FR, IE, IN, IT, JP, KR, MX, NO, NZ, PK, SE, SG, UK, US |
| `{CC}_FIRST_NAME` | AU, CA, DE, ES, FI, FR, IE, IN, IT, JP, KR, MX, NO, NZ, PK, SE, SG, UK, US |
| `{CC}_LAST_NAME` | AU, CA, DE, ES, FI, FR, IE, IN, IT, JP, KR, MX, NO, NZ, PK, SE, SG, UK, US |
| `{CC}_NAME` | AU, CA, DE, ES, FI, FR, IE, IN, IT, JP, KR, MX, NO, NZ, PK, SE, SG, UK, US |

**Examples:**

- `US_FIRST_NAME` - American first names (John, Mary, Michael)
- `JP_LAST_NAME` - Japanese surnames (Tanaka, Suzuki, Sato)
- `DE_NAME` - German full names (Hans Mueller, Anna Schmidt)
- `UK_CITY` - UK cities with counties (London, Greater London)

---

#### Country-Specific ID Numbers

| Pattern | Description | Format |
|---------|-------------|--------|
| `AU_TFN` | Australian Tax File Number | 9 digits |
| `CA_SIN` | Canadian Social Insurance Number | XXX-XXX-XXX |
| `DE_STEUERID` | German tax ID | 11 digits |
| `ES_NIF` | Spanish NIF/DNI | 8 digits + letter |
| `FI_HETU` | Finnish personal identity code | DDMMYY-XXXC |
| `FR_NIR` | French social security number | 15 digits |
| `IE_PPS` | Irish PPS number | 7 digits + 1-2 letters |
| `IN_AADHAAR` | Indian Aadhaar | 12 digits |
| `IN_PAN` | Indian PAN | AAAAA9999A |
| `IT_CF` | Italian Codice Fiscale | 16 alphanumeric |
| `JP_MYNUMBER` | Japanese My Number | 12 digits |
| `KR_RRN` | South Korean Resident Registration | YYMMDD-XXXXXXX |
| `MX_CURP` | Mexican CURP | 18 alphanumeric |
| `NO_FNR` | Norwegian Fødselsnummer | 11 digits |
| `NZ_IRD` | New Zealand IRD number | 8-9 digits |
| `PK_CNIC` | Pakistani CNIC | XXXXX-XXXXXXX-X |
| `SE_PNR` | Swedish personnummer | YYMMDD-XXXX |
| `SG_NRIC` | Singaporean NRIC | Letter + 7 digits + letter |
| `US_SSN` | US Social Security Number | XXX-XX-XXXX |

---

## Custom Patterns

You can define custom patterns in a separate YAML file. Custom patterns can
either reference built-in generators or define format-based patterns that
generate data based on format strings.

### Pattern File Format

```yaml
patterns:
  - name: PATTERN_NAME
    replacement: "generator name or format hint"
    note: "Human-readable description"
```

### Format-Based Patterns

Format-based patterns allow you to define custom data formats using format
codes. This is useful when you need to generate data in a specific format
that isn't covered by the built-in patterns.

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

#### Auto-Detection

The `type` field is optional. If omitted, the type is auto-detected based
on the format string:

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

### Example Custom Patterns File

See `examples/custom_patterns.yaml` for a comprehensive example of custom
format patterns.

### Using Custom Patterns

Reference your custom patterns file in the configuration:

```yaml
patterns:
  user_path: ./custom-patterns.yaml

columns:
  - column: public.employees.employee_id
    pattern: EMPLOYEE_ID
```

Or via command line:

```bash
pgedge-anonymizer run --config config.yaml --patterns custom-patterns.yaml
```

!!! warning
    Custom pattern names must not conflict with built-in patterns unless
    you set `disable_defaults: true` in your configuration.

## Pattern Selection Guide

Choose patterns based on your data type:

### Generic Patterns

| Data Type | Recommended Pattern |
|-----------|-------------------|
| Full names | `PERSON_NAME` or `WORLDWIDE_NAME` |
| First names only | `PERSON_FIRST_NAME` or `WORLDWIDE_FIRST_NAME` |
| Last names only | `PERSON_LAST_NAME` or `WORLDWIDE_LAST_NAME` |
| Email addresses | `EMAIL` |
| Phone numbers (various) | `WORLDWIDE_PHONE` |
| Street addresses | `ADDRESS` |
| City names | `CITY` or `WORLDWIDE_CITY` |
| Mixed postal codes | `WORLDWIDE_POSTCODE` |
| Credit card numbers | `CREDIT_CARD` |
| Card expiry dates | `CREDIT_CARD_EXPIRY` |
| CVV codes | `CREDIT_CARD_CVV` |
| Passport numbers | `PASSPORT` |
| Birth dates (any age) | `DOB` |
| Birth dates (13+) | `DOB_OVER_13` |
| Birth dates (16+) | `DOB_OVER_16` |
| Birth dates (18+) | `DOB_OVER_18` |
| Birth dates (21+) | `DOB_OVER_21` |
| Notes/comments | `LOREMIPSUM` |
| IPv4 addresses | `IPV4_ADDRESS` |
| IPv6 addresses | `IPV6_ADDRESS` |
| Hostnames/FQDNs | `HOSTNAME` |

### Country-Specific Patterns

For data that should match a specific country's format, use the
country-specific patterns. The country code prefix indicates the country:

| Country | Code | Phone | Postcode | Address | Names | ID |
|---------|------|-------|----------|---------|-------|-----|
| Australia | AU | `AU_PHONE` | `AU_POSTCODE` | `AU_ADDRESS` | `AU_NAME` | `AU_TFN` |
| Canada | CA | `CA_PHONE` | `CA_POSTCODE` | `CA_ADDRESS` | `CA_NAME` | `CA_SIN` |
| Finland | FI | `FI_PHONE` | `FI_POSTCODE` | `FI_ADDRESS` | `FI_NAME` | `FI_HETU` |
| France | FR | `FR_PHONE` | `FR_POSTCODE` | `FR_ADDRESS` | `FR_NAME` | `FR_NIR` |
| Germany | DE | `DE_PHONE` | `DE_POSTCODE` | `DE_ADDRESS` | `DE_NAME` | `DE_STEUERID` |
| India | IN | `IN_PHONE` | `IN_POSTCODE` | `IN_ADDRESS` | `IN_NAME` | `IN_AADHAAR`, `IN_PAN` |
| Ireland | IE | `IE_PHONE` | `IE_POSTCODE` | `IE_ADDRESS` | `IE_NAME` | `IE_PPS` |
| Italy | IT | `IT_PHONE` | `IT_POSTCODE` | `IT_ADDRESS` | `IT_NAME` | `IT_CF` |
| Japan | JP | `JP_PHONE` | `JP_POSTCODE` | `JP_ADDRESS` | `JP_NAME` | `JP_MYNUMBER` |
| Mexico | MX | `MX_PHONE` | `MX_POSTCODE` | `MX_ADDRESS` | `MX_NAME` | `MX_CURP` |
| New Zealand | NZ | `NZ_PHONE` | `NZ_POSTCODE` | `NZ_ADDRESS` | `NZ_NAME` | `NZ_IRD` |
| Norway | NO | `NO_PHONE` | `NO_POSTCODE` | `NO_ADDRESS` | `NO_NAME` | `NO_FNR` |
| Pakistan | PK | `PK_PHONE` | `PK_POSTCODE` | `PK_ADDRESS` | `PK_NAME` | `PK_CNIC` |
| Singapore | SG | `SG_PHONE` | `SG_POSTCODE` | `SG_ADDRESS` | `SG_NAME` | `SG_NRIC` |
| South Korea | KR | `KR_PHONE` | `KR_POSTCODE` | `KR_ADDRESS` | `KR_NAME` | `KR_RRN` |
| Spain | ES | `ES_PHONE` | `ES_POSTCODE` | `ES_ADDRESS` | `ES_NAME` | `ES_NIF` |
| Sweden | SE | `SE_PHONE` | `SE_POSTCODE` | `SE_ADDRESS` | `SE_NAME` | `SE_PNR` |
| UK | UK | `UK_PHONE` | `UK_POSTCODE` | `UK_ADDRESS` | `UK_NAME` | `UK_NI`, `UK_NHS` |
| US | US | `US_PHONE` | `US_ZIP` | `US_ADDRESS` | `US_NAME` | `US_SSN` |

**Example: Anonymizing a UK customer database:**

```yaml
columns:
  - column: customers.full_name
    pattern: UK_NAME
  - column: customers.phone
    pattern: UK_PHONE
  - column: customers.address
    pattern: UK_ADDRESS
  - column: customers.postcode
    pattern: UK_POSTCODE
  - column: customers.ni_number
    pattern: UK_NI
```

## Consistency

Within a single anonymization run, the same input value always produces
the same output value. This ensures:

- Referential integrity is maintained
- Related records stay linked
- Reports and analytics remain meaningful

For example, if "john.smith@example.com" is anonymized to
"michael.williams.e3b0c4@example.com", every occurrence of that email in
the database will be replaced with the same anonymized value.
