package main

import (
	"log"

	"github.com/yokitheyo/guardian-metrics/internal/agent"
	"github.com/yokitheyo/guardian-metrics/internal/agent/collector"
	"github.com/yokitheyo/guardian-metrics/internal/agent/sender"
	"github.com/yokitheyo/guardian-metrics/internal/config"
)

func main() {
	cfg := config.LoadAgentConfig()

	if cfg.Address == "" {
		log.Fatal("address is not set")
	}

	coll := collector.NewRuntimeCollector()
	snd := sender.NewHTTPSender("http://" + cfg.Address)

	a := agent.NewAgent(
		coll,
		snd,
		cfg.PollInterval,
		cfg.ReportInterval,
		"http://"+cfg.Address,
	)

	log.Println("Starting agent...")
	a.Run()
}
