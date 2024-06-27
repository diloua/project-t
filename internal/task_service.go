package internal

import "errors"

type Task struct {
    ID          int
    Name        string
    Description string
    Status      string
}

var tasks = []Task{}

var nextID = 1

func CreateTask(name string, description string) (Task, error) {
    if name == "" {
        return Task{}, errors.New("Name is required")
    }

    task := Task{
        ID:          nextID,
        Name:        name,
        Description: description,
        Status:      "Pending",
    }

    nextID++
    tasks = append(tasks, task)

    return task, nil
}

func GetTasks() []Task {
    return tasks
}

func DeleteTask(id int) error {
    for i, task := range tasks {
        if task.ID == id {
            tasks = append(tasks[:i], tasks[i+1:]...)
            return nil
        }
    }

    return errors.New("Task not found")
}
