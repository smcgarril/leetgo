-- Insert a new problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, problem_seed, difficulty) 
VALUES (
    1,
    'Palindrome',
    'Check if a string is a palindrome',
    'A phrase is a palindrome if, after converting all uppercase letters into lowercase letters and removing all non-alphanumeric characters, it reads the same forward and backward. Alphanumeric characters include letters and numbers.<br><br>Given a string <code>s</code>, return <code>true</code> if it is a palindrome, or <code>false</code> otherwise.',
    'func Palindrome(s string) bool {
    
}',
    'easy'
    );

-- Insert test cases for the "Palindrome" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, input_type, expected_output, output_type)
VALUES
(1, 1, '"radar"', '"string"', '"true"', '"bool"'),
(2, 1, '"hello"', '"string"', '"false"', '"bool"'),
(3, 1, '"A man a plan a canal Panama"', '"string"', '"true"', '"bool"');

-- Insert a new problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, problem_seed, difficulty) 
VALUES (
    2, 
    'Sum', 
    'Return the sum of two integers', 
    'Write a function that returns the sum of two integers.', 
    'func Sum(x, y int) int {
    
}', 
    'easy'
);

-- Insert test cases for the "Sum" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, input_type, expected_output, output_type)
VALUES
(4, 2, '[1, 2]', '"[]int"', '3', '"int"'),
(5, 2, '[-1, 2]', '"[]int"', '1', '"int"'),
(6, 2, '[0, 0]', '"[]int"', '0', '"int"');

-- Insert sample user
INSERT OR IGNORE INTO users (id, username, email, password)
VALUES (1, 'Test User', 'test@nowhere.com', '123456');
