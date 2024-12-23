package api

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/iancoleman/orderedmap"
)

// Unmarshal JSON and format for unit test execution
func FormatTestJSON(inputJSON string) string {
	o := orderedmap.New()
	err := json.Unmarshal([]byte(inputJSON), &o)
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}

	keys := o.Keys()
	outputString, sep := "", ""
	for _, k := range keys {
		v, _ := o.Get(k)
		s, ok := v.(string)
		if ok {
			outputString += sep + s
			sep = ", "
		} else {
			log.Fatalf("JSON value is not a string.")
		}
	}

	outputString = strings.TrimSuffix(outputString, ", ")
	outputString = NormalizeString(outputString)

	return outputString
}

// Regex to remove leading or trailing backslashes
func RemoveOuterBackslashes(s string) string {
	re := regexp.MustCompile(`^\\|\\$`)
	return re.ReplaceAllString(s, "")
}

// Replace escaped quotes (\") with actual quotes (")
func NormalizeString(s string) string {
	s = strings.ReplaceAll(s, `\"`, `"`)
	return s
}
