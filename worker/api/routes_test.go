package api

import (
	"strings"
	"testing"
)

func TestGenerateTestHarness(t *testing.T) {
	tests := []struct {
		name      string
		userCode  string
		testCalls string
		expected  []string
	}{
		{
			name: "Simple test harness",
			userCode: `
			func add(a, b int) int {
				return a + b
			}`,
			testCalls: `
			results = append(results, Result{
				Test:   1,
				Result: "PASS",
				Output: fmt.Sprint(add(2, 3)),
			})
			results = append(results, Result{
				Test:   2,
				Result: "FAIL",
				Output: fmt.Sprint(add(1, -1)),
			})`,
			expected: []string{
				"package main",
				"type Result struct {",
				"func add(a, b int) int {",
				"results = append(results, Result{",
				"fmt.Sprint(add(2, 3))",
				"Output: fmt.Sprint(add(1, -1))",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateTestHarness(tt.userCode, tt.testCalls)

			// Ensure the generated code contains the expected strings
			for _, expectedFragment := range tt.expected {
				if !strings.Contains(result, expectedFragment) {
					t.Errorf("Expected generated code to contain: %q, but it was missing", expectedFragment)
				}
			}

			// Ensure the result is a valid Go program structure (basic syntax check)
			if !strings.HasPrefix(result, "\n\t\tpackage main\n\t\timport") {
				t.Errorf("Generated code does not have a valid package declaration or import statement")
			}
		})
	}
}
