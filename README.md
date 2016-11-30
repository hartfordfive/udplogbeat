# Udplogbeat

Udplogbeat is a custom beats application intended to allow developers to log events to be indexed in Elasticsearch.  Log entries are sent to a local UDP socket and then shipped out via the selected output.
The intended purpose of this tool is to allow any application to easily log messages locally without writting to disk and taking advantage of the beats framework's various built-in outputs and features.


Ensure that this folder is at the following location:
`${GOPATH}/github.com/hartfordfive`

## Getting Started with Udplogbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7
* github.com/pquerna/ffjson/ffjson
* github.com/xeipuuv/gojsonschema

### Configuration Options

- `udplogbeat.port` : The UDP port on which the process will listen (Default = 5000)
- `udplogbeat.max_message_size` : The maximum accepted message size (Default = 1024)
- `udplogbeat.enable_syslog_format_only` : Boolean value indicating if only syslog messages should be accepted. (Default = false)
- `udplogbeat.enable_json_validation` : Boolean value indicating if JSON schema validation should be applied for `json` format messages (Default = false)
- `udplogbeat.publish_failed_json_invalid` : Boolean value indicating if JSON objects should be sent serialized in the event of a failed validation.  This will add the `_udplogbeat_jspf` tag. (Default = false)
- `udplogbeat.json_document_type_schema` :  A hash consisting of the Elasticsearch type as the key, and the absolute local schema file path as the value.

### Configuration Example

Sample configuration for a syslog replacement
```
udplogbeat:
  port: 5000
  max_message_size: 4096
  enable_syslog_format_only: false
```

Sample configuration that enforces schemas for JSON format events
```
udplogbeat:
  port: 5001
  max_message_size: 2048
  enable_json_validation: true
  json_document_type_schema: 
    email_contact: "/etc/udplogbeat/app1_schema.json"
    stock_item: "/etc/udplogbeat/app2_schema.json"
```

JSON schemas can be automatically generated from an object here: http://jsonschema.net/.  You can also view the included sample schemas `app1_schema.json` and `app2_schema.json` as examples.

#### Considerations

If you intend on using this as a drop-in replacement to logging with Rsyslog, this method will not persist your data to a file on disk.  
If udplogbeat is down for any given reason, messages sent to the configured UDP port will never be processed or sent to your ELK cluster.
If you need 100% guarantee each message will be delivered at least once, this may not be the best solution for you.  
If some potential loss of log events is acceptable for you, than this may be a reasonable solution for you.


### Log Structure

In order for the udplogbeat application to accept events, when not in syslog format only mode (*enable_syslog_format_only: false*), they must be structured in the following format:

**[FORMAT]:[ES_TYPE]:[EVENT_DATA]**

* FORMAT : Either `json` or `plain`.  JSON encoded entries will be automatically parsed.
* ES_TYPE : The type to be used in Elasticsearch
* EVENT_DATA : The log entry itself

*Example:*

Plain encoded event:
```
plain:syslog:Nov 26 18:51:42 my-web-host01 dhclient: DHCPACK of 10.2.1.2 from 10.2.1.3
```

JSON encoded event:
```
json:my_application:{"message":"This is a test JSON message", "application":"my\_application", "log\_level":"INFO"}
```

*Please note the current date/time is automatically added to each log entry.*

### Sample Clients

Please see the `sample_clients/` directory for examples of clients in various languages.


### Init Project
To get running with Udplogbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Udplogbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/hartfordfive/udplogbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Udplogbeat run the command below. This will generate a binary
in the same directory with the name udplogbeat.

```
make
```

Or to build the zipped binaries for OSX, Windows and Linux:

./build_os_binaries.sh "[VERSION_NUMBER]"

These will be placed in the `bin/` directory.

### Run

To run Udplogbeat with debugging output enabled, run:

```
./udplogbeat -c udplogbeat.yml -e -d "*"
```


### Test

To test Udplogbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/udplogbeat.template.json and etc/udplogbeat.asciidoc

```
make update
```


### Cleanup

To clean  Udplogbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Udplogbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/hartfordfive
cd ${GOPATH}/github.com/hartfordfive
git clone https://github.com/hartfordfive/udplogbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.

## Author

Alain Lefebvre <hartfordfive 'at' gmail.com>

## License

Covered under the Apache License, Version 2.0
Copyright (c) 2016 Alain Lefebvre