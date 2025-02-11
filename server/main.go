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

	// Connect to SQLite database
	db, err = sql.Open("sqlite3", "./db/db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Seed database and test query
	log.Printf("Successfully connected to SQLite database!")
	log.Printf("Seeding database file...")
	api.SeedFiles(db)

	log.Printf("Testing database query...")
	api.QueryProblems(db)

	// Create router
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/problems", api.GetAllProblemsHandler(db)).Methods("GET")
	router.HandleFunc("/problems/names", api.GetProblemNamesHandler(db)).Methods("GET")
	router.HandleFunc("/problems/{id}", api.GetProblemDetailsHandler(db)).Methods("GET")
	router.HandleFunc("/execute", api.ExecuteCodeHandler(db)).Methods("POST")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public"))))

	// Enable CORS for all origins (for development purposes)
	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)

	log.Printf("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
