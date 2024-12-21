package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// convertToType dynamically converts a string to its Go type
func ConvertToType(input string) (interface{}, error) {
	input = strings.TrimSpace(input)

	// Try unmarshaling as JSON (handles arrays, maps, etc.)
	var result interface{}
	if err := json.Unmarshal([]byte(input), &result); err == nil {
		return result, nil
	}

	// If JSON unmarshaling fails, try parsing primitive types
	if i, err := strconv.Atoi(input); err == nil {
		return i, nil // Integer
	}

	if f, err := strconv.ParseFloat(input, 64); err == nil {
		return f, nil // Float
	}

	// Assume it's a string if no other parsing succeeded
	return input, nil
}

func FormatFunctionArguments(input interface{}) (string, error) {
	// Handle the case where the input is a string
	if str, ok := input.(string); ok {
		args := []string{}
		for _, char := range str {
			args = append(args, fmt.Sprintf("%#v", string(char)))
		}
		return strings.Join(args, ", "), nil
	}

	// Handle the case where the input is a slice of interface{}
	slice, ok := input.([]interface{})
	if !ok {
		return "", fmt.Errorf("input is neither a string nor a slice: %v", input)
	}

	args := []string{}
	for _, val := range slice {
		args = append(args, fmt.Sprintf("%#v", val))
	}
	return strings.Join(args, ", "), nil
}
