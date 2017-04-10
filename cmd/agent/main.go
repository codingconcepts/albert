package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/codingconcepts/albert/pkg/agent"
	"github.com/codingconcepts/albert/pkg/model"
	nats "github.com/nats-io/go-nats"
)

func main() {
	config := mustLoadConfig("config.json")
	initLogger(config)
	mustConnectToNATS(config)

	a := agent.NewAgent(config)
	go a.Start()

	config.Logger.Info("agent started successfully")
	fmt.Scanln()
}

func mustLoadConfig(path string) (c *agent.Config) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	buffer := new(bytes.Buffer)
	if _, err = io.Copy(buffer, file); err != nil {
		log.Fatal(err)
	}

	c = new(agent.Config)
	if err = json.Unmarshal(buffer.Bytes(), c); err != nil {
		log.Fatal(err)
	}

	return
}

func initLogger(config *agent.Config) {
	config.Logger = model.NewLogger(os.Stdout, config.LogLevel.Level)
}

func mustConnectToNATS(config *agent.Config) {
	opts := nats.Options{
		User:     config.NatsUser,
		Password: config.NatsPass,
		Servers:  config.NatsHosts,
	}

	var err error
	if config.Connection, err = opts.Connect(); err != nil {
		config.Logger.Fatal(err)
	}
}
