package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codingconcepts/albert/pkg/model"
	"github.com/codingconcepts/albert/pkg/orchestrator"
	nats "github.com/nats-io/go-nats"
)

func main() {
	config, err := orchestrator.NewConfigFromFile("config.json")
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

	o := orchestrator.NewOrchestrator(config, conn, logger)
	go o.Start()
	defer o.Stop()

	for _, a := range o.Applications {
		logger.WithFields(logrus.Fields{
			"name":       a.Name,
			"schedule":   a.Schedule,
			"percentage": a.Percentage,
		}).Info("application schedule")
	}

	logger.Info("orchestrator started successfully")
	fmt.Scanln()
}
