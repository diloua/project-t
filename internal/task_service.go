package internal

import "errors"
var tasks = []Task{}

var nextID = 1
type TaskService struct{}

func (ts *TaskService) CreateTask(task Task) (Task, error) {
    if task.Name == "" {
        return Task{}, errors.New("Name is required")
    }
    if task.Complexity == "" {
        return Task{}, errors.New("Complexity is required")
    }
    task.Id = nextID
    task.Status = "Pending"
    nextID++
    tasks = append(tasks, task)

    return task, nil
}

func (ts *TaskService) GetTasks() []Task {
    return tasks
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
