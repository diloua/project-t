package handler

import (
    "project-t/internal"
    "net/http"
    "github.com/labstack/echo/v4"
    "strconv"
)
var taskService = &internal.TaskService{}

func Home(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

func GetTasks(c echo.Context) error {
    if category := c.QueryParam("category"); category != "" {
        tasks, err := taskService.GetTasksByCategory(category)
        if err != nil {
            return c.JSON(http.StatusBadRequest, echo.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(http.StatusOK, tasks)
    }
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

func UpdateTask (c echo.Context) error {
    idParam := c.Param("id")   
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid task ID",
        })
    }
    var updatedTask internal.Task
    if err := c.Bind(&updatedTask); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid task data",
        })
    }
    task, err := taskService.UpdateTask(id, updatedTask)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }
    return c.JSON(http.StatusOK, task)
}

func DeleteTask(c echo.Context) error {
    return c.String(http.StatusOK, "Delete a task")
}
