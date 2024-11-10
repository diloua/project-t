document.addEventListener('DOMContentLoaded', () => {
    const apiBase = 'http://localhost:8080';

    const taskModal = document.getElementById('taskModal');
    const createTaskBtn = document.getElementById('createTaskBtn');
    const closeModal = document.querySelector('.close');
    const taskForm = document.getElementById('taskForm');

    const todoTasks = document.getElementById('todoTasks');
    const inProgressTasks = document.getElementById('inProgressTasks');
    const reviewTasks = document.getElementById('reviewTasks');
    const doneTasks = document.getElementById('doneTasks');

    // Open Modal
    createTaskBtn.onclick = () => {
        taskModal.style.display = 'block';
    };

    // Close Modal
    closeModal.onclick = () => {
        taskModal.style.display = 'none';
    };

    window.onclick = (event) => {
        if (event.target == taskModal) {
            taskModal.style.display = 'none';
        }
    };

    // Fetch and Display Tasks
    const fetchTasks = async () => {
        try {
            const response = await fetch(`${apiBase}/tasks`);
            const tasks = await response.json();
            renderTasks(tasks);
        } catch (error) {
            console.error('Error fetching tasks:', error);
        }
    };

    const renderTasks = (tasks) => {
        // Clear existing tasks
        todoTasks.innerHTML = '';
        inProgressTasks.innerHTML = '';
        reviewTasks.innerHTML = '';
        doneTasks.innerHTML = '';

        tasks.forEach(task => {
            const taskElement = document.createElement('div');
            taskElement.className = 'task';
            taskElement.draggable = true;
            taskElement.dataset.id = task.id;

            taskElement.innerHTML = `
                <h3>${task.name}</h3>
                <p>${task.description}</p>
                <p><strong>Complexity:</strong> ${task.complexity}</p>
                <button onclick="deleteTask(${task.id})">Delete</button>
            `;

            // Add drag event listeners
            taskElement.addEventListener('dragstart', dragStart);
            taskElement.addEventListener('dragend', dragEnd);

            switch (task.category) {
                case 'To Do':
                    todoTasks.appendChild(taskElement);
                    break;
                case 'In Progress':
                    inProgressTasks.appendChild(taskElement);
                    break;
                case 'Review':
                    reviewTasks.appendChild(taskElement);
                    break;
                case 'Done':
                    doneTasks.appendChild(taskElement);
                    break;
                default:
                    break;
            }
        });
    };

    // Create Task
    taskForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(taskForm);
        const taskData = {
            name: formData.get('name'),
            description: formData.get('description'),
            complexity: formData.get('complexity'),
            category: formData.get('category')
        };

        try {
            const response = await fetch(`${apiBase}/tasks`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(taskData)
            });

            const result = await response.json();

            if (response.ok) {
                taskForm.reset();
                taskModal.style.display = 'none';
                fetchTasks();
            } else {
                alert(`Error: ${result.error}`);
            }
        } catch (error) {
            console.error('Error creating task:', error);
        }
    });

    // Delete Task
    window.deleteTask = async (id) => {
        if (!confirm('Are you sure you want to delete this task?')) return;

        try {
            const response = await fetch(`${apiBase}/tasks/${id}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                fetchTasks();
            } else {
                const result = await response.json();
                alert(`Error: ${result.error}`);
            }
        } catch (error) {
            console.error('Error deleting task:', error);
        }
    };

    // Drag and Drop Handlers
    let draggedTask = null;

    const dragStart = (e) => {
        draggedTask = e.target;
        setTimeout(() => {
            e.target.style.display = 'none';
        }, 0);
    };

    const dragEnd = (e) => {
        setTimeout(() => {
            e.target.style.display = 'block';
            draggedTask = null;
        }, 0);
    };

    const columns = document.querySelectorAll('.tasks');

    columns.forEach(column => {
        column.addEventListener('dragover', (e) => {
            e.preventDefault();
        });

        column.addEventListener('drop', async (e) => {
            e.preventDefault();
            if (draggedTask) {
                const taskId = draggedTask.dataset.id;
                const newCategory = e.currentTarget.parentElement.dataset.category;

                try {
                    const response = await fetch(`${apiBase}/tasks/${taskId}`, {
                        method: 'PUT',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            name: draggedTask.querySelector('h3').innerText,
                            description: draggedTask.querySelector('p').innerText,
                            complexity: draggedTask.querySelectorAll('p')[1].innerText.replace('Complexity: ', ''),
                            category: newCategory
                        })
                    });

                    const result = await response.json();

                    if (response.ok) {
                        fetchTasks();
                    } else {
                        alert(`Error: ${result.error}`);
                    }
                } catch (error) {
                    console.error('Error updating task category:', error);
                }
            }
        });
    });

    // Initial Fetch
    fetchTasks();
});
