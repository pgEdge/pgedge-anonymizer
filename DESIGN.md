# Project Design Notes

## Overview

pgEdge Anonymizer is a tool for anonymizing data such as Personally 
Identifieable Information (PPI) to meet the requirements of GDPR and other 
regulations when cloning production data for development purposes.

## Architecture

The application is built in GoLang, and consists of a single binary
executable, a configuration file defining standard data patterns, and a 
configuration file defining an anonymization task.

## Data Patterns

When run, the tool will first load the standard data patterns. These are 
defined as a list in a YAML file with the following fields:

- Name
- Replacement pattern (regexp-style)
- Note (optional, for notation only)

Default patterns will be provided to match the following, and more:

- US Phone numbers
- UK Phone numbers
- International phone numbers (with country code)
- Worldwide phone numbers (any format, with or without country code etc.)
- People's names
- People's first names
- People's last names
- People's addresses (multi-line)
- Email addresses
- Credit card numbers
- Credit card expiry dates
- Credit card CVV numbers
- US Social Security numbers
- UK NI numbers
- UK NHS numbers
- Passport numbers
- Date of birth (all ages)
- Date of birth (over 13's)
- Date of birth (over 16's)
- Date of birth (over 18's)
- Date of birth (over 21's)

A special "LORUMIPSUM" pattern will be made available, in which the tool will
generate lorem ipsum text of a suitable length for the column being anonymised.

In addition to the standard set of patterns, the user may provide an additional
pattern file if they choose, the contents of which will be merged with the 
default patterns at runtime. Optionally, the default patterns may be disabled
through the configuration file or on the command line.

The default patterns will be stored in pgedge-anonymizer-patterns.yaml, which
will be searched for in (in order of priority), an optional path specified on 
the command line, an optional path specified in the configuration file, 
/etc/pgedge, or the directory containing the binary. A configuration option
will be provided to disable the use of the standard patterns.

The path to the user defined patterns file will be specified in (in order of 
priority) an optional path specified on the command line or an optional path 
specified in the configuration file.

## Configuration File

The configuration file contains the preferences for an anonymization run. It
will be searched for in (in order of priority) an optional path specified on 
the command line, /etc/pgedge/pgedge-anonymizer.yaml, or 
pgedge-anonymizer.yaml in the directory containing the binary.

The configuration file will contain PostgreSQL database connection parameters,
including an optional password field (which shouild be ignored if the server
doesn't require a password), and optional fields for the paths to certificates
required for certificate authentication. Command line options will also be 
offered for all database connection parameters, which will override those in 
the configuration file (if present). If database connection parameters are
not included on the command line, or in the configuration file, fall back to 
using the standard libpq environment variables for connection information.

The configuration file will also contain a list of columns on which 
anonymization should be performed, in schema.table.column format. For each
column, the name of the anonymization pattern to use will also be included.

## Operation

When invoked, the tool will load the configuration file. If it cannot be 
located, and error will be given showing usage information and an error 
message, and the tool will exit.

If the configuration is found, the tool will then attempt to load the 
anonymization patterns. If the default patterns file is enabled in the 
configuration, it will load and validate the file, exiting with an error
if needed. It will then load and merge the user patterns, if a user pattern
file is specified. It will exit with an error is the user patterns file
cannot be read, contains errors, or contains names of patterns which conflict 
with any loaded from the default patterns file.

If all configuration and pattern loading is successful, the tool will then
attempt to connect to the database. If the connection is not successful, it
will exit with a message explaining what the connection problem was.

If the database connection succeeds, the tool will then confirm that all the 
schema.table.column names exist in the database. If any do not, a complete
list will be output with an error message, and the tool will exit.

If the column validation succeeds, the tool will then analyse the database
to understand foreign key structures in which any column is involved, and any
CASCADE or similar options used that may affect column updates. It will 
re-order the columns to be processed as needed to ensure that updates will
succeed, for example, by not attempting to update the child value before the 
parent in a foreign key relationship.

The tool will then start a transaction, and iterate through every row in each
table that requires anonymization. Each column in each table that is being 
anonymized will have a new random value generated that matches the pattern 
specified for the table. The following requirements must be respected:

- Each unique value encountered must always be anonymized with the same new
    value. For example, a name such as "John Doe" might be anonymized to "Fred
    Bloggs" in one column. Every other occurance of "John Doe" in any column
    specified must also be anonymized to "Fred Blogs" to ensure data 
    consistency.
- Parent values in foreign key relationships may be set to CASCADE to child 
    tables. In such cases, no attempt to update the child table should be 
    made (even if specified as an anonymization target) as it is not necessary.
- Very large databases might include millions of values that need to be
    anonymized consistently. The tool should build a temporary dictionary to
    keep track of generated values, and use whatever techniques are appropriate
    to ensure performance is kept to a maximum.
- Anonymized data should look authentic. For fields such as names, it may be 
    necessary to use external data sources to assist with realistic data 
    generation.

If an error occurs at any point during the anonymization process, the 
transaction must be rolled back immediately, and the tool should exit with an
error message explaining what the problem was, identifying the exact datum
affected (where appropriate) and the column identifier so the user can 
investigate the issue.

If no errors occur, the transaction should be committed, and the tool should 
exit after displaying a table of statistics showing how many values were 
anonymized in which columns, and a summary of the number of columns and the
total number of anonymizations.