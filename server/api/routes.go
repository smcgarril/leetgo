package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Define a function that returns an http.HandlerFunc
func GetProblemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetProblems(db, w, r)
	}
}

// Define a function that returns an http.HandlerFunc
func ExecuteCodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ExecuteCode(db, w, r)
	}
}

func ExecuteCode(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var codeSubmission CodeSubmission

	// Read the body into a byte slice
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read request body"}`, http.StatusInternalServerError)
		return
	}

	// Unmarshal JSON into codeSubmission struct
	err = json.Unmarshal(bodyBytes, &codeSubmission)
	if err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Log the decoded payload
	log.Printf("User submission: %+v", codeSubmission)

	// Retrieve problem examples from the database
	examples, err := GetProblemExamples(db, codeSubmission.ProblemID)
	if err != nil {
		http.Error(w, `{"error":"Failed to retrieve problem examples"}`, http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	// Log the retrieved examples
	log.Printf("Retrieved problem examples: %+v", examples)

	// Attach problem examples to the code submission
	codeSubmission.ProblemExamples = examples

	// Marshal the updated CodeSubmission into JSON for the worker API
	workerRequestBody, err := json.Marshal(codeSubmission)
	if err != nil {
		http.Error(w, `{"error":"Failed to prepare worker request"}`, http.StatusInternalServerError)
		return
	}

	// Log the request payload for debugging
	log.Printf("Worker request body: %s", string(workerRequestBody))

	// Send the code submission to the worker service
	workerURL := GetWorkerURL()
	resp, err := http.Post(workerURL, "application/json", bytes.NewBuffer(workerRequestBody))
	if err != nil {
		http.Error(w, `{"error":"Failed to connect to worker service"}`, http.StatusInternalServerError)
		log.Printf("Worker service error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Read the worker service's response
	workerResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read worker response"}`, http.StatusInternalServerError)
		return
	}

	// Log the worker response for debugging
	log.Printf("Worker response: %s", string(workerResponseBody))

	// Forward the worker response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(workerResponseBody)
}
