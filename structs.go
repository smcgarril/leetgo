package main

// Problem represents a LeetCode-style problem
type Problem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	Difficulty       string `json:"difficulty"`
	Attempts         string `json:"attempts"`
	Solves           string `json:"solves"`
}

// ProblemExample represents a single input to a user function and the expected return
type ProblemExample struct {
	ID             string `json:"id"`
	PromblemID     string `json:"problem_id"`
	Input          string `json:"input"`
	InputType      string `json:"input_type"`
	ExpectedOutput string `json:"expected_output"`
	OutputType     string `json:"output_type"`
}

// CodeReqeuest represents a user generated code snippet with test validation
type CodeSubmission struct {
	Code      string `json:"code"`
	ProblemID string `json:"problem_id"`
	Problem   string `json:"problem"`
}

// CodeResponse respresents the results of a test execution
type CodeResponse struct {
	TestCount  int    `json:"testCount"`
	TestPassed int    `json:"testPassed"`
	Output     string `json:"output"`
}
