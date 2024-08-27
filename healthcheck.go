package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func StartHealthCheck() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		logrus.Info("Health check OK")
	})
	logrus.Info("Starting health check server on :8080")
	http.ListenAndServe(":8080", nil)
}
