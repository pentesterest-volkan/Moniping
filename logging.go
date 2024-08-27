package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logMutex       sync.Mutex
	lastLogged     = make(map[string]time.Time)
	logBuffer      []string
	logTicker      = time.NewTicker(10 * time.Second)
	maxBufferSize  = 1000 // Maximum size of the log buffer before forcing a flush
	adaptiveTicker = time.NewTicker(5 * time.Second)
)

func init() {
	go func() {
		for range logTicker.C {
			flushLogs()
		}
	}()

	go adaptiveFlush()
}

// Initialize logging to both stdout and file
func initLogging() {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	file, err := os.OpenFile("alerts_backup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Warn("Failed to log to file, using default stdout")
		return
	}
	log.SetOutput(file)
}

func LogAlert(server *Server, status string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	now := time.Now()

	// Check if the last log for this server was less than a minute ago
	if lastTime, exists := lastLogged[server.IP]; exists && now.Sub(lastTime) < time.Minute {
		return // Skip logging if it has been less than a minute
	}

	// Prepare the log entry as a string
	logEntry := fmt.Sprintf("Server %s (%s) is %s", server.Name, server.IP, status)

	// Add the log entry to the buffer
	logBuffer = append(logBuffer, logEntry)

	// Force a flush if buffer size exceeds maxBufferSize
	if len(logBuffer) >= maxBufferSize {
		flushLogs()
	}

	// Update the last logged time for this server
	lastLogged[server.IP] = now
}

func flushLogs() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if len(logBuffer) > 0 {
		// Log to stdout
		for _, entry := range logBuffer {
			logrus.Info(entry)
		}

		// Log to file
		file, err := os.OpenFile("alerts_backup.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logrus.Fatal(err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, entry := range logBuffer {
			writer.WriteString(entry + "\n")
		}
		writer.Flush()

		// Clear the buffer after flushing
		logBuffer = nil

		// Optionally trim the log file if necessary
		trimLogFile(file)
	}
}

func trimLogFile(file *os.File) {
	// Trim the log file to the last 1,000 lines
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > 1000 {
			lines = lines[1:]
		}
	}

	file.Truncate(0)
	file.Seek(0, 0)

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	writer.Flush()
}

func adaptiveFlush() {
	for range adaptiveTicker.C {
		logMutex.Lock()
		if len(logBuffer) > maxBufferSize/2 {
			flushLogs()
		}
		logMutex.Unlock()
	}
}
