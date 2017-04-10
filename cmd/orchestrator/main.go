package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/model"
	"github.com/codingconcepts/albert/pkg/orchestrator"
	nats "github.com/nats-io/go-nats"
)

func main() {
	config := mustLoadConfig("config.json")
	initLogger(config)
	mustConnectToNATS(config)

	o := orchestrator.NewOrchestrator(config)
	//go o.Start()
	//defer o.Stop()

	for _, a := range o.Applications {
		config.Logger.WithFields(logrus.Fields{
			"name":       a.Name,
			"schedule":   a.Schedule,
			"percentage": a.Percentage,
		}).Info("application configured")
	}

	// trigger every 5s
	for range time.NewTicker(time.Second * 5).C {
		for _, a := range o.Applications {
			o.Process(a)
		}
	}

	config.Logger.Info("orchestrator started successfully")
	fmt.Scanln()
}

func mustLoadConfig(path string) (c *orchestrator.Config) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	buffer := new(bytes.Buffer)
	if _, err = io.Copy(buffer, file); err != nil {
		log.Fatal(err)
	}

	c = new(orchestrator.Config)
	if err = json.Unmarshal(buffer.Bytes(), c); err != nil {
		log.Fatal(err)
	}

	return
}

func initLogger(config *orchestrator.Config) {
	config.Logger = model.NewLogger(os.Stdout, config.LogLevel.Level)
}

func mustConnectToNATS(config *orchestrator.Config) {
	var err error
	opts := nats.Options{
		User:     config.NatsUser,
		Password: config.NatsPass,
		Servers:  config.NatsHosts,
	}

	if config.Connection, err = opts.Connect(); err != nil {
		config.Logger.Fatal(err)
	}
}
