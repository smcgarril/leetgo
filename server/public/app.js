const initialCode = 
`func main() {
    fmt.Println("Welcome to LeetGo!")
}`;

let editor;
let currentProblem = null;

// Seed editor with welcome message
function seedEditor(data) {
    editor.setValue(data);
}

// Function to fetch problems from the backend
async function fetchProblems() {
    try {
        const response = await fetch('/problems');
        if (!response.ok) {
            throw new Error("Failed to fetch problems");
        }

        const problems = await response.json();
        renderProblems(problems);
    } catch (error) {
        console.error('Error fetching problems:', error);
    }
}

// Function to render the fetched problems in the DOM as a dropdown
function renderProblems(problems) {
    const problemsDiv = document.getElementById('problems');
    problemsDiv.innerHTML = '';

    // Create the dropdown element
    const select = document.createElement('select');
    select.id = 'problem-select';
    select.onchange = () => displayProblemDescription(problems);

    // Add the default "Select a problem" option
    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.disabled = true;
    defaultOption.selected = true;
    defaultOption.innerText = 'Select a problem';
    select.appendChild(defaultOption);

    // Add an option for each problem
    problems.forEach(problem => {
        const option = document.createElement('option');
        option.value = problem.id;
        option.innerText = problem.name;
        select.appendChild(option);
    });

    problemsDiv.appendChild(select);
}

// Display problem description in left hand column
function displayProblemDescription(problems) {
    const select = document.getElementById('problem-select');
    const selectedProblemId = select.value;

    currentProblem = problems.find(p => p.id === selectedProblemId);

    const descriptionDiv = document.getElementById('problem-description');
    if (currentProblem) {
        descriptionDiv.innerHTML = `<h3>${currentProblem.name}</h3><p>${currentProblem.long_description}</p>`;
        renderExamples(currentProblem.examples);
        loadProblem();
        clearResults();
    } else {
        descriptionDiv.innerHTML = '';
    }
}

// Clear results
function clearResults() {
    const resultsElement = document.getElementById('results');
    const failureDetailsElement = document.getElementById('failure-details');

    if (resultsElement) {
        resultsElement.style.display = 'none';
        document.getElementById('result').innerText = '';
        document.getElementById('testPassed').innerText = '';
        document.getElementById('testCount').innerText = '';
    }

    if (failureDetailsElement) {
        failureDetailsElement.style.display = 'none';
        document.getElementById('failure-input').innerText = '';
        document.getElementById('failure-expected').innerText = '';
        document.getElementById('failure-actual').innerText = '';
    }
}

// Load problem examples
function renderExamples(examples) {
    console.log("examples is: ", examples);
    try {
        examples = JSON.parse(examples);
    } catch (error) {
        console.error("Failed to parse examples:", error);
    }

    const examplesDiv = document.getElementById('examples');
    examplesDiv.innerHTML = ''; // Clear previous examples

    // Check if examples is an array
    if (Array.isArray(examples)) {
        examples.forEach((example, index) => {
            const exampleContainer = document.createElement('div');
            exampleContainer.style.marginBottom = '1em';

            exampleContainer.innerHTML = `
                <p><strong>Example ${index + 1}:</strong></p>
                <p style="margin-left: 1em;">Input: ${example.input}</p>
                <p style="margin-left: 1em;">Output: ${example.output}</p>
                <p style="margin-left: 1em;">Explanation: ${example.explanation}</p>
            `;

            examplesDiv.appendChild(exampleContainer);
        });
    } else {
        console.error("Examples is not an array:", examples);
    }
}

// Function to load selected problem into the editor
function loadProblem() {
    if (currentProblem) {
        editor.setValue(currentProblem.problem_seed);
    }
}

// Dark Mode Toggle
const darkModeToggle = document.getElementById('darkModeToggle');
const darkModeIcon = document.getElementById('darkModeIcon');

darkModeToggle.addEventListener('click', () => {
    // document.body.classList.toggle('dark-mode');
    const darkModeEnabled = document.body.classList.toggle("dark-mode");
    localStorage.setItem("darkMode", darkModeEnabled);

    // Toggle between moon and sun icons
    if (document.body.classList.contains('dark-mode')) {
        darkModeIcon.classList.replace('fa-moon', 'fa-sun');
    } else {
        darkModeIcon.classList.replace('fa-sun', 'fa-moon');
    }
});

// Send code input to back end
async function runCode() {
    if (!currentProblem) {
        alert("Please select a problem first.");
        return;
    }

    const code = editor.getValue();

    const problem_id = currentProblem.id;
    const problem = currentProblem.name;

    try {
        const payload = {
            code: code,    
            problem_id: problem_id,
            problem: problem
        };

        console.log(payload);
        console.log(JSON.stringify({ payload }));

        const response = await fetch('/execute', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            console.log(response);
            throw new Error('Failed to execute code');
        }

        const rawText = await response.text();
        console.log(rawText)
        const data = JSON.parse(rawText);
        console.log(data)

        if (data.result) {
            // Show the results
            const resultsElement = document.getElementById('results');
            if (resultsElement) {
                resultsElement.style.display = 'block'; // Show failure details
            }

            // Update result element
            const resultElement = document.getElementById('result');
            resultElement.innerText = data.result.trim();
            
            // Add success or failure class
            if (data.result === "PASSED") {
                resultElement.classList.remove("failure");
                resultElement.classList.add("success");
                
                // Hide the failure details if result is PASSED
                const failureDetailsElement = document.getElementById('failure-details');
                if (failureDetailsElement) {
                    failureDetailsElement.style.display = 'none'; // Hide failure details
                }
            } else if (data.result === "FAILED") {
                resultElement.classList.remove("success");
                resultElement.classList.add("failure");
                
                // Show the failure details if result is FAILED
                const failureDetailsElement = document.getElementById('failure-details');
                if (failureDetailsElement) {
                    failureDetailsElement.style.display = 'block'; // Show failure details
                }
                
                // Parse and format the input if it exists, otherwise render an empty string
                let formattedInput = "";
                if (data.input || data.input === 0) {
                    try {
                        const parsedInput = JSON.parse(data.input);
                        formattedInput = Object.entries(parsedInput)
                            .map(([key, value]) => `${key} = ${value}`)
                            .join(', ');
                    } catch (e) {
                        formattedInput = "";
                    }
                }

                // Parse and extract the expected output value if it exists, otherwise render an empty string
                let formattedExpectedOutput = "";
                if (data.expected) {
                    try {
                        const parsedExpectedOutput = JSON.parse(data.expected);
                        formattedExpectedOutput = Object.values(parsedExpectedOutput).join(', ');
                    } catch (e) {
                        formattedExpectedOutput = "";
                    }
                }

                // Update failure details
                document.getElementById('failure-input').innerText = formattedInput;
                document.getElementById('failure-expected').innerText = formattedExpectedOutput;
                document.getElementById('failure-actual').innerText = data.output;
            }
        } else {
            document.getElementById('result').innerText = 'No result received';
        }
        
        if (data.testPassed || data.testPassed === 0) {
            document.getElementById('testPassed').innerText = data.testPassed;
        } else {
            document.getElementById('testPassed').innerText = 'No test passed information received';
        }
        
        if (data.testCount) {
            document.getElementById('testCount').innerText = data.testCount;
        } else {
            document.getElementById('testCount').innerText = 'No test count received';
        }

    } catch (error) {
        console.error('Error running code:', error);
        document.getElementById('result').innerText = 'Error running code: ' + error.message;
    }
}


// Call fetchProblems when the page loads
fetchProblems();

document.addEventListener("DOMContentLoaded", function () {
    fetchProblems();

    const body = document.body;
    const isDarkMode = localStorage.getItem("darkMode") === "true";
    if (isDarkMode) {
        body.classList.add("dark-mode");
        darkModeIcon.classList.replace('fa-sun', 'fa-moon');
    }

    // Initialize CodeMirror editor for the user to input code
    editor = CodeMirror(document.getElementById('editor'), {
        mode: 'text/x-go',
        lineNumbers: true,
        theme: 'default',
        indentUnit: 4,
        tabSize: 4,
    });

    seedEditor(initialCode);
});