package internal

import "errors"
import "database/sql"
import "fmt"
import "log"

var tasks = []Task{}

var nextID = 1
var allowedCategories = map[string]bool{
    "To Do": true,
    "In Progress": true,
    "Done": true,
}

type TaskService struct{
    DB *sql.DB
}

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
            return Task{}, errors.New("Category is invalid")
        }
    }
    
    // Use standard SQL insert without RETURNING clause
    query := `
    INSERT INTO tasks (name, description, complexity, category)
    VALUES (?, ?, ?, ?)
    `
    result, err := ts.DB.Exec(query, task.Name, task.Description, task.Complexity, task.Category)
    if err != nil {
        return Task{}, err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return Task{}, err
    }
    log.Printf("Task inserted with ID: %d", id)   
    task.Id = int(id)
    
    return task, nil
}

func (ts *TaskService) GetTasks() []Task {
    var tasks []Task
    
    rows, err := ts.DB.Query("SELECT id, name, description, complexity, category FROM tasks")
    if err != nil {
        // Log the error but return an empty slice rather than failing
        fmt.Println("Error fetching tasks:", err)
        return tasks
    }
    defer rows.Close()
    
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.Id, &task.Name, &task.Description, &task.Complexity, &task.Category)
        if err != nil {
            // Log the error but continue with other rows
            fmt.Println("Error scanning task row:", err)
            continue
        }
        tasks = append(tasks, task)
    }
    log.Printf("Retrieved %d tasks from database", len(tasks))
    
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
