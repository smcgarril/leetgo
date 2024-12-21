package api

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// Execute seed data from different files
func SeedFiles(db *sql.DB) {
	seedFiles := []string{
		"db/create_tables.sql",
		"db/seed_data.sql",
	}

	for _, file := range seedFiles {
		if err := ExecuteSQLFromFile(db, file); err != nil {
			log.Fatal("Error executing seed file: ", err)
		}
	}
}

// Execute SQL statements from a file
func ExecuteSQLFromFile(db *sql.DB, filename string) error {
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
func QueryProblems(db *sql.DB) error {
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
