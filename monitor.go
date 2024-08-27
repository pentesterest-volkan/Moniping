package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func MonitorServers(ctx context.Context, config *Config) {
	rateLimiter := time.Tick(time.Second / time.Duration(config.RateLimit))
	workerCount := config.WorkerCount

	semaphore := make(chan struct{}, workerCount)

	for _, server := range config.Servers {
		semaphore <- struct{}{}
		go func(srv Server) {
			defer func() { <-semaphore }()
			for {
				select {
				case <-ctx.Done():
					logrus.Infof("Stopping monitoring for server: %s (%s)", srv.Name, srv.IP)
					return
				case <-rateLimiter:
					logrus.Infof("Pinging server: %s (%s)", srv.Name, srv.IP)
					up := PingServer(&srv)
					srv.mutex.Lock()
					if !up {
						srv.FailCount++
						if srv.FailCount >= srv.MaxRetries && !srv.WasDown {
							SendAlert(config, &srv, "down")
							srv.WasDown = true
						}
					} else {
						if srv.WasDown {
							SendAlert(config, &srv, "up")
						}
						srv.FailCount = 0 // Reset fail count on success
						srv.WasDown = false
					}
					srv.mutex.Unlock()
				}
			}
		}(server)
	}
}
