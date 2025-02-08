package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Transform the input JSON into a formatted string based on the given key order.
func FormatArgs(input string, keyOrder []string) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %w", err)
	}

	var formattedArgs []string
	for _, key := range keyOrder {
		value, exists := args[key]
		if !exists {
			return "", fmt.Errorf("missing key: %s", key)
		}

		formattedValue, err := formatValue(value)
		if err != nil {
			return "", fmt.Errorf("error formatting key %s: %v", key, err)
		}
		formattedArgs = append(formattedArgs, formattedValue)
	}

	return strings.Join(formattedArgs, ", "), nil
}

// Format a single value based on its type.
func formatValue(value interface{}) (string, error) {
	switch v := value.(type) {
	case float64:
		// Return as integer if it has no fractional part.
		if v == float64(int(v)) {
			return strconv.Itoa(int(v)), nil
		}
		return fmt.Sprintf("%g", v), nil
	case string:
		return fmt.Sprintf("%q", v), nil
	case bool:
		return fmt.Sprintf("%t", v), nil
	case []interface{}:
		return formatArray(v)
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}

// Format an array into a string representation.
func formatArray(array []interface{}) (string, error) {
	var elements []string
	isIntArray := true

	for _, elem := range array {
		formattedElem, err := formatValue(elem)
		if err != nil {
			return "", fmt.Errorf("unsupported array element type: %v", err)
		}
		elements = append(elements, formattedElem)

		// Check if all elements are integers
		if _, ok := elem.(float64); !ok || elem != float64(int(elem.(float64))) {
			isIntArray = false
		}
	}

	arrayType := "int"
	if !isIntArray {
		arrayType = "float64"
	}

	return fmt.Sprintf("[]%s{%s}", arrayType, strings.Join(elements, ", ")), nil
}

// Transform the expected output JSON into a formatted string.
func FormatExpectedOutput(output string) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return "", fmt.Errorf("failed to parse expected output JSON: %w", err)
	}

	for _, value := range result {
		return formatValue(value) // Return the formatted first value.
	}
	return "", fmt.Errorf("no value found in expected output")
}

// Count the number of passing tests in the given output string.
func CountPassingTests(output string) int {
	testPassed := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "PASSED") {
			testPassed++
		}
	}
	return testPassed
}
