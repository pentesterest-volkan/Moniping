package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func PingServer(server *Server) bool {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	success := false
	for i := 0; i < server.PingCount; i++ {
		cmd := exec.Command("ping", "-c", "1", "-W", fmt.Sprint(server.PingTimeout), server.IP)
		output, err := cmd.Output()
		if err == nil && strings.Contains(string(output), "1 packets received") {
			success = true
			break
		}
	}

	if success {
		server.FailCount = 0
	} else {
		server.FailCount++
	}

	return success
}
