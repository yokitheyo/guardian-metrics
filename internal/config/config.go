package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

type AgentConfig struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

type ServerConfig struct {
	Address string
}

func LoadAgentConfig() *AgentConfig {
	conf := &AgentConfig{}
	var reportInterval, pollInterval int

	flag.StringVar(&conf.Address, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()

	conf.ReportInterval = time.Duration(reportInterval) * time.Second
	conf.PollInterval = time.Duration(pollInterval) * time.Second

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		conf.Address = envAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if v, err := strconv.Atoi(envReportInterval); err == nil {
			conf.ReportInterval = time.Duration(v) * time.Second
		} else {
			log.Printf("invalid REPORT_INTERVAL: %v, using default", err)
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if v, err := strconv.Atoi(envPollInterval); err == nil {
			conf.PollInterval = time.Duration(v) * time.Second
		} else {
			log.Printf("invalid POLL_INTERVAL: %v, using default", err)
		}
	}

	return conf
}

func LoadServerConfig() *ServerConfig {
	conf := &ServerConfig{}

	flag.StringVar(&conf.Address, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		conf.Address = envAddr
	}

	return conf
}
