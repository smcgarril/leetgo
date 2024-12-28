package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Define a function that returns an http.HandlerFunc
func GetProblemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetProblems(db, w, r)
	}
}

// Fetch all problems
func GetProblems(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query(`
    SELECT 
        id, 
        name, 
        short_description, 
        long_description, 
        problem_seed, 
        REPLACE(examples, '\\"', "'") AS examples, 
        difficulty, 
        attempts, 
        solves 
    FROM problems
`)
	if err != nil {
		http.Error(w, "Error fetching problems from database", http.StatusInternalServerError)
		log.Printf("Query error: %v\n", err)
		return
	}
	defer rows.Close()

	var problems []Problem

	for rows.Next() {
		var p Problem
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.ShortDescription,
			&p.LongDescription,
			&p.ProblemSeed,
			&p.Examples,
			&p.Difficulty,
			&p.Attempts,
			&p.Solves,
		)
		if err != nil {
			http.Error(w, "Error scanning problems from database", http.StatusInternalServerError)
			log.Printf("Row scan error: %v\n", err)
			return
		}
		problems = append(problems, p)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating over problems", http.StatusInternalServerError)
		log.Printf("Row iteration error: %v\n", err)
		return
	}

	json.NewEncoder(w).Encode(problems)
}

func GetProblemExamples(db *sql.DB, problemID string) ([]ProblemExample, error) {
	rows, err := db.Query(`
		SELECT id, problem_id, input, input_order, expected_output 
		FROM problem_examples 
		WHERE problem_id = ?`, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examples []ProblemExample

	for rows.Next() {
		var example ProblemExample
		err := rows.Scan(
			&example.ID,
			&example.PromblemID,
			&example.Input,
			&example.InputOrder,
			&example.ExpectedOutput,
		)
		if err != nil {
			return nil, err
		}
		examples = append(examples, example)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return examples, nil
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

	// Log the raw request body
	log.Printf("Raw request body: %s", string(bodyBytes))

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

	// Initialize slice of test calls
	var testCalls []string

	for i, example := range examples {
		inputJSON := example.Input
		input_order := []string{}
		err := json.Unmarshal([]byte(example.InputOrder), &input_order)
		if err != nil {
			fmt.Println("Error unmarshalling input_order:", err)
			return
		}

		formattedArgs, err := FormatArgs(inputJSON, input_order)
		if err != nil {
			fmt.Println("Error formatting arguments:", err)
			return
		}

		fmt.Printf("The formattedArgs are: %s", formattedArgs)

		outputJSON := example.ExpectedOutput
		expectedOutput, err := FormatExpectedOutput(outputJSON)
		if err != nil {
			fmt.Println("Error formatting expected output:", err)
			return
		}

		fmt.Printf("The expectedOutput are: %s", expectedOutput)

		// Append a single test call to the list
		testCalls = append(testCalls, fmt.Sprintf(`
			output%d := %s(%s)
			expected%d := %s
			if fmt.Sprint(output%d) == fmt.Sprint(expected%d) {
				results = append(results, "PASSED")
			} else {
				results = append(results, "FAILED")
			}
		`, i+1, codeSubmission.Problem, formattedArgs, i+1, expectedOutput, i+1, i+1))

	}

	// Generate the complete test harness
	harness := `
		package main
		import (
			"fmt"
		)

		// User function
		%s

		func main() {
			results := []string{}

			%s

			for i, result := range results {
				fmt.Printf("Test %%d: %%s\n", i+1, result)
			}
		}`

	// Combine all test calls into the harness
	harnessCode := fmt.Sprintf(harness, codeSubmission.Code, strings.Join(testCalls, "\n"))

	// Save the generated code to a temporary file
	codeFile := "temp_code.go"

	err = os.WriteFile(codeFile, []byte(harnessCode), 0644)
	if err != nil {
		fmt.Println("Failed to save code:", err)
		return
	}

	// Execute the Go code
	cmd := exec.Command("go", "run", codeFile)
	testReturn, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing test harness: %v\n", err)
	}
	fmt.Printf("Test return: %s\n", testReturn)

	// Delete temp_code.go
	err = os.Remove(codeFile)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete tmp file"}`, http.StatusInternalServerError)
		return
	}

	// Process the output to count PASSED and FAILED
	testCount := 0
	testPassed := 0
	output := string(testReturn)
	result := "PASSED"

	// Split the output into lines and count PASSED and FAILED results
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Test") {
			testCount++
			if strings.Contains(line, "PASSED") {
				testPassed++
			} else {
				result = "FAILED"
			}
		}
	}

	// Prepare the response
	response := CodeResponse{
		TestCount:  testCount,
		TestPassed: testPassed,
		Output:     output,
		Result:     result,
	}

	// Print response for debugging
	fmt.Println(response)

	// Encode and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
