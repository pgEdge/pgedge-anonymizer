# Example: Configuration File

The following code sample shows a complete configuration file that contains all sections:

```yaml
# pgEdge Anonymizer Configuration

# Database connection
database:
  host: db.example.com
  port: 5432
  database: production_copy
  user: anonymizer
  password: ""  # Will prompt or use PGPASSWORD
  sslmode: require

# Pattern configuration
patterns:
  user_path: ./custom-patterns.yaml
  disable_defaults: false

# Columns to anonymize
columns:
  # Customer PII
  - column: public.customers.name
    pattern: PERSON_NAME

  - column: public.customers.email
    pattern: EMAIL

  - column: public.customers.phone
    pattern: US_PHONE

  - column: public.customers.dob
    pattern: DOB_OVER_18

  # Employee data
  - column: hr.employees.ssn
    pattern: US_SSN

  - column: hr.employees.passport_number
    pattern: PASSPORT

  # UK-specific data
  - column: uk.customers.ni_number
    pattern: UK_NI

  - column: uk.customers.nhs_number
    pattern: UK_NHS

  # Notes and free text
  - column: public.support_tickets.description
    pattern: LOREMIPSUM
```
