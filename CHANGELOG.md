## Changelog

0.1.0
-----
* First initial beta release.

0.2.0
-----
* Added configuration option `json_document_type_schema` which gives ability to optionally enable and enforce a JSON schema validation for JSON format messages.
* Added option `enable_syslog_format_only`, which gives the ability to run process in a mode that only accepts syslog type messages.  This could be used as a logging replacement for processes that log to local syslog address.
