package main

import (
    "project-t/router"
    "project-t/internal"
    "database/sql"
    "log"
    _ "modernc.org/sqlite"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

var db *sql.DB

func initDB() *sql.DB {
    database, err := sql.Open("sqlite", "./tasks.db")
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT,
        complexity TEXT NOT NULL,
        category TEXT NOT NULL
    );`
    _, err = database.Exec(createTableQuery)
    if err != nil {
        log.Fatalf("Error creating table: %v", err)
    }
    log.Println("Database initialized")
    return database
}

func main() {
    db = initDB()
    defer db.Close()
	taskService := &internal.TaskService{DB: db}
    e := router.Init(taskService)
    // Middleware
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"http://localhost:3000"},
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
    	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},		
    }))

    e.Logger.Fatal(e.Start(":8080"))
}
