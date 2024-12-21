let editor;
let currentProblem = null;

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

function displayProblemDescription(problems) {
    const select = document.getElementById('problem-select');
    const selectedProblemId = select.value;

    currentProblem = problems.find(p => p.id === selectedProblemId);

    const descriptionDiv = document.getElementById('problem-description');
    if (currentProblem) {
        descriptionDiv.innerHTML = `<h3>${currentProblem.name}</h3><p>${currentProblem.short_description}</p>`;
    } else {
        descriptionDiv.innerHTML = '';
    }
}

// Dark Mode Toggle
const darkModeToggle = document.getElementById('darkModeToggle');
darkModeToggle.addEventListener('change', () => {
    document.body.classList.toggle('dark-mode', darkModeToggle.checked);
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
        const data = JSON.parse(rawText);

        if (data.output) {
            // Update output element
            const outputElement = document.getElementById('output');
            outputElement.innerText = data.output.trim();
            // Add success or failure class
            if (data.output === "PASSED") {
                outputElement.classList.remove("failure");
                outputElement.classList.add("success");
            } else {
                outputElement.classList.remove("success");
                outputElement.classList.add("failure");
            }
        } else {
            document.getElementById('output').innerText = 'No output received';
        }
        if (data.testPassed) {
            document.getElementById('testPassed').innerText = data.testPassed;
        } else {
            document.getElementById('output').innerText = 'No output received';
        }
        if (data.testCount) {
            document.getElementById('testCount').innerText = data.testCount;
        } else {
            document.getElementById('output').innerText = 'No output received';
        }

    } catch (error) {
        console.error('Error running code:', error);
        document.getElementById('output').innerText = 'Error running code: ' + error.message;
    }
}

// Call fetchProblems when the page loads
fetchProblems();

document.addEventListener("DOMContentLoaded", function () {
    fetchProblems();

    // Initialize CodeMirror editor for the user to input code
    editor = CodeMirror(document.getElementById('editor'), {
        mode: 'javascript',
        lineNumbers: true,
        theme: 'default',
    });
});