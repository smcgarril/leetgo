package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mocks

type TestStruct struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func mockGetProblemExamples(db *sql.DB, problemID string) ([]ProblemExample, error) {
	if problemID == "1" {
		return []ProblemExample{{Input: "test input", ExpectedOutput: "test output"}}, nil
	}
	return nil, errors.New("database error")
}

func mockCallWorkerService(codeSubmission CodeSubmission) (CodeOutput, error) {
	if codeSubmission.Code == "fail" {
		return CodeOutput{}, errors.New("worker service error")
	}
	return CodeOutput{Result: "PASSED"}, nil
}

// Tests

func TestExecuteCode(t *testing.T) {
	// Save original functions and restore them at the end
	originalGetProblemExamples := GetProblemExamplesWrapper
	originalCallWorkerService := callWorkerServiceWrapper

	// Mock functions
	GetProblemExamplesWrapper = mockGetProblemExamples
	callWorkerServiceWrapper = mockCallWorkerService

	defer func() {
		GetProblemExamplesWrapper = originalGetProblemExamples
		callWorkerServiceWrapper = originalCallWorkerService
	}()

	tests := []struct {
		name               string
		input              CodeSubmission
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "SuccessfulExecution",
			input: CodeSubmission{
				ProblemID: "1",
				Code:      "valid code",
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "PASSED",
		},
		{
			name: "DatabaseError",
			input: CodeSubmission{
				ProblemID: "999", // ProblemID that triggers database error
				Code:      "valid code",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Failed to retrieve problem examples",
		},
		{
			name: "WorkerServiceError",
			input: CodeSubmission{
				ProblemID: "1",
				Code:      "fail", // Code that triggers worker service error
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "Failed to execute code",
		},
		{
			name:               "InvalidRequest",
			input:              CodeSubmission{}, // Missing fields in the request
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/execute", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			ExecuteCode(nil, rec, req)

			equals(t, tt.expectedStatusCode, rec.Code)

			if !bytes.Contains(rec.Body.Bytes(), []byte(tt.expectedResponse)) {
				t.Errorf("Expected response to contain %q, got %q", tt.expectedResponse, rec.Body.String())
			}
		})
	}
}

func TestDecodeRequest(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expected      TestStruct
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid JSON",
			requestBody: `{"name":"John Doe", "email":"john@example.com"}`,
			expected:    TestStruct{Name: "John Doe", Email: "john@example.com"},
			expectError: false,
		},
		{
			name:          "Invalid JSON",
			requestBody:   `{"name":"John Doe", "email":}`, // Invalid JSON
			expectError:   true,
			errorContains: "invalid character",
		},
		{
			name:          "Empty Body",
			requestBody:   ``,
			expectError:   true,
			errorContains: "unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")

			var result TestStruct
			err := decodeRequest(req, &result)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got nil")
				} else if tt.errorContains != "" && !bytes.Contains([]byte(err.Error()), []byte(tt.errorContains)) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %+v, got %+v", tt.expected, result)
				}
			}
		})
	}
}

func TestBuildCodeOutput(t *testing.T) {
	tests := []struct {
		name           string
		codeOutput     CodeOutput
		examples       []ProblemExample
		expectedResult CodeOutput
	}{
		{
			name: "SuccessfulTest",
			codeOutput: CodeOutput{
				TestCount:  1,
				TestPassed: 1,
				Output:     "2",
				Result:     "PASSED",
			},
			examples: []ProblemExample{
				{ID: 1, Input: "1", ExpectedOutput: "2"},
			},
			expectedResult: CodeOutput{
				TestCount:  1,
				TestPassed: 1,
				Output:     "2",
				Input:      "1",
				Expected:   "2",
				Result:     "PASSED",
			},
		},
		{
			name: "FailedTest",
			codeOutput: CodeOutput{
				TestCount:  1,
				TestPassed: 0,
				Output:     "-1",
				Result:     "FAILED",
			},
			examples: []ProblemExample{
				{ID: 1, Input: "1", ExpectedOutput: "2"},
			},
			expectedResult: CodeOutput{
				TestCount:  1,
				TestPassed: 0,
				Output:     "-1",
				Input:      "1",
				Expected:   "2",
				Result:     "FAILED",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildCodeOutput(tt.codeOutput, tt.examples)

			equals(t, tt.expectedResult.TestCount, result.TestCount)
			equals(t, tt.expectedResult.TestPassed, result.TestPassed)
			equals(t, tt.expectedResult.Result, result.Result)
		})
	}
}
