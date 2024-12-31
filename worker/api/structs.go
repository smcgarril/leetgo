package api

type CodeSubmission struct {
	Code            string           `json:"code"`
	Problem         string           `json:"problem"`
	ProblemExamples []ProblemExample `json:"problem_examples"`
}

type ProblemExample struct {
	ID             int    `json:"id"`
	Input          string `json:"input"`
	InputOrder     string `json:"input_order"`
	ExpectedOutput string `json:"expected_output"`
}

// CodeOutput respresents the results of a test execution
type CodeOutput struct {
	TestCount  int    `json:"testCount"`
	TestPassed int    `json:"testPassed"`
	Output     string `json:"output"`
	Input      string `json:"input"`
	Expected   string `json:"expected"`
	Result     string `json:"result"`
}
