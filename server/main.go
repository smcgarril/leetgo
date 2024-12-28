package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/smcgarril/leetgo/api"
)

var db *sql.DB

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
	api.SeedFiles(db)

	// log.Printf("Database seeded!")
	log.Printf("Testing database query...")

	// Sanity check
	api.QueryProblems(db)

	// Create a new router using gorilla/mux
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/problems", api.GetProblemsHandler(db)).Methods("GET")
	router.HandleFunc("/execute", api.ExecuteCodeHandler(db)).Methods("POST")

	// Serve static frontend files (HTML, JS, CSS)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public"))))

	// Enable CORS for all origins (for development purposes)
	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)

	// Start the server on port 8080
	log.Printf("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
