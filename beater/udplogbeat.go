package beater

import (
	"fmt"
	"net"
	"os"
	"os/signal"
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
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Udplogbeat{
		done:               make(chan struct{}),
		config:             config,
		jsonDocumentSchema: map[string]gojsonschema.JSONLoader{},
	}

	for name, path := range config.JsonDocumentTypeSchema {
		logp.Info("Loading JSON schema %s from %s", name, path)
		schemaLoader := gojsonschema.NewReferenceLoader("file://" + path)
		ds := schemaLoader
		bt.jsonDocumentSchema[name] = ds
	}

	bt.config.Addr = fmt.Sprintf("127.0.0.1:%d", bt.config.Port)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		logp.Info("Caught interrupt signal, terminating udplogbeat.")
		os.Exit(0)
	}()

	return bt, nil
}

func (bt *Udplogbeat) Run(b *beat.Beat) error {
	logp.Info("udplogbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	counter := 1

	addr, err := net.ResolveUDPAddr("udp", bt.config.Addr)
	l, err := net.ListenUDP(addr.Network(), addr)

	logp.Info("Listening on %s (UDP)", bt.config.Addr)

	if err != nil {
		return err
	}
	udpBuf := make([]byte, bt.config.MaxMessageSize)
	var event common.MapStr

	logp.Info("Loaded schemas: %v", bt.jsonDocumentSchema)

	for {

		select {
		case <-bt.done:
			return nil
		default:
		}

		logp.Info("Reading from UDP socket...")

		// Events should be in the format of: [FORMAT]:[ES_TYPE]:[EVENT_DATA]
		logSize, _, err := l.ReadFrom(udpBuf)

		if err != nil {
			e, ok := err.(net.Error)
			if ok && e.Timeout() {
				logp.Err("Timeout reading from socket: %v", err)
				continue
			}
		}

		logFormat, logType, logData, err := udploglib.GetLogItem(udpBuf[:logSize])
		if err != nil {
			logp.Err("Error parsing log item: %v", err)
			continue
		}

		logp.Info("Total log item bytes: %d", logSize)
		logp.Info("Format: %s", logFormat)
		logp.Info("ES Type: %s", logType)
		logp.Info("Data: %s", logData)
		logp.Info("---------------------------------------------------")

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
			logp.Info("Event: %v", event)
		} else {
			event["message"] = logData
		}

	SendFailedMsg:
		event["@timestamp"] = common.Time(time.Now())
		event["type"] = logType
		event["counter"] = counter

		bt.client.PublishEvent(event)
		logp.Info("Event sent")
		counter++

	}
}

func (bt *Udplogbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
