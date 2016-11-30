package udploglib

import (
	//"errors"
	"fmt"
	"strings"
	//"github.com/xeipuuv/gojsonschema"
)

// GetLogItem returns the log entry format, elasticsearch type, message and error (if any)
func GetLogItem(buf []byte) ([]string, error) {

	parts := strings.SplitN(string(buf), ":", 3)
	if len(parts) != 3 {
		return []string{"", "", ""}, fmt.Errorf("Invalid log item")
	}
	if parts[0] != "json" && parts[0] != "plain" {
		return []string{"", "", ""}, fmt.Errorf("Log format %s is invalid", parts[0])
	}
	if parts[1] == "" {
		return []string{"", "", ""}, fmt.Errorf("A log type must be specified")
	}
	if parts[2] == "" {
		return []string{"", "", ""}, fmt.Errorf("Log data is empty")
	}

	return []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])}, nil
}
