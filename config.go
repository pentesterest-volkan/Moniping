package main

import (
	"encoding/json"
	"os"
	"sync"
)

type Server struct {
	IP          string
	Name        string
	FailCount   int
	WasDown     bool
	MaxRetries  int
	Recipients  []string
	PingTimeout int
	PingCount   int
	mutex       sync.Mutex
}

type EmailConfig struct {
	From               string
	Password           string
	SMTPHost           string
	SMTPPort           string
	DefaultRecipients  []string
	InsecureSkipVerify bool // Option to disable TLS verification
}

type Config struct {
	ClientName  string
	Servers     []Server
	Email       EmailConfig
	WorkerCount int
	RateLimit   int
}

func LoadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	// Set default values
	for i := range config.Servers {
		if config.Servers[i].PingTimeout == 0 {
			config.Servers[i].PingTimeout = 5
		}
		if config.Servers[i].PingCount == 0 {
			config.Servers[i].PingCount = 3
		}
	}

	return &config, nil
}
