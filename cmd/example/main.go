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

	fmt.Println("Example application started, press Enter to quit . . .")
	fmt.Scanln()

	a.Stop()
}
