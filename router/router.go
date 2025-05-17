package router

import (
    "github.com/labstack/echo/v4"
    "project-t/handler"
    "project-t/internal"
)

    func Init(taskService *internal.TaskService) *echo.Echo {
        e := echo.New()
        h := handler.NewHandler(taskService)
        // Routes
        e.GET("/", handler.Home)
        e.GET("/tasks", h.GetTasks)
        e.POST("/tasks", h.CreateTask)
        e.PUT("/tasks/:id", h.UpdateTask)
        e.DELETE("/tasks/:id", h.DeleteTask)

        return e
    }
