package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	// Load configuration
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logging
	initLogging()

	// Start monitoring service
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logrus.Info("Service is starting...")

	go MonitorServers(ctx, config)
	go StartHealthCheck()
	go WaitForShutdown(cancel)

	<-ctx.Done()
	logrus.Info("Service stopped")
}
