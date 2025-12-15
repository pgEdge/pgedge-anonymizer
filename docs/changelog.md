# Changelog

All notable changes to the pgEdge Anonymizer will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0-beta1] - 2025-12-15

### Changed

- Comprehensive documentation restructure for improved navigation and
  readability
- Split monolithic documentation into focused topic pages:
    - Best practices guide
    - Installation guide
    - Configuration reference
    - Quickstart tutorial
    - Usage guide
    - Custom patterns guide
    - Built-in patterns reference
    - Sample configuration
    - Troubleshooting guide
- Updated MkDocs navigation structure with logical groupings
- Simplified README with links to detailed documentation

## [1.0.0-alpha3] - 2025-12-08

### Fixed

- Fixed config file error handling - commands that require configuration now
  show a clear error message when no config file is found
- Fixed issue where viper would incorrectly try to parse the binary as a
  config file when searching in the binary's directory
- The `version` and `help` commands no longer attempt to load a config file

### Changed

- Config file search now explicitly looks for `pgedge-anonymizer.yaml` to
  avoid accidentally matching other files (such as the binary itself)

## [1.0.0-alpha2] - 2025-11-27

### Added

- Initial release of pgEdge Anonymizer

- **Core Functionality**

    - YAML-based configuration for database connection and column mappings
    - Pattern-based anonymization with built-in and custom patterns
    - Consistent value replacement within a single run (same input produces
      same output)
    - Single transaction processing with rollback on error
    - Foreign key analysis with CASCADE handling

- **100+ Built-in Patterns**

    - Personal: `PERSON_NAME`, `PERSON_FIRST_NAME`, `PERSON_LAST_NAME`,
      plus worldwide and country-specific variants
    - Date of Birth: `DOB`, `DOB_OVER_13`, `DOB_OVER_16`, `DOB_OVER_18`,
      `DOB_OVER_21`
    - Contact: `EMAIL`, `ADDRESS`, `CITY`, plus country-specific addresses
    - Phone: `US_PHONE`, `UK_PHONE`, `INTERNATIONAL_PHONE`, `WORLDWIDE_PHONE`,
      plus 17 additional country-specific formats
    - Postal Codes: `US_ZIP`, `UK_POSTCODE`, `CA_POSTCODE`,
      `WORLDWIDE_POSTCODE`, plus 16 additional country-specific formats
    - Financial: `CREDIT_CARD`, `CREDIT_CARD_EXPIRY`, `CREDIT_CARD_CVV`
    - Government IDs: `US_SSN`, `UK_NI`, `UK_NHS`, `PASSPORT`,
      plus country-specific IDs for AU, CA, DE, ES, FI, FR, IE, IN, IT, JP,
      KR, MX, NO, NZ, PK, SE, SG
    - Network: `IPV4_ADDRESS`, `IPV6_ADDRESS`, `HOSTNAME`
    - Text: `LOREMIPSUM`

- **Country-Specific Support for 19 Countries**

    - Australia (AU), Canada (CA), Finland (FI), France (FR), Germany (DE),
      India (IN), Ireland (IE), Italy (IT), Japan (JP), Mexico (MX),
      New Zealand (NZ), Norway (NO), Pakistan (PK), Singapore (SG),
      South Korea (KR), Spain (ES), Sweden (SE), United Kingdom (UK),
      United States (US)

- **Format Preservation**

    - Phone number format detection (dashes, dots, spaces, parentheses)
    - Date format detection (ISO, US, EU, short, long)
    - Credit card separator preservation
    - Case preservation for names and emails

- **Database Support**

    - PostgreSQL connection via pgx driver
    - Support for standard libpq environment variables (PGHOST, PGPORT, etc.)
    - SSL connection options (disable, require, verify-ca, verify-full)
    - Server-side cursors for efficient batch processing
    - CTID-based batch updates for performance

- **Large Database Support**

    - Tiered caching with LRU in-memory cache
    - SQLite spillover for very large value dictionaries
    - Configurable batch size (default 10,000 rows)

- **CLI Commands**

    - `run` - Execute anonymization with progress output
    - `validate` - Validate configuration without making changes
    - `version` - Display version information
    - `--quiet` flag to suppress progress output
    - Database connection flags to override configuration

- **Documentation**

    - Comprehensive user guide with examples
    - Configuration reference
    - Pattern reference with format examples
