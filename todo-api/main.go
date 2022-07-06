package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoService interface {
	GetTasks() []Task
	GetTask(id int) (Task, error)
	CreateTask(title string) int
	DeleteTask(id int) error
	UpdateTask(id int, title string, completed bool) error
}

type todoServer struct {
	service TodoService
}

var ErrNotFound = errors.New("resource was not found")

func NewTodoServer(service TodoService) *todoServer {
	return &todoServer{service: service}
}

func (ts *todoServer) getTasks(c *gin.Context) {
	tasks := ts.service.GetTasks()
	c.JSON(http.StatusOK, tasks)
}

func (ts *todoServer) getTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := ts.service.GetTask(id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ts *todoServer) createTask(c *gin.Context) {
	type TaskInput struct {
		Title string `json:"title" binding:"required,gte=3,lte=50"`
	}

	var ti TaskInput
	if err := c.ShouldBindJSON(&ti); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ts.service.CreateTask(ti.Title)

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (ts *todoServer) deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ignore if error is ErrNotFound since we want to return 204 even
	// if the task does not exist
	err = ts.service.DeleteTask(id)
	if err != nil && !errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (ts *todoServer) updateTask(c *gin.Context) {
	type TaskInput struct {
		ID        int    `json:"id"` // Not required, but need to match URL id if set
		Title     string `json:"title" binding:"required,gte=3,lte=50"`
		Completed bool   `json:"completed"` // Will be false if not set
	}

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ti TaskInput
	if err := c.ShouldBindJSON(&ti); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ti.ID != 0 {
		if ti.ID != id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id in body does not match URL id"})
			return
		}
	}

	// Return error if the task does not exist
	err = ts.service.UpdateTask(id, ti.Title, ti.Completed)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (ts *todoServer) setupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	api.GET("/tasks", ts.getTasks)
	api.POST("/tasks", ts.createTask)
	api.GET("/tasks/:id", ts.getTask)
	api.DELETE("/tasks/:id", ts.deleteTask)
	api.PUT("/tasks/:id", ts.updateTask)

	return router
}

func main() {
	service := NewInMemoryTodoService()
	server := NewTodoServer(service)
	router := server.setupRouter()

	router.Run(":8000")
}
