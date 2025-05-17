package handler

import (
    "project-t/internal"
    "net/http"
    "github.com/labstack/echo/v4"
    "strconv"
	"log"
)
type Handler struct {
    TaskService *internal.TaskService
}

func NewHandler(taskService *internal.TaskService) *Handler {
    return &Handler{TaskService: taskService}
}

func Home(c echo.Context) error {
    return c.String(http.StatusOK, "Hello, World!")
}

func (h *Handler) GetTasks(c echo.Context) error {
    if category := c.QueryParam("category"); category != "" {
        tasks, err := h.TaskService.GetTasksByCategory(category)
        if err != nil {
            return c.JSON(http.StatusBadRequest, echo.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(http.StatusOK, tasks)
    }
    tasks := h.TaskService.GetTasks()
    return c.JSON(http.StatusOK, tasks)
}

func (h *Handler) CreateTask(c echo.Context) error {
    var task internal.Task
    if err := c.Bind(&task); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid task data",
        })
    }

    createdTask, err := h.TaskService.CreateTask(task)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }
    log.Printf("Task created: %+v", createdTask)
    return c.JSON(http.StatusCreated, createdTask)
}

func (h *Handler) UpdateTask (c echo.Context) error {
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
    task, err := h.TaskService.UpdateTask(id, updatedTask)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }
    return c.JSON(http.StatusOK, task)
}

func (h *Handler) DeleteTask(c echo.Context) error {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid task ID",
        })
    }
    err = h.TaskService.DeleteTask(id)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }
    return c.JSON(http.StatusOK, map[string]string{
        "message": "Task deleted successfully",
    })
}
