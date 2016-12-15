## Changelog

0.2.1
-----
* Added `GetSyslogMsgDetails` function to compute the syslog facility and severity when used as syslog receiver (when `enable_syslog_format_only` set to `true`)
* UDP connection is now properly closed when process is stoped.
* Fixed issue where 0 byte length message was being processed upon closing of the UDP connection.

0.2.0
-----
* Added configuration option `json_document_type_schema` which gives ability to optionally enable and enforce a JSON schema validation for JSON format messages.
* Added option `enable_syslog_format_only`, which gives the ability to run process in a mode that only accepts syslog type messages.  This could be used as a logging replacement for processes that log to local syslog address.

0.1.0
-----
* First initial beta release.
