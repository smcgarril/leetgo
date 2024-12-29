package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

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