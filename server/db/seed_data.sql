-- Insert "Palindrome" problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, problem_seed, examples, difficulty) 
VALUES (
    1,
    'Palindrome',
    'Check if a string is a palindrome',
    'A phrase is a palindrome if, after converting all uppercase letters into lowercase letters and removing all non-alphanumeric characters, it reads the same forward and backward. Alphanumeric characters include letters and numbers.<br><br>Given a string <code>s</code>, return <code>true</code> if it is a palindrome, or <code>false</code> otherwise.',
    'func Palindrome(s string) bool {
    
}',
    '[
    {
        "input": "s = \\"radar\\"",
        "output": "true",
        "explanation": "\\"radar\\" is a palindrome."
    },
    {
        "input": "s = \\"hello\\"",
        "output": "false",
        "explanation": "\\"hello\\" is not a palindrome."
    },
    {
        "input": "s = \\"A man, a plan, a canal: Panama\\"",
        "output": "true",
        "explanation": "\\"amanaplanacanalpanama\\" is a palindrome."
    }
]',
    'easy'
    );

-- Insert test cases for the "Palindrome" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, input_order, expected_output)
VALUES
(1, 1, '{"s": "radar"}', '["s"]', '{"result": true}'),
(2, 1, '{"s": "hello"}', '["s"]', '{"result": false}'),
(3, 1, '{"s": "A man a plan a canal Panama"}', '["s"]', '{"result": true}');

-- Insert "Sum" problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, problem_seed, examples, difficulty) 
VALUES (
    2, 
    'Sum', 
    'Return the sum of two integers', 
    'Write a function that returns the sum of two integers.', 
    'func Sum(x, y int) int {
    
}', 
    '[
    {
        "input": "[1, 2]",
        "output": "3",
        "explanation": "1 plus 2 equals 3."
    },
    {
        "input": "[-1, 2]",
        "output": "1",
        "explanation": "-1 plus 2 equals 1."
    },
    {
        "input": "[0, 0]",
        "output": "0",
        "explanation": "0 plus 0 equals 0."
    }
]',
    'easy'
);

-- Insert test cases for the "Sum" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, input_order, expected_output)
VALUES
(4, 2, '{"x": 1, "y": 2}', '["x", "y"]', '{"result": 3}'),
(5, 2, '{"x": -1, "y": 2}', '["x", "y"]', '{"result": 1}'),
(6, 2, '{"x": 0, "y": 0}', '["x", "y"]', '{"result": 0}');

-- Insert "Two Sum" problem
INSERT OR IGNORE INTO problems (id, name, short_description, long_description, problem_seed, examples, difficulty) 
VALUES (
    3, 
    'TwoSum', 
    'Return the indexes of the two numbers that sum to the target', 
    'Given an array of integers <code>nums</code> and an integer <code>target</code>, return <i>indices of the two numbers such that they add up to <code>target</code></i>.<br></br>You may assume that each input would have <strong><i>exactly</i> one solution</strong>, and you may not use the <i>same</i> element twice.<br></br>You can return the answer in any order.', 
    'func TwoSum(nums []int, target int) []int {
    
}', 
    '[
    {
        "input": "nums = [2,7,11,15], target = 9",
        "output": "[0,1]",
        "explanation": "Because nums[0] + nums[1] == 9, we return [0, 1]."
    },
    {
        "input": "nums = [3,2,4], target = 6",
        "output": "[1,2]",
        "explanation": "Because nums[1] + nums[2] == 6, we return [1, 2]."
    },
    {
        "input": "nums = [3,3], target = 6",
        "output": "[0,1]",
        "explanation": "Because nums[0] + nums[1] == 6, we return [0, 1]."
    }
]',
    'easy'
);

-- Insert test cases for the "Two Sum" problem
INSERT OR IGNORE INTO problem_examples (id, problem_id, input, input_order, expected_output)
VALUES
(7, 3, '{"nums": [2, 7, 11, 15], "target": 9}', '["nums", "target"]', '{"indices": [0, 1]}'),
(8, 3, '{"nums": [3, 2, 4], "target": 6}', '["nums", "target"]', '{"indices": [1, 2]}'),
(9, 3, '{"nums": [3, 3], "target": 6}', '["nums", "target"]', '{"indices": [0, 1]}');

-- Insert sample user
INSERT OR IGNORE INTO users (id, username, email, password)
VALUES (1, 'Test User', 'test@nowhere.com', '123456');
