package route

import (
    "github.com/labstack/echo/v4"
    "project-t/handler"
)

    func Init() *echo.Echo {
        e := echo.New()

        // Routes
        e.GET("/", handler.Home)
        e.GET("/tasks", handler.GetTasks)
        e.POST("/tasks", handler.CreateTask)
        e.PUT("/tasks/:id", handler.UpdateTask)
        e.DELETE("/tasks/:id", handler.DeleteTask)

        return e
    }
