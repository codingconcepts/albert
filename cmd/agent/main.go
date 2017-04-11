package main

import (
	"fmt"
	"log"

	"github.com/codingconcepts/albert/pkg/agent"
)

func main() {
	a, err := agent.NewAgentFromConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	go a.Start()

	a.Logger.Info("agent started successfully")
	fmt.Scanln()

	a.Stop()
}
