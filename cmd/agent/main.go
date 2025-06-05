package main

import (
	"log"
	"time"

	"github.com/yokitheyo/guardian-metrics/internal/agent"
)

func main() {
	collector := agent.NewRuntimeCollector()
	sender := agent.NewHTTPSender("http://localhost:8080")

	a := agent.NewAgent(
		collector,
		sender,
		2*time.Second,
		10*time.Second,
		"http://localhost:8080",
	)

	log.Println("Starting agent...")
	a.Run()
}
