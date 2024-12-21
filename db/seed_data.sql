-- Insert a new problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, difficulty) 
VALUES (1, 'Palindrome', 'Check if a string is a palindrome', 'Write a function that checks if a string is the same forwards and backwards, ignoring spaces and case.', 'easy');

-- Insert test cases for the "Palindrome" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, input_type, expected_output, output_type)
VALUES
(1, 1, '"radar"', '"string"', '"true"', '"bool"'),
(2, 1, '"hello"', '"string"', '"false"', '"bool"'),
(3, 1, '"A man a plan a canal Panama"', '"string"', '"true"', '"bool"');

-- Insert a new problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, difficulty) 
VALUES (2, 'Sum', 'Return the sum of two integers', 'Write a function that returns the sum of two integers', 'easy');

-- Insert test cases for the "Sum" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, input_type, expected_output, output_type)
VALUES
(4, 2, '[1, 2]', '"[]int"', '3', '"int"'),
(5, 2, '[-1, 2]', '"[]int"', '1', '"int"'),
(6, 2, '[0, 0]', '"[]int"', '0', '"int"');

-- Insert sample user
INSERT OR IGNORE INTO users (id, username, email, password)
VALUES (1, 'Test User', 'test@nowhere.com', '123456');
