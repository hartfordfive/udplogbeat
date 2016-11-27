// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period         time.Duration `config:"period"`
	Port           int           `config:"port"`
	MaxMessageSize int           `config:"max_message_size"`
	Addr           string
}

var DefaultConfig = Config{
	Period:         1 * time.Second,
	Port:           5000,
	MaxMessageSize: 1024,
}
