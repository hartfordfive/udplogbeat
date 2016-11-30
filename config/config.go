// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period                            time.Duration `config:"period"`
	Port                              int           `config:"port"`
	MaxMessageSize                    int           `config:"max_message_size"`
	Addr                              string
	EnableSyslogFormatOnly            bool              `config:"enable_syslog_format_only"`
	EnableJsonValidation              bool              `config:"enable_json_validation"`
	PublishFailedJsonSchemaValidation bool              `config:publish_failed_json_schema_validation`
	PublishFailedJsonInvalid          bool              `config:publish_failed_json_invalid`
	JsonDocumentTypeSchema            map[string]string `config:"json_document_type_schema"`
}

var DefaultConfig = Config{
	Period:                            1 * time.Second,
	Port:                              5000,
	MaxMessageSize:                    1024,
	EnableJsonValidation:              false,
	PublishFailedJsonSchemaValidation: false,
	PublishFailedJsonInvalid:          false,
	EnableSyslogFormatOnly:            false,
}
