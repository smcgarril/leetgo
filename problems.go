package main

// Problem represents a LeetCode-style problem
type Problem struct {
	ID           string                                  `json:"id"`
	Title        string                                  `json:"title"`
	Description  string                                  `json:"description"`
	Example      string                                  `json:"example"`
	FunctionName string                                  `json:"function_name"`
	TestInput    string                                  `json:"test_input"`
	Expected     string                                  `json:"expected"`
	Solution     func(inputs ...interface{}) interface{} `json:"-"` // Ignore this field during JSON marshaling
}

// CodeReqeuest represents a user generated code snippet with test validation
type CodeRequest struct {
	Code      string `json:"code"`
	TestInput string `json:"test_input"`
	Problem   string `json:"problem"`
	Expected  string `json:"expected"`
}

// A list of sample problems
var Problems = []Problem{
	{
		ID:           "1",
		Title:        "Add 4",
		Description:  "Given one integer, add 4",
		Example:      "Input: 1, 2 | Output: 3",
		FunctionName: "sum",
		TestInput:    "1",
		Expected:     "5",
		Solution:     Sum, // Reference the solution function here
	},
	{
		ID:           "2",
		Title:        "Palidrome",
		Description:  "Given a string, return whether or not it is a palindrome.",
		Example:      "Input: bob | Output: true",
		FunctionName: "palindrome",
		TestInput:    "bob",
		Expected:     "true",
		Solution:     IsPalindrome, // Reference the solution function here
	},
	// You can add more problems here as needed
}

// // UnitTest function to test the solution of the problem
// func (p *Problem) UnitTestInt(t *testing.T, solution func(int int) int, expected int, input int) bool {
// 	result := p.Solution(input1, input2)
// 	if result != expected {
// 		t.Errorf("Problem %s: expected %d, got %d", p.Title, expected, result)
// 		return false
// 	}
// 	return true
// }

// Generic unit test wrapper
func UnitTest(problem Problem, output string) bool {
	return false
}
