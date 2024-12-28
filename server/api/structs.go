package api

// Problem represents a LeetCode-style problem
type Problem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	Difficulty       string `json:"difficulty"`
	ProblemSeed      string `json:"problem_seed"`
	Examples         string `json:"examples"`
	Attempts         string `json:"attempts"`
	Solves           string `json:"solves"`
}

// ProblemExample represents a single input to a user function and the expected return
type ProblemExample struct {
	ID             string `json:"id"`
	PromblemID     string `json:"problem_id"`
	Input          string `json:"input"`
	InputOrder     string `json:"input_order"`
	ExpectedOutput string `json:"expected_output"`
}

// CodeReqeuest represents a user generated code snippet with test validation
type CodeSubmission struct {
	Code            string           `json:"code"`
	ProblemID       string           `json:"problem_id"`
	Problem         string           `json:"problem"`
	ProblemExamples []ProblemExample `json:"problem_examples"`
}

// CodeResponse respresents the results of a test execution
type CodeResponse struct {
	TestCount  int    `json:"testCount"`
	TestPassed int    `json:"testPassed"`
	Output     string `json:"output"`
	Result     string `json:"result"`
}
