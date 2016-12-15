package beater

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/hartfordfive/udplogbeat/config"
	"github.com/hartfordfive/udplogbeat/udploglib"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/xeipuuv/gojsonschema"
)

type Udplogbeat struct {
	done               chan struct{}
	config             config.Config
	client             publisher.Client
	jsonDocumentSchema map[string]gojsonschema.JSONLoader
	conn               *net.UDPConn
}

// New creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Udplogbeat{
		done:   make(chan struct{}),
		config: config,
	}

	if bt.config.EnableJsonValidation {

		bt.jsonDocumentSchema = map[string]gojsonschema.JSONLoader{}

		for name, path := range config.JsonDocumentTypeSchema {
			logp.Info("Loading JSON schema %s from %s", name, path)
			schemaLoader := gojsonschema.NewReferenceLoader("file://" + path)
			ds := schemaLoader
			bt.jsonDocumentSchema[name] = ds
		}

	}

	bt.config.Addr = fmt.Sprintf("127.0.0.1:%d", bt.config.Port)

	return bt, nil
}

func (bt *Udplogbeat) Run(b *beat.Beat) error {
	logp.Info("udplogbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	counter := 1

	addr, err := net.ResolveUDPAddr("udp", bt.config.Addr)
	l, err := net.ListenUDP(addr.Network(), addr)
	bt.conn = l

	logp.Info("Listening on %s (UDP)", bt.config.Addr)

	if err != nil {
		return err
	}
	udpBuf := make([]byte, bt.config.MaxMessageSize)
	var event common.MapStr
	var now common.Time
	var logFormat, logType, logData string

	for {

		select {
		case <-bt.done:
			return nil
		default:
		}

		now = common.Time(time.Now())

		// Events should be in the format of: [FORMAT]:[ES_TYPE]:[EVENT_DATA]
		logSize, _, err := bt.conn.ReadFrom(udpBuf)

		if logSize == 0 {
			continue
		}

		if err != nil {
			e, ok := err.(net.Error)
			if ok && e.Timeout() {
				logp.Err("Timeout reading from socket: %v", err)
				continue
			}
		}

		if bt.config.EnableSyslogFormatOnly {
			logFormat = "plain"
			logType = "syslog"
			logData = strings.TrimSpace(string(udpBuf[:logSize]))
			if logData == "" {
				logp.Err("Syslog event is empty")
				continue
			}
		} else {
			parts, err := udploglib.GetLogItem(udpBuf[:logSize])
			logFormat = parts[0]
			logType = parts[1]
			logData = parts[2]
			if err != nil {
				logp.Err("Error parsing log item: %v", err)
				continue
			}
			logp.Info("Size, Format, ES Type: %d bytes, %s, %s", logSize, logFormat, logType)
		}

		event = common.MapStr{}

		if logFormat == "json" {

			if bt.config.EnableJsonValidation {

				if _, ok := bt.jsonDocumentSchema[logType]; !ok {
					logp.Err("No schema found for this type")
					continue
				}

				result, err := gojsonschema.Validate(bt.jsonDocumentSchema[logType], gojsonschema.NewStringLoader(logData))
				if err != nil {
					logp.Err("Error with JSON object: %s", err.Error())
					continue
				}

				if !result.Valid() {
					logp.Err("Invalid document type")
					event["message"] = logData
					event["tags"] = []string{"_udplogbeat_jspf"}
					goto SendFailedMsg
				}
			}

			if err := ffjson.Unmarshal([]byte(logData), &event); err != nil {
				logp.Err("Could not load json formated event: %v", err)
				event["message"] = logData
				event["tags"] = []string{"_udplogbeat_jspf"}
			}
		} else {
			if bt.config.EnableSyslogFormatOnly {
				msg, facility, severity, err := udploglib.GetSyslogMsgDetails(logData)
				if err == nil {
					event["facility"] = facility
					event["severity"] = severity
					event["message"] = msg
				}
			} else {
				event["message"] = logData
			}
		}

	SendFailedMsg:
		event["@timestamp"] = now
		event["type"] = logType
		event["counter"] = counter

		bt.client.PublishEvent(event)
		counter++

	}
}

func (bt *Udplogbeat) Stop() {
	if err := bt.conn.Close(); err != nil {
		logp.Err("Could not close UDP connection before terminating: %v", err)
	}
	bt.client.Close()
	close(bt.done)
}
