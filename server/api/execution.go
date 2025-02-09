package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

const workerTimeout = 10 * time.Second

// Handle a code execution request
func ExecuteCode(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var codeSubmission CodeSubmission

	if err := decodeRequest(r, &codeSubmission); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		log.Printf("Request decoding error: %v", err)
		return
	}

	log.Printf("User submission: %+v", codeSubmission)

	examples, err := GetProblemExamplesWrapper(db, codeSubmission.ProblemID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve problem examples")
		log.Printf("Database error: %v", err)
		return
	}

	log.Printf("Retrieved problem examples: %+v", examples)
	codeSubmission.ProblemExamples = examples

	codeOutput, err := callWorkerServiceWrapper(codeSubmission)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to execute code")
		log.Printf("Worker service error: %v", err)
		return
	}

	log.Printf("Worker response: %+v", codeOutput)

	response := buildCodeOutput(codeOutput, codeSubmission.ProblemExamples)
	respondWithJSON(w, http.StatusOK, response)
}

// Decode a request body into a struct
func decodeRequest(r *http.Request, v interface{}) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bodyBytes, v)
}

// Wrapper function to get problem examples
var callWorkerServiceWrapper func(codeSubmission CodeSubmission) (CodeOutput, error) = callWorkerService

// Send a code submission to the worker service
func callWorkerService(codeSubmission CodeSubmission) (CodeOutput, error) {
	var codeOutput CodeOutput

	workerRequestBody, err := json.Marshal(codeSubmission)
	if err != nil {
		return codeOutput, err
	}

	log.Printf("Worker request body: %s", string(workerRequestBody))

	ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)
	defer cancel()

	workerURL := GetWorkerURL()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, workerURL, bytes.NewBuffer(workerRequestBody))
	if err != nil {
		return codeOutput, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return codeOutput, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return codeOutput, err
	}

	workerResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return codeOutput, err
	}

	err = json.Unmarshal(workerResponseBody, &codeOutput)
	return codeOutput, err
}

// Build a response string from the code output and problem examples
func buildCodeOutput(codeOutput CodeOutput, examples []ProblemExample) CodeOutput {
	input, expectedOutput, actualOutput := BuildResponse(&codeOutput, examples)
	return CodeOutput{
		TestCount:  codeOutput.TestCount,
		TestPassed: codeOutput.TestPassed,
		Output:     actualOutput,
		Input:      input,
		Expected:   expectedOutput,
		Result:     codeOutput.Result,
	}
}

// Closure to build a response string from the code output and problem examples
func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}

// Send a JSON response
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}
