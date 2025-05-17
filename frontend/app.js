// app.js

document.addEventListener('DOMContentLoaded', () => {
    fetchTasks();
});

function fetchTasks() {
    fetch('http://localhost:8080/tasks')
        .then(response => response.json())
        .then(data => {
            console.log("Tasks received TEST:", data); // Log the data
            displayTasks(data);
        })
        .catch(error => {
            console.error('Error fetching tasks:', error);
        });
}

function displayTasks(tasks) {
    // Clear existing tasks
    document.getElementById('todo-column').innerHTML = '';
    document.getElementById('inprogress-column').innerHTML = '';
    document.getElementById('done-column').innerHTML = '';

    tasks.forEach(task => {
        const taskElement = createTaskElement(task);

        // Append task to the appropriate column based on its category
        if (task.category === 'To Do') {
            document.getElementById('todo-column').appendChild(taskElement);
        } else if (task.category === 'In Progress') {
            document.getElementById('inprogress-column').appendChild(taskElement);
        } else if (task.category === 'Done') {
            document.getElementById('done-column').appendChild(taskElement);
        }
    });
}

function createTaskElement(task) {
    const taskDiv = document.createElement('div');
    taskDiv.className = 'task';
    taskDiv.draggable = true; // Enable dragging
    taskDiv.innerHTML = `
        <h5>${task.name}</h5>
        <p>${task.description || ''}</p>
        <small>Complexity: ${task.complexity}</small>
        <button class="btn btn-danger btn-sm delete-task" title="Delete Task">Delete</button>

    `;
    taskDiv.querySelector('.delete-task').addEventListener('click', function(e) {
        e.stopPropagation(); // Prevent triggering other events
        deleteTask(task.id, taskDiv);
    });
    // Set data attribute for task ID
    taskDiv.setAttribute('data-task-id', task.id);
    taskDiv.addEventListener('dragstart', drag);
    return taskDiv;
}

// Add event listener for form submission
document.getElementById('task-form').addEventListener('submit', function (e) {
    e.preventDefault(); // Prevent default form submission

    // Collect form data
    const name = document.getElementById('task-name').value.trim();
    const description = document.getElementById('task-description').value.trim();
    const complexity = document.getElementById('task-complexity').value;

    // Create task object
    const newTask = {
        Name: name,
        Description: description,
        Complexity: complexity,
    };

    // Send POST request to backend
    fetch('http://localhost:8080/tasks', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(newTask),
    })
        .then((response) => {
            if (!response.ok) {
                return response.json().then((err) => {
                    throw err;
                });
            }
            return response.json();
        })
        .then((createdTask) => {
            // Clear form fields
            document.getElementById('task-form').reset();

            // Add the new task to the UI
            console.log('Created Task:', createdTask); // Add this line
            const taskElement = createTaskElement(createdTask);
            document.getElementById('todo-column').appendChild(taskElement);
        })
        .catch((error) => {
            console.error('Error creating task:', error);
            alert('Error creating task: ' + error.error);
        });
});

function allowDrop(event) {
        event.preventDefault();
        event.currentTarget.classList.add('drag-over');
}


function dragLeave(event) {
    event.currentTarget.classList.remove('drag-over');
}

function drag(event) {
    event.dataTransfer.setData('text/plain', event.target.getAttribute('data-task-id'));
}

function drop(event, newCategory) {
    event.preventDefault();
    event.currentTarget.classList.remove('drag-over');
    const taskId = event.dataTransfer.getData('text/plain');
    const taskElement = document.querySelector(`.task[data-task-id='${taskId}']`);

   // Move the task element to the new column
    event.currentTarget.appendChild(taskElement);

        // Update the task's category in the backend
    updateTaskCategory(taskId, newCategory);
}

function updateTaskCategory(taskId, newCategory) {
        const updatedTask = {
                    category: newCategory
                };

        fetch(`http://localhost:8080/tasks/${taskId}`, {
                    method: 'PUT',
                    headers: {
                                    'Content-Type': 'application/json',
                                },
                    body: JSON.stringify(updatedTask),
                })
        .then(response => {
                    if (!response.ok) {
                                    return response.json().then(err => { throw err; });
                                }
                    return response.json();
                })
        .then(data => {
                    console.log('Task updated:', data);
                })
        .catch(error => {
                    console.error('Error updating task:', error);
                    alert('Error updating task: ' + error.error);
                });
}

function deleteTask(taskId, taskElement) {
    if (!confirm('Are you sure you want to delete this task?')) {
        return;
    }

    fetch(`http://localhost:8080/tasks/${taskId}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (!response.ok) {
            return response.json().then(err => { throw err; });
        }
        return response.json();
    })
    .then(() => {
        // Remove the task element from the UI
        taskElement.remove();
    })
    .catch(error => {
        console.error('Error deleting task:', error);
        alert('Error deleting task: ' + (error.error || 'Unknown error'));
    });
}
