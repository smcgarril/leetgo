package main

// Sum takes two ints and returns their sum
func Sum(inputs ...interface{}) interface{} {
	a := inputs[0].(int)
	b := inputs[1].(int)
	return a + b
}

// IsPalindrome checks if a given string is a palindrome
func IsPalindrome(inputs ...interface{}) interface{} {
	str := inputs[0].(string)
	length := len(str)
	for i := 0; i < length/2; i++ {
		if str[i] != str[length-i-1] {
			return false
		}
	}
	return true
}
