package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codingconcepts/albert/pkg/agent"
	"github.com/codingconcepts/albert/pkg/model"
	nats "github.com/nats-io/go-nats"
)

func main() {
	config, err := agent.NewConfigFromFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	logger := model.NewLogger(os.Stdout, config.LogLevel.Level)

	opts := nats.Options{
		User:     config.NatsUser,
		Password: config.NatsPass,
		Servers:  config.NatsHosts,
	}

	conn, err := opts.Connect()
	if err != nil {
		logger.Fatal(err)
	}

	a := agent.NewAgent(config, conn, logger)
	go a.Start()

	logger.Info("agent started successfully")
	fmt.Scanln()
}
