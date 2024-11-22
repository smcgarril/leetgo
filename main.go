package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Fetch all problems
func getProblems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Problems)
}

func executeCode(w http.ResponseWriter, r *http.Request) {
	// Decode the user-submitted code from the request body
	var codeRequest CodeRequest

	// Read the body into a byte slice
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read request body"}`, http.StatusInternalServerError)
		return
	}

	// Log the raw request body
	log.Printf("Raw request body: %s", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &codeRequest)
	if err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Log the decoded payload
	log.Printf("Decoded request: %+v", codeRequest)

	// Test Harness Template
	harness := `
	package main
	import (
		"fmt"
	)

	// User's function
	%s

	func main() {

		// Test input
		input := "%s"

		// Call the user's function and check the output
		output := %s(input) 
		expected := %s

		if fmt.Sprint(output) == fmt.Sprint(expected) {
			fmt.Println("PASSED")
		} else {
			fmt.Printf("FAILED: Got %%v, Expected %%v\n", output, expected)
		}
	}`

	// Wrap the submitted code in the test harness
	code := fmt.Sprintf(harness, codeRequest.Code, codeRequest.TestInput, codeRequest.Problem, codeRequest.Expected)

	// Save the submitted code to a temporary file
	codeFile := "temp_code.go"

	err = os.WriteFile(codeFile, []byte(code), 0644)
	if err != nil {
		http.Error(w, `{"error":"Failed to save code"}`, http.StatusInternalServerError)
		return
	}

	// Execute the Go code
	cmd := exec.Command("go", "run", codeFile)
	output, err := cmd.CombinedOutput()

	// Check for execution errors
	if err != nil {
		http.Error(w, `{"error":"Execution error", "details":"`+string(output)+`"}`, http.StatusInternalServerError)
		return
	}

	// Prepare the JSON response
	response := map[string]string{
		"output": string(output),
	}
	fmt.Println(response)

	if err != nil {
		// Include the error message in the response
		response["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Encode and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	// Delete temp_code.go
	err = os.Remove(codeFile)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete tmp file"}`, http.StatusInternalServerError)
		return
	}
}

func main() {
	// Create a new router using gorilla/mux
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/problems", getProblems).Methods("GET")
	router.HandleFunc("/execute", executeCode).Methods("POST")

	// Serve static frontend files (HTML, JS, CSS)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public"))))

	// Enable CORS for all origins (for development purposes)
	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)

	// Start the server on port 8080
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
