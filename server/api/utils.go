package api

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Retrieve worker service URL:PORT from env variables
func GetWorkerURL() string {
	workerHost := os.Getenv("WORKER_HOST")
	if workerHost == "" {
		workerHost = "http://localhost"
	}

	workerPort := os.Getenv("WORKER_PORT")
	if workerPort == "" {
		workerPort = "8081"
	}

	workerPath := os.Getenv("WORKER_PATH")
	if workerPath == "" {
		workerPath = "/process-code"
	}

	return fmt.Sprintf("%s:%s%s", workerHost, workerPort, workerPath)
}

// Return input, expectedOutput, and actualOutput from CodeOutput
func BuildResponse(codeOutput *CodeOutput, examples []ProblemExample) (input, expectedOutput, actualOutput string) {
	if codeOutput.Result == "FAILED" {
		lines := strings.Split(codeOutput.Output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "./temp_code.go") {
				actualOutput = getFailureError(line)
				break
			}
			if strings.Contains(line, "FAILED") {
				re := regexp.MustCompile(`\d+`)
				id := re.FindString(line)
				testID, err := strconv.Atoi(id)
				if err != nil {
					fmt.Println("Error converting number: ", err)
					return
				}
				actualOutput = getOutputValue(line)
				input, expectedOutput = getInputAndExpectedOutputByID(examples, testID)
				break
			}
		}
	}
	return input, expectedOutput, actualOutput
}

// Return Input and ExpectedOutput by ID
func getInputAndExpectedOutputByID(examples []ProblemExample, id int) (string, string) {
	for _, example := range examples {
		if example.ID == id {
			return example.Input, example.ExpectedOutput
		}
	}
	return "", ""
}

// Return value provided after Output
func getOutputValue(line string) string {
	prefix := "Output: "
	startIndex := strings.Index(line, prefix)
	if startIndex == -1 {
		return ""
	}

	output := line[startIndex+len(prefix):]

	return strings.TrimSpace(output)
}

// Return error after code file & line numbers
func getFailureError(line string) string {
	re := regexp.MustCompile(`.*?:\d+:\d+: (.+)`)
	match := re.FindStringSubmatch(line)
	if len(match) > 1 {
		errorDetail := match[1]
		return errorDetail
	} else {
		return ""
	}
}
