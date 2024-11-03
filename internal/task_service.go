package internal

import "errors"
var tasks = []Task{}

var nextID = 1
var allowedCategories = map[string]bool{
    "To Do": true,
    "In Progress": true,
    "Done": true,
}

type TaskService struct{}

func (ts *TaskService) CreateTask(task Task) (Task, error) {
    if task.Name == "" {
        return Task{}, errors.New("Name is required")
    }
    if task.Complexity == "" {
        return Task{}, errors.New("Complexity is required")
    }
    if task.Category == "" {
        task.Category = "To Do"
    } else {
        if !allowedCategories[task.Category] {
            return Task{}, errors.New("Category is required")
        }
    }
    task.Id = nextID
    nextID++
    tasks = append(tasks, task)

    return task, nil
}

func (ts *TaskService) GetTasks() []Task {
    return tasks
}

func (ts *TaskService) GetTasksByCategory(category string) ([]Task, error) {
    if !allowedCategories[category] {
        return nil, errors.New("Category is invalid")
    }
    var filteredTasks []Task
    for _, task := range tasks {
        if task.Category == category {
            filteredTasks = append(filteredTasks, task)
        }
    }
    return filteredTasks, nil
}

func (ts *TaskService) UpdateTask(id int, updatedTask Task) (Task, error) {
    for i, task := range tasks {
        if task.Id == id {
            if updatedTask.Name != "" {
                tasks[i].Name = updatedTask.Name
            }
            if updatedTask.Complexity != "" {
                tasks[i].Complexity = updatedTask.Complexity
            }
            if updatedTask.Description != "" {
                tasks[i].Description = updatedTask.Description
            }
            if updatedTask.Category != "" {
                tasks[i].Category = updatedTask.Category
            }
            return tasks[i], nil
        }
    }
    return Task{}, errors.New("Task not found")
}

func (ts *TaskService) DeleteTask(id int) error {
    for i, task := range tasks {
        if task.Id == id {
            tasks = append(tasks[:i], tasks[i+1:]...)
            return nil
        }
    }

    return errors.New("Task not found")
}
