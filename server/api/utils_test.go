package api

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Helper functions for assertions

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		failTest(tb, msg, v...)
	}
}

// ok fails the test if an error is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		failTest(tb, "unexpected error: %s", err.Error())
	}
}

// equals fails the test if expected is not equal to actual.
func equals(tb testing.TB, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		failTest(tb, "\n\texpected: %#v\n\tgot: %#v", expected, actual)
	}
}

func failTest(tb testing.TB, msg string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	tb.FailNow()
}

// Tests

func TestGetWorkerURL(t *testing.T) {
	// Save original environment variables and restore them at the end
	originalWorkerHost := os.Getenv("WORKER_HOST")
	originalWorkerPort := os.Getenv("WORKER_PORT")
	originalWorkerPath := os.Getenv("WORKER_PATH")
	defer func() {
		os.Setenv("WORKER_HOST", originalWorkerHost)
		os.Setenv("WORKER_PORT", originalWorkerPort)
		os.Setenv("WORKER_PATH", originalWorkerPath)
	}()

	tests := []struct {
		name           string
		envVars        map[string]string
		expectedResult string
	}{
		{"NoEnvVars", map[string]string{}, "http://localhost:8081/process-code"},
		{"OnlyWorkerHostSet", map[string]string{"WORKER_HOST": "http://example.com"}, "http://example.com:8081/process-code"},
		{"OnlyWorkerPortSet", map[string]string{"WORKER_PORT": "9090"}, "http://localhost:9090/process-code"},
		{"OnlyWorkerPathSet", map[string]string{"WORKER_PATH": "/custom-path"}, "http://localhost:8081/custom-path"},
		{"AllEnvVarsSet", map[string]string{"WORKER_HOST": "http://example.com", "WORKER_PORT": "9090", "WORKER_PATH": "/custom-path"}, "http://example.com:9090/custom-path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			equals(t, tt.expectedResult, GetWorkerURL())

			for key := range tt.envVars {
				os.Unsetenv(key)
			}
		})
	}
}

func TestBuildResponse(t *testing.T) {
	codeOutput := &CodeOutput{
		Result: "FAILED",
		Output: "Test 1: FAILED, Output: -1",
	}
	examples := []ProblemExample{
		{ID: 1, Input: "1", ExpectedOutput: "2"},
		{ID: 2, Input: "2", ExpectedOutput: "4"},
	}

	tests := []struct {
		name              string
		codeResult        string
		expectedInput     string
		expectedOutput    string
		expectedActualOut string
	}{
		{"FailedTest", "FAILED", "1", "2", "-1"},
		{"NoFailedTest", "SUCCESS", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeOutput.Result = tt.codeResult
			input, expected, actual := BuildResponse(codeOutput, examples)
			equals(t, tt.expectedInput, input)
			equals(t, tt.expectedOutput, expected)
			equals(t, tt.expectedActualOut, actual)
		})
	}
}

func TestGetInputAndExpectedOutputByID(t *testing.T) {
	examples := []ProblemExample{
		{ID: 1, Input: "1", ExpectedOutput: "2"},
		{ID: 2, Input: "2", ExpectedOutput: "4"},
	}

	tests := []struct {
		name           string
		id             int
		expectedInput  string
		expectedOutput string
	}{
		{"IDExists", 1, "1", "2"},
		{"IDDoesNotExist", 3, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, output := getInputAndExpectedOutputByID(examples, tt.id)
			equals(t, tt.expectedInput, input)
			equals(t, tt.expectedOutput, output)
		})
	}
}

func TestGetOutputValue(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		expectedOutput string
	}{
		{"OutputExists", "Test 4: FAILED, Output: -1", "-1"},
		{"OutputDoesNotExist", "Test 4: FAILED, Output: ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equals(t, tt.expectedOutput, getOutputValue(tt.line))
		})
	}
}

func TestGetFailureError(t *testing.T) {
	tests := []struct {
		name          string
		line          string
		expectedError string
	}{
		{"ErrorExists", "./temp_code.go:4:5: invalid operation: mismatched types int and string", "invalid operation: mismatched types int and string"},
		{"ErrorDoesNotExist", "Test 4: FAILED, Output: -1", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equals(t, tt.expectedError, getFailureError(tt.line))
		})
	}
}
