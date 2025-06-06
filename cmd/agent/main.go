package main

import (
	"log"

	"github.com/yokitheyo/guardian-metrics/internal/agent"
	"github.com/yokitheyo/guardian-metrics/internal/config"
)

func main() {
	cfg := config.LoadAgentConfig()

	if cfg.Address == "" {
		log.Fatal("address is not set")
	}

	collector := agent.NewRuntimeCollector()
	sender := agent.NewHTTPSender("http://" + cfg.Address)

	a := agent.NewAgent(
		collector,
		sender,
		cfg.PollInterval,
		cfg.ReportInterval,
		"http://"+cfg.Address,
	)

	log.Println("Starting agent...")
	a.Run()
}
