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

func processCode(submission CodeSubmission) (CodeResponse, error) {
	// Log the retrieved examples
	log.Printf("Retrieved problem examples: %+v", submission.ProblemExamples)

	// Initialize slice of test calls
	var testCalls []string

	for i, example := range submission.ProblemExamples {
		inputJSON := example.Input
		input_order := []string{}
		err := json.Unmarshal([]byte(example.InputOrder), &input_order)
		if err != nil {
			fmt.Println("Error unmarshalling input_order:", err)
			return CodeResponse{}, err
		}

		formattedArgs, err := FormatArgs(inputJSON, input_order)
		if err != nil {
			fmt.Println("Error formatting arguments:", err)
			return CodeResponse{}, err
		}

		fmt.Printf("The formattedArgs are: %s", formattedArgs)

		outputJSON := example.ExpectedOutput
		expectedOutput, err := FormatExpectedOutput(outputJSON)
		if err != nil {
			fmt.Println("Error formatting expected output:", err)
			return CodeResponse{}, err
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
		`, i+1, submission.Problem, formattedArgs, i+1, expectedOutput, i+1, i+1))

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
	harnessCode := fmt.Sprintf(harness, submission.Code, strings.Join(testCalls, "\n"))

	// Save the generated code to a temporary file
	codeFile := "temp_code.go"

	err := os.WriteFile(codeFile, []byte(harnessCode), 0644)
	if err != nil {
		fmt.Println("Failed to save code:", err)
		return CodeResponse{}, err
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
		fmt.Printf("Error deleting tmp file: ", err)
		return CodeResponse{}, err
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

	return response, nil
}
