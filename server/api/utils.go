package api

import (
	"fmt"
	"os"
)

// Retrieve worker service URL:PORT from env variables
func GetWorkerURL() string {
	workerHost := os.Getenv("WORKER_HOST")
	if workerHost == "" {
		workerHost = "http://localhost"
	}

	workerPort := os.Getenv("WORKER_PORT")
	if workerPort == "" {
		workerPort = "8081"
	}

	workerPath := os.Getenv("WORKER_PATH")
	if workerPath == "" {
		workerPath = "/process-code"
	}

	return fmt.Sprintf("%s:%s%s", workerHost, workerPort, workerPath)
}
