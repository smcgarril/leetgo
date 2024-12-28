package main

import (
	"fmt"
	"net/http"

	"github.com/smcgarril/leetgo-worker/api"
)

func main() {
	http.HandleFunc("/process-code", api.ProcessCodeHandler)
	fmt.Println("Worker service running on port 8081")
	http.ListenAndServe(":8081", nil)
}
