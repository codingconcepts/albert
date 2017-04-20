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
	config := mustLoadConfig("config.json")
	logger := model.NewLogger(os.Stdout, config.LogLevel.Level)

	opts := nats.Options{
		User:     config.NatsUser,
		Password: config.NatsPass,
		Servers:  config.NatsHosts,
	}

	conn, err := opts.Connect()
	if err != nil {
		return
	}

	processor := newNatsProcessor(conn)
	killer := newCmdKiller()
	a, err := agent.NewAgent(config, processor, killer, logger)
	if err != nil {
		log.Fatal(err)
	}
	go a.Start()

	a.Logger.Info("agent started successfully")
	fmt.Scanln()

	a.Stop()
}

func mustLoadConfig(path string) (config *agent.Config) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	if config, err = agent.NewConfigFromReader(file); err != nil {
		panic(err)
	}

	return
}
