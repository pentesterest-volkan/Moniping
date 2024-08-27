package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func WaitForShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	logrus.Info("Shutdown signal received, flushing logs and terminating...")
	flushLogs() // Flush logs before shutting down
	cancel()
}
