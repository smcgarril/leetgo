-- Problems table: stores problem metadata
CREATE TABLE if NOT EXISTS problems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    short_description TEXT,
    long_description TEXT,
    problem_seed TEXT,
    examples TEXT,
    difficulty TEXT,
    attempts INTEGER DEFAULT 0,
    solves INTEGER DEFAULT 0
);

-- Problem examples table: stores inputs and expected outputs for validation
CREATE TABLE IF NOT EXISTS problem_examples (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    problem_id INTEGER,
    input TEXT NOT NULL,
    input_order TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

-- Problem images table (optional): stores images related to the problem
CREATE TABLE IF NOT EXISTS problem_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    problem_id INTEGER,
    image_url TEXT,
    description TEXT,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

-- Users table: stores users who are solving problems
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    password TEXT NOT NULL
);

-- User solutions table: stores the solutions submitted by users
CREATE TABLE IF NOT EXISTS user_solutions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    problem_id INTEGER,
    user_id INTEGER,
    solution_code TEXT NOT NULL,
    status TEXT,
    date_submitted DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


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
