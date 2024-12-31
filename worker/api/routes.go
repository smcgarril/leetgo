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

func ProcessCodeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var submission CodeSubmission
	err := json.NewDecoder(r.Body).Decode(&submission)
	if err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Generate and test the code
	codeResponse, err := processCode(submission)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Respond with the test results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(codeResponse)
}

func processCode(submission CodeSubmission) (CodeOutput, error) {
	// Log the retrieved examples
	log.Printf("Retrieved problem examples: %+v", submission.ProblemExamples)

	// Initialize slice of test calls
	var testCalls []string

	for _, example := range submission.ProblemExamples {
		inputJSON := example.Input
		input_order := []string{}
		err := json.Unmarshal([]byte(example.InputOrder), &input_order)
		if err != nil {
			fmt.Println("Error unmarshalling input_order:", err)
			return CodeOutput{}, err
		}

		formattedArgs, err := FormatArgs(inputJSON, input_order)
		if err != nil {
			fmt.Println("Error formatting arguments:", err)
			return CodeOutput{}, err
		}

		outputJSON := example.ExpectedOutput
		expectedOutput, err := FormatExpectedOutput(outputJSON)
		if err != nil {
			fmt.Println("Error formatting expected output:", err)
			return CodeOutput{}, err
		}

		testID := example.ID

		// Append a single test call to the list
		testCalls = append(testCalls, fmt.Sprintf(`
			output%d := %s(%s)
			expected%d := %s
			if fmt.Sprint(output%d) == fmt.Sprint(expected%d) {
				results = append(results, Result{%d, "PASSED", ""})
			} else {
				results = append(results, Result{%d, "FAILED", fmt.Sprint(output%d)})
			}
		`, testID, submission.Problem, formattedArgs, testID, expectedOutput, testID, testID, testID, testID, testID))

	}

	// Generate the complete test harness
	harness := `
		package main
		import (
			"fmt"
		)
		type Result struct {
			Test int
			Result string
			Output string
		}

		// User function
		%s

		func main() {
			var results []Result

			%s

			for _, result := range results {
				fmt.Printf("Test %%d: %%s, Output: %%s\n", result.Test, result.Result, result.Output)
			}
		}`

	// Combine all test calls into the harness
	harnessCode := fmt.Sprintf(harness, submission.Code, strings.Join(testCalls, "\n"))

	// Save the generated code to a temporary file
	codeFile := "temp_code.go"

	err := os.WriteFile(codeFile, []byte(harnessCode), 0644)
	if err != nil {
		fmt.Println("Failed to save code:", err)
		return CodeOutput{}, err
	}

	// Execute the Go code
	cmd := exec.Command("go", "run", codeFile)
	testReturn, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing test harness: %v\n", err)
	}

	// Delete temp_code.go
	err = os.Remove(codeFile)
	if err != nil {
		fmt.Println("Error deleting tmp file: ", err)
		return CodeOutput{}, err
	}

	// Process the output to count PASSED and FAILED
	output := string(testReturn)
	testCount := len(submission.ProblemExamples)
	testPassed := CountPassingTests(output)

	result := "FAILED"
	if testCount == testPassed {
		result = "PASSED"
	}

	// Prepare the response
	response := CodeOutput{
		TestCount:  testCount,
		TestPassed: testPassed,
		Output:     output,
		Result:     result,
	}

	// Print response for debugging
	fmt.Println(response)
	log.Printf("Response: %+v", response)

	return response, nil
}
