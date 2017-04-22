package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/model"
	"github.com/codingconcepts/albert/pkg/orchestrator"
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

	natsProcessor := newNatsProcessor(conn, config)
	o, err := orchestrator.NewOrchestrator(config, natsProcessor, logger)
	if err != nil {
		logger.Fatal(err)
	}

	go o.Start()
	defer o.Stop()

	for _, a := range config.Applications {
		logger.WithFields(logrus.Fields{
			"name":       a.Name,
			"schedule":   a.Schedule,
			"percentage": a.Percentage,
		}).Info("application schedule")
	}

	logger.Info("orchestrator started successfully")
	fmt.Scanln()
}

// mustLoadConfig loads a configuration file from disk and
// fails if it cannot be loaded (as the application won't be
// able to continue without it).
func mustLoadConfig(path string) (config *orchestrator.Config) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	if config, err = orchestrator.NewConfigFromReader(file); err != nil {
		panic(err)
	}

	return
}
