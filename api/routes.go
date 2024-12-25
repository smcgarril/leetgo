package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/traefik/yaegi/interp"
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
		SELECT id, problem_id, input, expected_output 
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

// parseGoExpression takes a Go expression in string form and evaluates it into an actual Go value.
func parseGoExpression(expr string) (interface{}, error) {
	// Use the Go parser to validate and parse the expression
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %w", err)
	}

	// Handle specific cases based on the expression type
	switch v := node.(type) {
	case *ast.CompositeLit: // Handle arrays/slices
		fmt.Println("In the *ast.CompositeLit")
		return parseCompositeLit(v)
	case *ast.BasicLit: // Handle basic literals
		fmt.Println("In the *ast.BasicLit")
		return parseBasicLit(v)
	case *ast.UnaryExpr: // Handle unary expressions (e.g., -5, +0)
		return parseUnaryExpr(v)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", v)
	}
}

// parseCompositeLit parses a Go composite literal (e.g., `[]int{1, 2, 3}`)
func parseCompositeLit(lit *ast.CompositeLit) (interface{}, error) {
	// Example assumes array literals of type []int
	var result []int
	for _, elt := range lit.Elts {
		if basicLit, ok := elt.(*ast.BasicLit); ok {
			var value int
			fmt.Sscanf(basicLit.Value, "%d", &value) // Basic conversion from string
			result = append(result, value)
		} else {
			return nil, fmt.Errorf("unsupported composite element type: %T", elt)
		}
	}
	return result, nil
}

// parseBasicLit parses a Go basic literal (e.g., "5")
func parseBasicLit(lit *ast.BasicLit) (interface{}, error) {
	switch lit.Kind {
	case token.INT: // Integer
		fmt.Println("In the token.INT")
		var value int
		fmt.Sscanf(lit.Value, "%d", &value)
		return value, nil
	case token.STRING: // String
		fmt.Println("In the token.STRING")
		return lit.Value[1 : len(lit.Value)-1], nil // Remove quotes
	default:
		return nil, fmt.Errorf("unsupported literal kind: %s", lit.Kind)
	}
}

// parseUnaryExpr parses a Go unary expression (e.g., -5, +0)
func parseUnaryExpr(expr *ast.UnaryExpr) (interface{}, error) {
	// Only handle basic literals as the operand
	if basicLit, ok := expr.X.(*ast.BasicLit); ok {
		value, err := parseBasicLit(basicLit)
		if err != nil {
			return nil, err
		}

		// Apply the unary operator
		switch expr.Op {
		case token.SUB: // Negative numbers
			if intValue, ok := value.(int); ok {
				return -intValue, nil
			}
		case token.ADD: // Positive numbers (no-op)
			return value, nil
		default:
			return nil, fmt.Errorf("unsupported unary operator: %s", expr.Op)
		}
	}

	return nil, fmt.Errorf("unsupported operand type for unary expression: %T", expr.X)
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

	// Initialize slice of strings for results
	results := []string{}

	testHarness := `
		package main

		%s
	`
	userCode := fmt.Sprintf(testHarness, codeSubmission.Code)
	println(userCode)

	i := interp.New(interp.Options{})

	_, err = i.Eval(userCode)
	if err != nil {
		panic(err)
	}

	testFunc := fmt.Sprintf("main.%s", codeSubmission.Problem)
	println(testFunc)
	v, err := i.Eval(testFunc)
	if err != nil {
		panic(err)
	}

	sum := v.Interface().(func(int, int) int)

	for _, example := range examples {
		inputJSON := example.Input
		var data map[string]interface{}
		err := json.Unmarshal([]byte(inputJSON), &data)
		if err != nil {
			panic(err)
		}

		// Iterate over the map and parse values
		parsedData := make(map[string]interface{})
		for key, value := range data {
			// Assuming all values are strings needing parsing
			if strValue, ok := value.(string); ok {
				fmt.Println("HEEEEREEEE")
				parsedValue, err := parseGoExpression(strValue)
				if err != nil {
					fmt.Printf("Failed to parse key %s: %v\n", key, err)
					continue
				}
				parsedData[key] = parsedValue
			}
		}

		// Extract the values for x and y
		x, okX := parsedData["x"].(int)
		y, okY := parsedData["y"].(int)

		if !okX || !okY {
			panic("x or y is not an integer")
		}

		// Call the function dynamically
		output := sum(x, y)

		outputJSON := example.ExpectedOutput
		expectedOutput := FormatTestJSON(outputJSON)

		fmt.Println(output)
		fmt.Println(expectedOutput)

		println(fmt.Sprint(output) == fmt.Sprint(expectedOutput))

		if fmt.Sprint(output) == expectedOutput {
			results = append(results, "PASSED")
		} else {
			results = append(results, "FAILED")
		}
	}

	// for _, example := range examples {

	// 	inputJSON := example.Input
	// 	formattedArgs := FormatTestJSON(inputJSON)

	// 	fmt.Printf("The formattedArgs are: %s", formattedArgs)

	// 	outputJSON := example.ExpectedOutput
	// 	expectedOutput := FormatTestJSON(outputJSON)

	// 	fmt.Printf("The expectedOutput are: %s", expectedOutput)

	// 	// Generate the test harness
	// 	harness := `
	// 	package main
	// 	import (
	// 		"fmt"
	// 	)

	// 	// User's function
	// 	%s

	// 	func main() {
	// 		// Call the user's function with deconstructed inputs
	// 		output := %s(%s)
	// 		expected := %#v

	// 		if fmt.Sprint(output) == fmt.Sprint(expected) {
	// 			fmt.Println("PASSED")
	// 		} else {
	// 			fmt.Printf("FAILED")
	// 		}
	// 	}`

	// 	// Wrap the submitted code in the test harness
	// 	code := fmt.Sprintf(harness, codeSubmission.Code, codeSubmission.Problem, formattedArgs, expectedOutput)

	// 	// Save the submitted code to a temporary file
	// 	codeFile := "temp_code.go"

	// 	err = os.WriteFile(codeFile, []byte(code), 0644)
	// 	if err != nil {
	// 		http.Error(w, `{"error":"Failed to save code"}`, http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Execute the Go code
	// 	cmd := exec.Command("go", "run", codeFile)

	// 	testReturn, err := cmd.CombinedOutput()
	// 	fmt.Printf("Test return: %s", testReturn)

	// 	// Check for execution errors
	// 	if err != nil {
	// 		http.Error(w, `{"error":"Execution error", "details":"`+string(testReturn)+`"}`, http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Delete temp_code.go
	// 	err = os.Remove(codeFile)
	// 	if err != nil {
	// 		http.Error(w, `{"error":"Failed to delete tmp file"}`, http.StatusInternalServerError)
	// 		return
	// 	}

	// 	results = append(results, string(testReturn))
	// }

	fmt.Printf("results: %v", results)

	testCount := len(results)
	testPassed := 0
	for _, t := range results {
		normalized := strings.TrimSpace(t)
		if normalized == "PASSED" {
			testPassed++
		}
	}
	didPass := "PASSED"
	if testCount != testPassed {
		didPass = "FAILED"
	}

	// Sanity logging:
	fmt.Printf("Number of tests: %d\n", testCount)
	fmt.Printf("Number passed: %d\n", testPassed)
	fmt.Printf("Results: %s\n", didPass)

	response := CodeResponse{
		TestCount:  testCount,
		TestPassed: testPassed,
		Output:     didPass,
	}
	fmt.Println(response)

	// if err != nil {
	// 	// Include the error message in the response
	// 	response["error"] = err.Error()
	// 	w.WriteHeader(http.StatusInternalServerError)
	// }

	// Encode and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
