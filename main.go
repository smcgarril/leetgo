package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Define a function that returns an http.HandlerFunc
func getProblemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getProblems(db, w, r)
	}
}

// Fetch all problems
func getProblems(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Query to fetch all problems
	rows, err := db.Query(`
		SELECT id, name, short_description, long_description, difficulty, attempts, solves 
		FROM problems
	`)
	if err != nil {
		http.Error(w, "Error fetching problems from database", http.StatusInternalServerError)
		log.Printf("Query error: %v\n", err)
		return
	}
	defer rows.Close()

	// Slice to hold the results
	var problems []Problem

	// Iterate over rows
	for rows.Next() {
		var p Problem
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.ShortDescription,
			&p.LongDescription,
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

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		http.Error(w, "Error iterating over problems", http.StatusInternalServerError)
		log.Printf("Row iteration error: %v\n", err)
		return
	}

	// Encode the result as JSON
	json.NewEncoder(w).Encode(problems)
}

func getProblemExamples(db *sql.DB, problemID string) ([]ProblemExample, error) {
	// Query the database for problem examples with the given problem_id
	rows, err := db.Query(`
		SELECT id, problem_id, input, expected_output 
		FROM problem_examples 
		WHERE problem_id = ?`, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var examples []ProblemExample

	// Iterate through the rows and populate the ProblemExample slice
	for rows.Next() {
		var example ProblemExample
		err := rows.Scan(&example.ID, &example.PromblemID, &example.Input, &example.ExpectedOutput)
		if err != nil {
			return nil, err
		}
		examples = append(examples, example)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return examples, nil
}

// convertToType dynamically converts a string to its Go type
func convertToType(input string) (interface{}, error) {
	// Trim surrounding whitespace
	input = strings.TrimSpace(input)

	// Try unmarshaling as JSON (handles arrays, maps, etc.)
	var result interface{}
	if err := json.Unmarshal([]byte(input), &result); err == nil {
		return result, nil
	}

	// If JSON unmarshaling fails, try parsing primitive types
	if i, err := strconv.Atoi(input); err == nil {
		return i, nil // Integer
	}

	if f, err := strconv.ParseFloat(input, 64); err == nil {
		return f, nil // Float
	}

	// Assume it's a string if no other parsing succeeded
	return input, nil
}

// Create a comma-separated list of arguments from the input
func formatFunctionArguments(input interface{}) (string, error) {
	slice, ok := input.([]interface{})
	if !ok {
		return "", fmt.Errorf("input is not a slice: %v", input)
	}

	args := []string{}
	for _, val := range slice {
		args = append(args, fmt.Sprintf("%#v", val))
	}
	return strings.Join(args, ", "), nil
}

func executeCode(w http.ResponseWriter, r *http.Request) {
	// Decode the user-submitted code from the request body
	var codeSubmission CodeSubmission

	// Read the body into a byte slice
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read request body"}`, http.StatusInternalServerError)
		return
	}

	// Log the raw request body
	log.Printf("Raw request body: %s", string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &codeSubmission)
	if err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	// Log the decoded payload
	log.Printf("User submission: %+v", codeSubmission)

	// Retrieve problem examples from the database
	examples, err := getProblemExamples(db, codeSubmission.ProblemID)
	if err != nil {
		http.Error(w, `{"error":"Failed to retrieve problem examples"}`, http.StatusInternalServerError)
		log.Printf("Database error: %v", err)
		return
	}

	// Log the retrieved examples
	log.Printf("Retrieved problem examples: %+v", examples)

	// Convert the input and expected output to proper Go types
	testInput, err := convertToType(examples[0].Input)
	if err != nil {
		http.Error(w, `{"error":"Failed to parse input"}`, http.StatusInternalServerError)
		return
	}

	expectedOutput, err := convertToType(examples[0].ExpectedOutput)
	if err != nil {
		http.Error(w, `{"error":"Failed to parse expected output"}`, http.StatusInternalServerError)
		return
	}

	// Format the input as arguments for the function call
	formattedArgs, err := formatFunctionArguments(testInput)
	if err != nil {
		http.Error(w, `{"error":"Failed to format function arguments"}`, http.StatusInternalServerError)
		return
	}

	// Generate the test harness
	harness := `
	package main
	import (
		"fmt"
	)
	
	// User's function
	%s
	
	func main() {
		// Call the user's function with deconstructed inputs
		output := %s(%s)
		expected := %#v
	
		if fmt.Sprint(output) == fmt.Sprint(expected) {
			fmt.Println("PASSED")
		} else {
			fmt.Printf("FAILED: Got %%v, Expected %%v\n", output, expected)
		}
	}`

	// Wrap the submitted code in the test harness
	code := fmt.Sprintf(harness, codeSubmission.Code, codeSubmission.Problem, formattedArgs, expectedOutput)

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

// Execute seed data from different files
func seedFiles(db *sql.DB) {
	seedFiles := []string{
		"db/create_tables.sql", // Create tables
		"db/seed_data.sql",     // Seed data
	}

	for _, file := range seedFiles {
		if err := executeSQLFromFile(db, file); err != nil {
			log.Fatal("Error executing seed file: ", err)
		}
	}
}

// Execute SQL statements from a file
func executeSQLFromFile(db *sql.DB, filename string) error {
	// Read the SQL file
	sqlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	// Split the file contents by semicolons to get individual SQL statements
	statements := strings.Split(string(sqlData), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		// Execute the SQL statement
		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing statement from file %s: %v", filename, err)
		}
	}
	log.Printf("Seed data from %s executed successfully.\n", filename)
	return nil
}

// Retrieve all problems from DB for testing
func queryProblems(db *sql.DB) error {
	rows, err := db.Query("SELECT id, name, short_description, long_description, difficulty FROM problems")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, short_description, long_description, difficulty string
		if err := rows.Scan(&id, &name, &short_description, &long_description, &difficulty); err != nil {
			return err
		}
		log.Printf("Problem number: %d\nName: %s\nShort Description: %s\nLong Description: %s\nDifficulty: %s\n\n", id, name, short_description, long_description, difficulty)
	}

	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	var err error

	// Open SQLite database (it will be created if it doesn't exist)
	db, err = sql.Open("sqlite3", "./db/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully connected to SQLite database!")
	log.Printf("Seeding database file...")

	// Seed database with basic info for testing
	seedFiles(db)

	log.Printf("Database seeded!")
	log.Printf("Testing database query...")

	// Sanity check
	queryProblems(db)

	// Create a new router using gorilla/mux
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/problems", getProblemsHandler(db)).Methods("GET")
	router.HandleFunc("/execute", executeCode).Methods("POST")

	// Serve static frontend files (HTML, JS, CSS)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public"))))

	// Enable CORS for all origins (for development purposes)
	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)

	// Start the server on port 8080
	log.Printf("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
