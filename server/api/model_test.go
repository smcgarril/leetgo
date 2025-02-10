package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAllProblems(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func(mock sqlmock.Sqlmock)
		expectedCode int
		expectedBody string
	}{
		{
			name: "SuccessfulFetch",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "short_description", "long_description", "problem_seed", "examples", "difficulty", "attempts", "solves"}).
					AddRow("1", "Problem 1", "Short Desc 1", "Long Desc 1", "Seed 1", "Example 1", "Easy", "10", "5").
					AddRow("2", "Problem 2", "Short Desc 2", "Long Desc 2", "Seed 2", "Example 2", "Medium", "20", "10")
				mock.ExpectQuery("^SELECT (.+) FROM problems$").WillReturnRows(rows)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"id":"1","name":"Problem 1","short_description":"Short Desc 1","long_description":"Long Desc 1","difficulty":"Easy","problem_seed":"Seed 1","examples":"Example 1","attempts":"10","solves":"5"},{"id":"2","name":"Problem 2","short_description":"Short Desc 2","long_description":"Long Desc 2","difficulty":"Medium","problem_seed":"Seed 2","examples":"Example 2","attempts":"20","solves":"10"}]`,
		},
		{
			name: "DatabaseQueryError",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM problems$").WillReturnError(errors.New("query error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "Error fetching problems from database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			tt.mockSetup(mock)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/problems", nil)

			GetAllProblems(db, rec, req)

			equals(t, rec.Code, tt.expectedCode)
			equals(t, strings.TrimSpace(rec.Body.String()), tt.expectedBody)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
