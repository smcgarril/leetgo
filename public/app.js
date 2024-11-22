let editor;

document.addEventListener("DOMContentLoaded", function () {
    fetchProblems(); // Fetch problems when the page is loaded

    // Initialize CodeMirror editor for the user to input code
    editor = CodeMirror(document.getElementById('editor'), {
        mode: 'javascript',
        lineNumbers: true,
        theme: 'default',
    });
});


// Function to fetch problems from the backend
async function fetchProblems() {
    try {
        const response = await fetch('/problems'); // Fetch problems from the Go backend (same origin)
        if (!response.ok) {
            throw new Error("Failed to fetch problems");
        }

        const problems = await response.json(); // Parse the JSON data
        renderProblems(problems); // Pass the problems to the render function
    } catch (error) {
        console.error('Error fetching problems:', error);
    }
}

// Function to render the fetched problems in the DOM as a dropdown
function renderProblems(problems) {
    const problemsDiv = document.getElementById('problems'); // The container for the dropdown
    problemsDiv.innerHTML = ''; // Clear previous content

    // Create the dropdown element
    const select = document.createElement('select');
    select.id = 'problem-select';
    select.onchange = () => displayProblemDescription(problems); // Attach onchange event

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
        option.value = problem.id; // Use problem ID as the value
        option.innerText = problem.title; // Display the title
        select.appendChild(option);
    });

    problemsDiv.appendChild(select);
}

// Function to display the selected problem's description
function displayProblemDescription(problems) {
    const select = document.getElementById('problem-select');
    const selectedProblemId = select.value;

    const problem = problems.find(p => p.id === selectedProblemId); // Find the selected problem
    const descriptionDiv = document.getElementById('problem-description');
    if (problem) {
        descriptionDiv.innerHTML = `<h3>${problem.title}</h3><p>${problem.description}</p>`;
    } else {
        descriptionDiv.innerHTML = '';
    }
}

// Call fetchProblems when the page loads
fetchProblems();

async function runCode() {
    const code = editor.getValue(); // Get the code from the editor
    try {
        const response = await fetch('/execute', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ code }), // Send the code to the backend
        });

        if (!response.ok) {
            throw new Error('Failed to execute code');
        }

        const rawText = await response.text();
        const data = JSON.parse(rawText);

        // Check if `data.output` exists and render it
        if (data.output) {
            document.getElementById('output').innerText = data.output.trim(); 
        } else {
            document.getElementById('output').innerText = 'No output received';
        }
    } catch (error) {
        console.error('Error running code:', error);
        document.getElementById('output').innerText = 'Error running code: ' + error.message;
    }
}
