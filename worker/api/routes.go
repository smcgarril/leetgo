package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Handler for processing code submissions
func ProcessCodeHandler(w http.ResponseWriter, r *http.Request) {
	var submission CodeSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	codeResponse, err := processCode(submission)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(codeResponse); err != nil {
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// Process the code submission
func processCode(submission CodeSubmission) (CodeOutput, error) {
	log.Printf("Retrieved problem examples: %+v", submission.ProblemExamples)

	var testCalls []string
	for _, example := range submission.ProblemExamples {
		formattedArgs, err := prepareTestCall(example, submission.Problem)
		if err != nil {
			return CodeOutput{}, fmt.Errorf("failed to prepare test call for example ID %d: %w", example.ID, err)
		}
		testCalls = append(testCalls, formattedArgs)
	}

	harnessCode := generateTestHarness(submission.Code, strings.Join(testCalls, "\n"))

	codeFile := "temp_code.go"
	if err := os.WriteFile(codeFile, []byte(harnessCode), 0644); err != nil {
		return CodeOutput{}, fmt.Errorf("failed to save code to file: %w", err)
	}
	defer os.Remove(codeFile) // Ensure the temp file is removed

	testReturn, err := exec.Command("go", "run", codeFile).CombinedOutput()
	if err != nil {
		log.Printf("Error executing test harness: %v", err)
	}

	output := string(testReturn)
	testCount := len(submission.ProblemExamples)
	testPassed := CountPassingTests(output)
	result := "FAILED"
	if testCount == testPassed {
		result = "PASSED"
	}

	response := CodeOutput{
		TestCount:  testCount,
		TestPassed: testPassed,
		Output:     output,
		Result:     result,
	}

	log.Printf("Response: %+v", response)
	return response, nil
}

// Prepare the test call for the given example
func prepareTestCall(example ProblemExample, problemName string) (string, error) {
	var inputOrder []string
	if err := json.Unmarshal([]byte(example.InputOrder), &inputOrder); err != nil {
		return "", fmt.Errorf("failed to unmarshal input order: %w", err)
	}

	formattedArgs, err := FormatArgs(example.Input, inputOrder)
	if err != nil {
		return "", fmt.Errorf("failed to format arguments: %w", err)
	}

	expectedOutput, err := FormatExpectedOutput(example.ExpectedOutput)
	if err != nil {
		return "", fmt.Errorf("failed to format expected output: %w", err)
	}

	return fmt.Sprintf(`
		output%d := %s(%s)
		expected%d := %s
		if fmt.Sprint(output%d) == fmt.Sprint(expected%d) {
			results = append(results, Result{%d, "PASSED", ""})
		} else {
			results = append(results, Result{%d, "FAILED", fmt.Sprint(output%d)})
		}
	`, example.ID, problemName, formattedArgs, example.ID, expectedOutput, example.ID, example.ID, example.ID, example.ID, example.ID), nil
}

// Generate the test harness code
func generateTestHarness(userCode, testCalls string) string {
	return fmt.Sprintf(`
		package main
		import (
			"fmt"
		)
		type Result struct {
			Test   int
			Result string
			Output string
		}

		%s

		func main() {
			var results []Result
			%s
			for _, result := range results {
				fmt.Printf("Test %%d: %%s, Output: %%s\n", result.Test, result.Result, result.Output)
			}
		}
	`, userCode, testCalls)
}
