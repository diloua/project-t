package route

import (
    "github.com/labstack/echo/v4"
    "github.com/username/project-t/handlers"
)

    func Init() *echo.Echo {
        e := echo.New()

        // Routes
        e.GET("/", handler.Home)
        e.GET("/tasks", handler.GetTasks)

        return e
    }
