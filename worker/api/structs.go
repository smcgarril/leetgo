package api

type CodeSubmission struct {
	Code            string           `json:"code"`
	Problem         string           `json:"problem"`
	ProblemExamples []ProblemExample `json:"problem_examples"`
}

type ProblemExample struct {
	Input          string `json:"input"`
	InputOrder     string `json:"input_order"`
	ExpectedOutput string `json:"expected_output"`
}

type CodeResponse struct {
	TestCount  int    `json:"testCount"`
	TestPassed int    `json:"testPassed"`
	Output     string `json:"output"`
	Result     string `json:"result"` // "PASSED" if all tests pass, otherwise "FAILED"
}
