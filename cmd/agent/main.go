package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/yokitheyo/guardian-metrics/internal/agent"
)

func main() {
	envAddr := os.Getenv("ADDRESS")
	if envAddr == "" {
		envAddr = "localhost:8080"
	}
	envReportInterval := os.Getenv("REPORT_INTERVAL")
	if envReportInterval == "" {
		envReportInterval = "10"
	}
	envPollInterval := os.Getenv("POLL_INTERVAL")
	if envPollInterval == "" {
		envPollInterval = "2"
	}

	addr := flag.String("a", envAddr, "address of HTTP server")
	reportIntervalStr := flag.String("r", envReportInterval, "report interval in seconds")
	pollIntervalStr := flag.String("p", envPollInterval, "poll interval in seconds")
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
