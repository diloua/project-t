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
	rows, err := ts.DB.Query("SELECT id, name, description, complexity, category FROM tasks WHERE category = ?", category)
	if err != nil {
		return nil, err
	}
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
	err := ts.DB.QueryRow("SELECT id, name, description, complexity, category FROM tasks WHERE id = ?", id).Scan(&existingTask.Id, &existingTask.Name, &existingTask.Description, &existingTask.Complexity, &existingTask.Category)
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

	_, err = ts.DB.Exec("UPDATE tasks SET name = ?, description = ?, complexity = ?, category = ? WHERE id = ?", existingTask.Name, existingTask.Description, existingTask.Complexity, existingTask.Category, id)

	if err != nil {
		return Task{}, err
	}
	log.Printf("Task with ID %d updated successfully", id)
	return existingTask, nil
}


func (ts *TaskService) DeleteTask(id int) error {
    // Execute DELETE query against the database
    result, err := ts.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
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

