package handler

import (
    "github.com/labstack/echo"
)
func Home(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

func GetTasks(c echo.Context) error {
    return c.String(http.StatusOK, "Get all tasks")
}
