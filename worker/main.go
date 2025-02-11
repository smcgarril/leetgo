package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/smcgarril/leetgo-worker/api"
)

func main() {
	// Create router
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/process-code", api.ProcessCodeHandler).Methods("POST")

	// Enable CORS for all origins (for development purposes)
	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)

	fmt.Println("Worker service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", corsHandler))
}
