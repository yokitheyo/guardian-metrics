package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/yokitheyo/guardian-metrics/internal/agent"
)

func main() {
	var (
		addr              = flag.String("a", "localhost:8080", "address of HTTP server")
		reportIntervalStr = flag.String("r", "10", "report interval in seconds")
		pollIntervalStr   = flag.String("p", "2", "poll interval in seconds")
	)
	flag.Parse()

	if len(flag.Args()) > 0 {
		log.Fatalf("unknown flag(s): %v", flag.Args())
	}

	reportIntervalSec, err := strconv.Atoi(*reportIntervalStr)
	if err != nil || reportIntervalSec <= 0 {
		log.Fatalf("invalid report interval: %v", *reportIntervalStr)
	}
	pollIntervalSec, err := strconv.Atoi(*pollIntervalStr)
	if err != nil || pollIntervalSec <= 0 {
		log.Fatalf("invalid poll interval: %v", *pollIntervalStr)
	}

	collector := agent.NewRuntimeCollector()
	sender := agent.NewHTTPSender("http://" + *addr)

	a := agent.NewAgent(
		collector,
		sender,
		time.Duration(pollIntervalSec)*time.Second,
		time.Duration(reportIntervalSec)*time.Second,
		"http://"+*addr,
	)

	log.Println("Starting agent...")
	a.Run()
}
