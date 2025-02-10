package api

import (
	"fmt"
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

// Mocks

// Tests

func TestFormatArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		keyOrder []string
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid input with all keys in order",
			input:    `{"name": "John", "age": 30, "active": true}`,
			keyOrder: []string{"name", "age", "active"},
			expected: `"John", 30, true`,
			wantErr:  false,
		},
		{
			name:     "Missing key in input",
			input:    `{"name": "John", "age": 30}`,
			keyOrder: []string{"name", "age", "active"},
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Unsupported value type",
			input:    `{"name": "John", "age": 30, "data": {"nested": "value"}}`,
			keyOrder: []string{"name", "age", "data"},
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty input",
			input:    `{}`,
			keyOrder: []string{"name"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatArgs(tt.input, tt.keyOrder)

			equals(t, tt.wantErr, (err != nil))
			equals(t, tt.expected, result)
		})
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "Float value",
			value:    42.5,
			expected: "42.5",
			wantErr:  false,
		},
		{
			name:     "String value",
			value:    "hello",
			expected: `"hello"`,
			wantErr:  false,
		},
		{
			name:     "Boolean value",
			value:    true,
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "Unsupported type",
			value:    map[string]interface{}{"key": "value"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatValue(tt.value)

			equals(t, tt.wantErr, (err != nil))
			equals(t, tt.expected, result)
		})
	}
}
