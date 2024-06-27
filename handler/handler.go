package handler

import (
    "internal/task_service"
    "net/http"
    "github.com/labstack/echo/v4"
)
func Home(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

func GetTasks(c echo.Context) error {
    return c.String(http.StatusOK, "Get all tasks")  
}

func CreateTask(c echo.Context) error {
    return c.String(http.StatusOK, "Create a new task")
}

func DeleteTask(c echo.Context) error {
    return c.String(http.StatusOK, "Delete a task")
}
