-- Insert a new problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, difficulty) 
VALUES (1, 'Palindrome', 'Check if a string is a palindrome', 'Write a function that checks if a string is the same forwards and backwards, ignoring spaces and case.', 'easy');

-- Insert test cases for the "Palindrome" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, expected_output)
VALUES
(1, 1, '"radar"', '"true"'),
(2, 1, '"hello"', '"false"'),
(3, 1, '"A man a plan a canal Panama"', '"true"');

-- Insert a new problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, difficulty) 
VALUES (2, 'Sum', 'Return the sum of two integers', 'Write a function that returns the sum of two integers', 'easy');

-- Insert test cases for the "Sum" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, expected_output)
VALUES
(4, 2, '[1, 2]', '3'),
(5, 2, '[-1, 2]', '1'),
(6, 2, '[0, 0]', '0');

-- Insert sample user
INSERT OR IGNORE INTO users (id, username, email, password)
VALUES (1, 'Test User', 'test@nowhere.com', '123456');
