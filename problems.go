package main

// Problem represents a LeetCode-style problem
type Problem struct {
	ID          string                                  `json:"id"`
	Title       string                                  `json:"title"`
	Description string                                  `json:"description"`
	Example     string                                  `json:"example"`
	Solution    func(inputs ...interface{}) interface{} `json:"-"` // Ignore this field during JSON marshaling
}

// A list of sample problems
var Problems = []Problem{
	{
		ID:          "1",
		Title:       "Two Sum",
		Description: "Given two integers, return their sum.",
		Example:     "Input: 1, 2 | Output: 3",
		Solution:    Sum, // Reference the solution function here
	},
	{
		ID:          "2",
		Title:       "Palidrome",
		Description: "Given a string, return whether or not it is a palindrome.",
		Example:     "Input: bob | Output: true",
		Solution:    IsPalindrome, // Reference the solution function here
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
