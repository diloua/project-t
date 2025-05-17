package main

import (
    "project-t/router"
    "project-t/internal"
    "database/sql"
    "log"
    "os"
    "fmt"
    _ "github.com/lib/pq"  // PostgreSQL driver

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

var db *sql.DB

func initDB() *sql.DB {
    // Get database connection details from environment variables or use defaults
    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "5432")
    dbUser := getEnv("DB_USER", "postgres")
    dbPassword := getEnv("DB_PASSWORD", "postgres")
    dbName := getEnv("DB_NAME", "taskboard")
    
    // Create connection string
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)
    
    log.Printf("Connecting to PostgreSQL at %s:%s", dbHost, dbPort)
    
    // Connect to database
    database, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    
    // Check connection
    err = database.Ping()
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    
    // Create table if it doesn't exist (PostgreSQL syntax)
    createTableQuery := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        complexity TEXT NOT NULL,
        category TEXT NOT NULL
    );`
    
    _, err = database.Exec(createTableQuery)
    if err != nil {
        log.Fatalf("Error creating table: %v", err)
    }
    
    log.Println("Database initialized successfully")
    return database
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}

func main() {
    db = initDB()
    defer db.Close()
    
    taskService := &internal.TaskService{DB: db}
    e := router.Init(taskService)
    
    // Middleware
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"}, // We'll update this for production later
        AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},        
    }))

    // Get port from environment variable or use default
    port := getEnv("PORT", "8080")
    log.Printf("Starting server on port %s", port)
    e.Logger.Fatal(e.Start(":" + port))
}
