const initialCode = `func main() {
    fmt.Println("Welcome to LeetGo!")
}`;

let editor;
let currentProblem = null;

document.addEventListener("DOMContentLoaded", init);

// Initialize the editor and fetch the problem list
function init() {
    initializeEditor();
    setupDarkMode();
    fetchProblemList();
}

// Initialize CodeMirror editor
function initializeEditor() {
    editor = CodeMirror(document.getElementById('editor'), {
        mode: 'text/x-go',
        lineNumbers: true,
        theme: 'default',
        indentUnit: 4,
        tabSize: 4,
    });
    seedEditor(initialCode);
}

// Seed editor with initial or problem seed code
function seedEditor(data) {
    editor.setValue(data);
}

// Fetch problem names and IDs from the backend
async function fetchProblemList() {
    try {
        const response = await fetch('/problems/names');
        if (!response.ok) throw new Error("Failed to fetch problem list");
        
        const problemList = await response.json();
        renderProblemsDropdown(problemList);
    } catch (error) {
        logError('Error fetching problem list:', error);
    }
}

// Render the dropdown with problem names
function renderProblemsDropdown(problems) {
    const problemsDiv = document.getElementById('problems');
    problemsDiv.innerHTML = '';

    const select = createDropdown(problems);
    problemsDiv.appendChild(select);
}

// Create a dropdown menu for problem selection
function createDropdown(problems) {
    const select = document.createElement('select');
    select.id = 'problem-select';
    select.onchange = () => fetchProblemDetails(select.value);

    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.disabled = true;
    defaultOption.selected = true;
    defaultOption.innerText = 'Select a problem';
    select.appendChild(defaultOption);

    problems.forEach(problem => {
        const option = document.createElement('option');
        option.value = problem.id;
        option.innerText = problem.name;
        select.appendChild(option);
    });

    return select;
}

// Fetch full problem details when a user selects a problem
async function fetchProblemDetails(problemId) {
    if (!problemId) return;

    try {
        const response = await fetch(`/problems/${problemId}`);
        if (!response.ok) throw new Error("Failed to fetch problem details");
        
        const problem = await response.json();
        currentProblem = problem[0];

        displayProblemDetails(currentProblem);   
        loadProblemSeed(currentProblem.problem_seed);
        clearResults();
    } catch (error) {
        logError('Error fetching problem details:', error);
    }
}

// Display the selected problem's details and examples
function displayProblemDetails(problem) {
    const descriptionDiv = document.getElementById('problem-description');
    if (problem) {
        descriptionDiv.innerHTML = `<h3>${problem.name}</h3><p>${problem.long_description}</p>`;
        renderExamples(problem.examples);
    } else {
        descriptionDiv.innerHTML = '';
    }
}

// Load problem seed code into the editor
function loadProblemSeed(seedCode) {
    seedEditor(seedCode);
}

// Render examples for the selected problem
function renderExamples(examples) {
    const examplesDiv = document.getElementById('examples');
    examplesDiv.innerHTML = '';

    try {
        const parsedExamples = JSON.parse(examples);
        if (Array.isArray(parsedExamples)) {
            parsedExamples.forEach((example, index) => {
                const exampleHtml = `
                    <div style="margin-bottom: 1em;">
                        <p><strong>Example ${index + 1}:</strong></p>
                        <p style="margin-left: 1em;">Input: ${example.input}</p>
                        <p style="margin-left: 1em;">Output: ${example.output}</p>
                        <p style="margin-left: 1em;">Explanation: ${example.explanation}</p>
                    </div>`;
                examplesDiv.insertAdjacentHTML('beforeend', exampleHtml);
            });
        } else {
            throw new Error("Examples data is not an array");
        }
    } catch (error) {
        logError("Failed to parse examples:", error);
    }
}

// Clear previous results from the UI
function clearResults() {
    ['results', 'failure-details'].forEach(id => {
        const element = document.getElementById(id);
        if (element) element.style.display = 'none';
    });

    ['result', 'testPassed', 'testCount', 'failure-input', 'failure-expected', 'failure-actual'].forEach(id => {
        const element = document.getElementById(id);
        if (element) element.innerText = '';
    });
}

// Set up dark mode toggle
function setupDarkMode() {
    const darkModeToggle = document.getElementById('darkModeToggle');
    const darkModeIcon = document.getElementById('darkModeIcon');

    const isDarkMode = localStorage.getItem("darkMode") === "true";
    if (isDarkMode) document.body.classList.add("dark-mode");

    darkModeToggle.addEventListener('click', () => toggleDarkMode(darkModeIcon));
}

// Toggle dark mode and update localStorage
function toggleDarkMode(icon) {
    const darkModeEnabled = document.body.classList.toggle("dark-mode");
    localStorage.setItem("darkMode", darkModeEnabled);

    icon.classList.replace(
        darkModeEnabled ? 'fa-moon' : 'fa-sun',
        darkModeEnabled ? 'fa-sun' : 'fa-moon'
    );
}

// Run the user-submitted code
async function runCode() {
    if (!currentProblem) return alert("Please select a problem first.");

    const code = editor.getValue();
    const payload = {
        code,
        problem_id: currentProblem.id,
        problem: currentProblem.name,
    };

    try {
        const response = await fetch('/execute', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });

        if (!response.ok) throw new Error('Failed to execute code');
        
        const data = await response.json();
        displayResults(data);
    } catch (error) {
        logError('Error running code:', error);
        document.getElementById('result').innerText = `Error running code: ${error.message}`;
    }
}

// Display code execution results
function displayResults(data) {
    const resultsElement = document.getElementById('results');
    const resultElement = document.getElementById('result');
    const failureDetailsElement = document.getElementById('failure-details');

    resultsElement.style.display = 'block';

    resultElement.innerText = data.result.trim();

    resultElement.classList.remove('success', 'failure');

    if (data.result === "PASSED") {
        resultElement.classList.add('success');
        failureDetailsElement.style.display = 'none';
    } else if (data.result === "FAILED") {
        resultElement.classList.add('failure');
        failureDetailsElement.style.display = 'block'; 
        displayFailureDetails(data);
    }

    document.getElementById('testPassed').innerText = data.testPassed ?? 'N/A';
    document.getElementById('testCount').innerText = data.testCount ?? 'N/A';
}

// Display failure details in a separate block
function displayFailureDetails(data) {
    const failureInputElement = document.getElementById('failure-input');
    const failureExpectedElement = document.getElementById('failure-expected');
    const failureActualElement = document.getElementById('failure-actual');

    failureInputElement.innerText = formatDataForDisplay(data.input);
    failureExpectedElement.innerText = formatDataForDisplay(data.expected);
    failureActualElement.innerText = data.output ?? '';
}

// Format data for display
function formatDataForDisplay(data) {
    if (!data && data !== 0) return '';
    try {
        const parsedData = JSON.parse(data);
        return Object.entries(parsedData)
            .map(([key, value]) => `${key} = ${value}`)
            .join(', ');
    } catch {
        return String(data);
    }
}

// Parse JSON to string for display
function parseJsonToString(jsonData) {
    try {
        const parsedData = JSON.parse(jsonData);
        return Object.entries(parsedData)
            .map(([key, value]) => `${key} = ${value}`)
            .join(', ');
    } catch {
        return jsonData || '';
    }
}

function logError(message, error) {
    console.error(`${message} ${error.message}`);
}