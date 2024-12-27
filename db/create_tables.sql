-- Problems table: stores problem metadata
CREATE TABLE if NOT EXISTS problems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    short_description TEXT,
    long_description TEXT,
    problem_seed TEXT,
    examples TEXT,
    params TEXT,
    difficulty TEXT,
    attempts INTEGER DEFAULT 0,
    solves INTEGER DEFAULT 0
);

-- Problem examples table: stores inputs and expected outputs for validation
CREATE TABLE IF NOT EXISTS problem_examples (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    problem_id INTEGER,
    input TEXT NOT NULL,
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
