package handler

import (
    "project-t/internal"
    "net/http"
    "github.com/labstack/echo/v4"
)
var taskService = &internal.TaskService{}

func Home(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

func GetTasks(c echo.Context) error {
    tasks := taskService.GetTasks()
    return c.JSON(http.StatusOK, tasks)
}

func CreateTask(c echo.Context) error {
    var task internal.Task
    if err := c.Bind(&task); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid task data",
        })
    }

    createdTask, err := taskService.CreateTask(task)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusCreated, createdTask)
}
func DeleteTask(c echo.Context) error {
    return c.String(http.StatusOK, "Delete a task")
}
