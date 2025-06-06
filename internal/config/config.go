package config

import (
	"flag"
	"log"
	"os"
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

	flag.StringVar(&conf.Address, "a", "localhost:8080", "address and port to run server")
	flag.DurationVar(&conf.ReportInterval, "r", 10*time.Second, "report interval")
	flag.DurationVar(&conf.PollInterval, "p", 2*time.Second, "poll interval")

	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		conf.Address = envAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if d, err := time.ParseDuration(envReportInterval + "s"); err == nil {
			conf.ReportInterval = d
		} else {
			log.Printf("invalid REPORT_INTERVAL: %v, using default", err)
		}
	}

	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if d, err := time.ParseDuration(envPollInterval + "s"); err == nil {
			conf.PollInterval = d
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
