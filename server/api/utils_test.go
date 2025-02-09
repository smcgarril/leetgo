package api

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestGetWorkerURL(t *testing.T) {
	// Save the current environment variables to restore later
	originalWorkerHost := os.Getenv("WORKER_HOST")
	originalWorkerPort := os.Getenv("WORKER_PORT")
	originalWorkerPath := os.Getenv("WORKER_PATH")

	t.Run("NoEnvVars", func(t *testing.T) {
		os.Unsetenv("WORKER_HOST")
		os.Unsetenv("WORKER_PORT")
		os.Unsetenv("WORKER_PATH")

		expected := "http://localhost:8081/process-code"
		got := GetWorkerURL()

		equals(t, expected, got)
	})

	t.Run("OnlyWorkerHostSet", func(t *testing.T) {
		os.Setenv("WORKER_HOST", "http://example.com")
		os.Unsetenv("WORKER_PORT")
		os.Unsetenv("WORKER_PATH")

		expected := "http://example.com:8081/process-code"
		got := GetWorkerURL()

		equals(t, expected, got)
	})

	t.Run("OnlyWorkerPortSet", func(t *testing.T) {
		os.Unsetenv("WORKER_HOST")
		os.Setenv("WORKER_PORT", "9090")
		os.Unsetenv("WORKER_PATH")

		expected := "http://localhost:9090/process-code"
		got := GetWorkerURL()

		equals(t, expected, got)
	})

	t.Run("OnlyWorkerPathSet", func(t *testing.T) {
		os.Unsetenv("WORKER_HOST")
		os.Unsetenv("WORKER_PORT")
		os.Setenv("WORKER_PATH", "/custom-path")

		expected := "http://localhost:8081/custom-path"
		got := GetWorkerURL()

		equals(t, expected, got)
	})

	t.Run("AllEnvVarsSet", func(t *testing.T) {
		os.Setenv("WORKER_HOST", "http://example.com")
		os.Setenv("WORKER_PORT", "9090")
		os.Setenv("WORKER_PATH", "/custom-path")

		expected := "http://example.com:9090/custom-path"
		got := GetWorkerURL()

		equals(t, expected, got)
	})

	// Restore the original environment variables
	os.Setenv("WORKER_HOST", originalWorkerHost)
	os.Setenv("WORKER_PORT", originalWorkerPort)
	os.Setenv("WORKER_PATH", originalWorkerPath)
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

	t.Run("FailedTest", func(t *testing.T) {
		expectedInput := "1"
		expectedExpectedOutput := "2"
		expectedActualOutput := "-1"

		input, expectedOutput, actualOutput := BuildResponse(codeOutput, examples)

		equals(t, expectedInput, input)
		equals(t, expectedExpectedOutput, expectedOutput)
		equals(t, expectedActualOutput, actualOutput)
	})

	t.Run("NoFailedTest", func(t *testing.T) {
		codeOutput.Result = "SUCCESS"

		expectedInput := ""
		expectedExpectedOutput := ""
		expectedActualOutput := ""

		input, expectedOutput, actualOutput := BuildResponse(codeOutput, examples)

		equals(t, expectedInput, input)
		equals(t, expectedExpectedOutput, expectedOutput)
		equals(t, expectedActualOutput, actualOutput)
	})
}

func TestGetInputAndExpectedOutputByID(t *testing.T) {
	examples := []ProblemExample{
		{ID: 1, Input: "1", ExpectedOutput: "2"},
		{ID: 2, Input: "2", ExpectedOutput: "4"},
	}

	t.Run("IDExists", func(t *testing.T) {
		expectedInput := "1"
		expectedExpectedOutput := "2"

		input, expectedOutput := getInputAndExpectedOutputByID(examples, 1)

		equals(t, expectedInput, input)
		equals(t, expectedExpectedOutput, expectedOutput)
	})

	t.Run("IDDoesNotExist", func(t *testing.T) {
		expectedInput := ""
		expectedExpectedOutput := ""

		input, expectedOutput := getInputAndExpectedOutputByID(examples, 3)

		equals(t, expectedInput, input)
		equals(t, expectedExpectedOutput, expectedOutput)
	})
}

func TestGetOutputValue(t *testing.T) {
	t.Run("OutputExists", func(t *testing.T) {
		line := "Test 4: FAILED, Output: -1"
		expectedOutput := "-1"

		output := getOutputValue(line)

		equals(t, expectedOutput, output)
	})

	t.Run("OutputDoesNotExist", func(t *testing.T) {
		line := "Test 4: FAILED, Output: "
		expectedOutput := ""

		output := getOutputValue(line)

		equals(t, expectedOutput, output)
	})
}

func TestGetFailureError(t *testing.T) {
	t.Run("ErrorExists", func(t *testing.T) {
		line := "./temp_code.go:4:5: invalid operation: mismatched types int and string"
		expectedError := "invalid operation: mismatched types int and string"

		errorDetail := getFailureError(line)

		equals(t, expectedError, errorDetail)
	})

	t.Run("ErrorDoesNotExist", func(t *testing.T) {
		line := "Test 4: FAILED, Output: -1"
		expectedError := ""

		errorDetail := getFailureError(line)

		equals(t, expectedError, errorDetail)
	})
}
