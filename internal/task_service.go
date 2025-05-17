package internal

import (
    "database/sql"
    "errors"
    "fmt"
    "log"
)

var allowedCategories = map[string]bool{
    "To Do":       true,
    "In Progress": true,
    "Done":        true,
}

type TaskService struct {
    DB *sql.DB
}

func (ts *TaskService) CreateTask(task Task) (Task, error) {
    // Validate task
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
    
    // Use PostgreSQL's RETURNING clause to get the ID
    var id int
    query := `
    INSERT INTO tasks (name, description, complexity, category)
    VALUES ($1, $2, $3, $4)
    RETURNING id
    `
    err := ts.DB.QueryRow(query, task.Name, task.Description, task.Complexity, task.Category).Scan(&id)
    if err != nil {
        return Task{}, err
    }
    
    log.Printf("Task inserted with ID: %d", id)
    task.Id = id
    
    return task, nil
}

func (ts *TaskService) GetTasks() []Task {
    var tasks []Task
    
    rows, err := ts.DB.Query("SELECT id, name, description, complexity, category FROM tasks")
    if err != nil {
        fmt.Println("Error fetching tasks:", err)
        return tasks
    }
    defer rows.Close()
    
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.Id, &task.Name, &task.Description, &task.Complexity, &task.Category)
        if err != nil {
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

    var tasks []Task
    rows, err := ts.DB.Query("SELECT id, name, description, complexity, category FROM tasks WHERE category = $1", category)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.Id, &task.Name, &task.Description, &task.Complexity, &task.Category)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }
    log.Printf("Retrieved %d tasks from category %s", len(tasks), category)
    return tasks, nil
}

func (ts *TaskService) UpdateTask(id int, updatedTask Task) (Task, error) {
    var existingTask Task
    err := ts.DB.QueryRow("SELECT id, name, description, complexity, category FROM tasks WHERE id = $1", id).Scan(
        &existingTask.Id, &existingTask.Name, &existingTask.Description, &existingTask.Complexity, &existingTask.Category)
    if err != nil {
        if err == sql.ErrNoRows {
            return Task{}, errors.New("Task not found")
        }
        return Task{}, err
    }
    
    if updatedTask.Name != "" {
        existingTask.Name = updatedTask.Name
    }
    if updatedTask.Complexity != "" {
        existingTask.Complexity = updatedTask.Complexity
    }
    if updatedTask.Description != "" {
        existingTask.Description = updatedTask.Description
    }
    if updatedTask.Category != "" {
        if !allowedCategories[updatedTask.Category] {
            return Task{}, errors.New("Category is invalid")
        }
        existingTask.Category = updatedTask.Category
    }

    _, err = ts.DB.Exec(
        "UPDATE tasks SET name = $1, description = $2, complexity = $3, category = $4 WHERE id = $5",
        existingTask.Name, existingTask.Description, existingTask.Complexity, existingTask.Category, id)

    if err != nil {
        return Task{}, err
    }
    log.Printf("Task with ID %d updated successfully", id)
    return existingTask, nil
}

func (ts *TaskService) DeleteTask(id int) error {
    // Execute DELETE query against the database
    result, err := ts.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
    if err != nil {
        log.Printf("Error deleting task: %v", err)
        return err
    }
    
    // Check if any rows were affected by the delete operation
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error getting rows affected: %v", err)
        return err
    }
    
    // If no rows were affected, the task wasn't found
    if rowsAffected == 0 {
        return errors.New("Task not found")
    }
    
    log.Printf("Task with ID %d deleted successfully", id)
    return nil
}
