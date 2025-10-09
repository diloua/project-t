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
        e.GET("/api/tasks", h.GetTasks)
        e.POST("/api/tasks", h.CreateTask)
        e.PUT("/api/tasks/:id", h.UpdateTask)
        e.DELETE("/api/tasks/:id", h.DeleteTask)

        return e
    }
