package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Helper function to transform JSON input to formatted arguments
func FormatArgs(input string, keyOrder []string) (string, error) {
	var args map[string]interface{}
	err := json.Unmarshal([]byte(input), &args)
	if err != nil {
		return "", err
	}

	var formattedArgs []string
	for _, key := range keyOrder {
		value, exists := args[key]
		if !exists {
			return "", fmt.Errorf("missing key: %s", key)
		}
		switch v := value.(type) {
		case float64:
			// Check if the number is an integer
			if v == float64(int(v)) {
				formattedArgs = append(formattedArgs, fmt.Sprintf("%d", int(v)))
			} else {
				formattedArgs = append(formattedArgs, fmt.Sprintf("%g", v))
			}
		case string:
			formattedArgs = append(formattedArgs, fmt.Sprintf("%q", v))
		case bool:
			formattedArgs = append(formattedArgs, fmt.Sprintf("%t", v))
		case []interface{}:
			formattedArray, err := FormatArray(v)
			if err != nil {
				return "", fmt.Errorf("error formatting array for key %s: %v", key, err)
			}
			formattedArgs = append(formattedArgs, formattedArray)
		default:
			return "", fmt.Errorf("unsupported type for key %s: %T", key, v)
		}
	}
	return strings.Join(formattedArgs, ", "), nil
}

// Helper function to format JSON arrays
func FormatArray(array []interface{}) (string, error) {
	var elements []string
	isInt := true

	for _, elem := range array {
		switch v := elem.(type) {
		case float64:
			// Check if the number is an integer
			if v == float64(int(v)) {
				elements = append(elements, fmt.Sprintf("%d", int(v)))
			} else {
				isInt = false
				elements = append(elements, fmt.Sprintf("%g", v))
			}
		case string:
			isInt = false
			elements = append(elements, fmt.Sprintf("%q", v))
		case bool:
			isInt = false
			elements = append(elements, fmt.Sprintf("%t", v))
		default:
			return "", fmt.Errorf("unsupported array element type: %T", v)
		}
	}

	arrayType := "int"
	if !isInt {
		arrayType = "float64"
	}

	return fmt.Sprintf("[]%s{%s}", arrayType, strings.Join(elements, ", ")), nil
}

// Helper function to transform JSON expected output
func FormatExpectedOutput(output string) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(output), &result)
	if err != nil {
		return "", err
	}

	for _, value := range result {
		switch v := value.(type) {
		case float64:
			if v == float64(int(v)) {
				return fmt.Sprintf("%d", int(v)), nil
			}
			return fmt.Sprintf("%g", v), nil
		case string:
			return fmt.Sprintf("%q", v), nil
		case bool:
			return fmt.Sprintf("%t", v), nil
		case []interface{}:
			return FormatArray(v)
		default:
			return "", fmt.Errorf("unsupported type: %T", v)
		}
	}
	return "", nil
}

// Helper function to count number of passing tests
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
