// app.js - Fixed Version

document.addEventListener('DOMContentLoaded', () => {
    // Set up form toggle
    const formContainer = document.getElementById('task-form-container');
    const showFormButton = document.getElementById('show-task-form');
    const closeFormButton = document.getElementById('close-form');
    const cancelButton = document.getElementById('cancel-task');
    
    showFormButton.addEventListener('click', () => {
        formContainer.classList.add('visible');
    });
    
    closeFormButton.addEventListener('click', () => {
        formContainer.classList.remove('visible');
    });
    
    cancelButton.addEventListener('click', () => {
        formContainer.classList.remove('visible');
    });
    
    // Load tasks from server
    fetchTasks();
});

function fetchTasks() {
    // Use the full URL with the correct port
    fetch('http://localhost:8080/tasks')
        .then(response => {
            if (!response.ok) {
                throw new Error(`Server returned ${response.status}: ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            console.log("Tasks received:", data);
            displayTasks(data);
            updateTaskCounts(data);
        })
        .catch(error => {
            console.error('Error fetching tasks:', error);
            showAlert('Could not load tasks. Please check if the server is running.', 'error');
        });
}

function displayTasks(tasks) {
    // Get column content divs
    const todoColumn = document.querySelector('#todo-column .column-content');
    const inProgressColumn = document.querySelector('#inprogress-column .column-content');
    const doneColumn = document.querySelector('#done-column .column-content');
    
    // Clear existing tasks (keeping the add button and empty messages)
    const todoAddButton = document.getElementById('show-task-form');
    const todoEmptyMessage = document.getElementById('todo-empty-message');
    const inProgressEmptyMessage = document.getElementById('inprogress-empty-message');
    const doneEmptyMessage = document.getElementById('done-empty-message');
    
    // Remove all tasks
    document.querySelectorAll('.task').forEach(el => el.remove());
    
    // Show/hide empty messages based on tasks
    const todoTasks = tasks.filter(task => task.category === 'To Do');
    const inProgressTasks = tasks.filter(task => task.category === 'In Progress');
    const doneTasks = tasks.filter(task => task.category === 'Done');
    
    todoEmptyMessage.style.display = todoTasks.length > 0 ? 'none' : 'block';
    inProgressEmptyMessage.style.display = inProgressTasks.length > 0 ? 'none' : 'block';
    doneEmptyMessage.style.display = doneTasks.length > 0 ? 'none' : 'block';
    
    // Add tasks to appropriate columns
    tasks.forEach(task => {
        const taskElement = createTaskElement(task);
        
        if (task.category === 'To Do') {
            // Insert before the Add Task button
            todoColumn.insertBefore(taskElement, todoAddButton);
        } else if (task.category === 'In Progress') {
            inProgressColumn.appendChild(taskElement);
        } else if (task.category === 'Done') {
            doneColumn.appendChild(taskElement);
        }
    });
    
    updateTaskCounts(tasks);
}

function createTaskElement(task) {
    const taskDiv = document.createElement('div');
    taskDiv.className = 'task';
    taskDiv.draggable = true;
    
    // Add color-coded border based on category
    if (task.category === 'To Do') {
        taskDiv.classList.add('todo-task');
    } else if (task.category === 'In Progress') {
        taskDiv.classList.add('progress-task');
    } else if (task.category === 'Done') {
        taskDiv.classList.add('done-task');
    }
    
    // Create task content
    taskDiv.innerHTML = `
        <h5>${task.name}</h5>
        <p>${task.description || 'No description provided'}</p>
        <div class="task-footer">
            <span class="complexity-badge complexity-${task.complexity}">${task.complexity}</span>
            <span class="task-id">ID: ${task.id}</span>
        </div>
        <button class="delete-task" title="Delete Task">
            <i class="fas fa-trash"></i>
        </button>
    `;
    
    // Set data attribute for task ID
    taskDiv.setAttribute('data-task-id', task.id);
    
    // Add event listeners
    taskDiv.addEventListener('dragstart', drag);
    taskDiv.querySelector('.delete-task').addEventListener('click', function(e) {
        e.stopPropagation(); // Prevent triggering other events
        deleteTask(task.id, taskDiv);
    });
    
    return taskDiv;
}

function updateTaskCounts(tasks) {
    const todoCount = tasks.filter(task => task.category === 'To Do').length;
    const inProgressCount = tasks.filter(task => task.category === 'In Progress').length;
    const doneCount = tasks.filter(task => task.category === 'Done').length;
    
    document.getElementById('todo-count').textContent = todoCount;
    document.getElementById('inprogress-count').textContent = inProgressCount;
    document.getElementById('done-count').textContent = doneCount;
}

// Form submission
document.getElementById('task-form').addEventListener('submit', function(e) {
    e.preventDefault();
    
    const name = document.getElementById('task-name').value.trim();
    const description = document.getElementById('task-description').value.trim();
    const complexity = document.getElementById('task-complexity').value;
    
    // Create task object - Use capitalized field names to match your backend
    const newTask = {
        Name: name,
        Description: description,
        Complexity: complexity,
        Category: 'To Do'
    };
    
    // Show loading state
    const submitButton = this.querySelector('button[type="submit"]');
    const originalButtonText = submitButton.innerHTML;
    submitButton.disabled = true;
    submitButton.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Creating...';
    
    // Send POST request to backend with full URL
    fetch('http://localhost:8080/tasks', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(newTask),
    })
    .then(response => {
        submitButton.disabled = false;
        submitButton.innerHTML = originalButtonText;
        
        if (!response.ok) {
            return response.json().then(err => { throw err; });
        }
        return response.json();
    })
    .then(createdTask => {
        // Reset form
        document.getElementById('task-form').reset();
        document.getElementById('task-form-container').classList.remove('visible');
        
        // Refresh tasks
        fetchTasks();
        
        // Show success message
        showAlert('Task created successfully!', 'success');
    })
    .catch(error => {
        console.error('Error creating task:', error);
        showAlert('Error creating task: ' + (error.error || 'Unknown error'), 'error');
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
    
    // Update the task's category in the backend
    updateTaskCategory(taskId, newCategory);
}

function updateTaskCategory(taskId, newCategory) {
    const updatedTask = {
        Category: newCategory  // Note: Using capital 'C' to match backend
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
        
        // Refresh tasks to ensure UI is consistent with backend
        fetchTasks();
        
        // Show success message
        showAlert(`Task moved to ${newCategory}`, 'success');
    })
    .catch(error => {
        console.error('Error updating task:', error);
        showAlert('Error updating task: ' + (error.error || 'Unknown error'), 'error');
        
        // Refresh tasks to ensure UI is consistent with backend
        fetchTasks();
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
        // Refresh tasks
        fetchTasks();
        
        // Show success message
        showAlert('Task deleted successfully!', 'success');
    })
    .catch(error => {
        console.error('Error deleting task:', error);
        showAlert('Error deleting task: ' + (error.error || 'Unknown error'), 'error');
    });
}

function showAlert(message, type) {
    // Create alert element
    const alertDiv = document.createElement('div');
    alertDiv.className = type === 'success' ? 'alert alert-success' : 'alert alert-danger';
    
    // Add content
    alertDiv.innerHTML = `
        <i class="fas ${type === 'success' ? 'fa-check-circle' : 'fa-exclamation-circle'}"></i>
        <div>
            <strong>${type === 'success' ? 'Success' : 'Error'}</strong>
            <p>${message}</p>
        </div>
        <button type="button" class="close ml-auto" aria-label="Close">
            <span aria-hidden="true">&times;</span>
        </button>
    `;
    
    // Add to document
    document.body.appendChild(alertDiv);
    
    // Add click event to close button
    alertDiv.querySelector('.close').addEventListener('click', () => {
        alertDiv.remove();
    });
    
    // Auto remove after a few seconds
    setTimeout(() => {
        if (alertDiv.parentNode) {
            alertDiv.remove();
        }
    }, type === 'success' ? 3000 : 5000);
}
