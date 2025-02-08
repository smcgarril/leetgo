package api

import (
	"database/sql"
	"net/http"
)

func GetAllProblemsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetAllProblems(db, w, r)
	}
}

func GetProblemNamesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetProblemNames(db, w, r)
	}
}

func GetProblemDetailsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		GetProblemDetails(db, w, r)
	}
}

func ExecuteCodeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ExecuteCode(db, w, r)
	}
}
