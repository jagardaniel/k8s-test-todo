package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTasks(t *testing.T) {
	service := NewInMemoryTodoService()
	server := NewTodoServer(service)
	router := server.setupRouter()

	t.Run("get all tasks with empty list", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		responseTasks := getTasksFromResponse(t, response.Body)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(server.service.GetTasks()), len(responseTasks))
	})

	t.Run("get all tasks with 2 entries", func(t *testing.T) {
		server.service.CreateTask("Task 1")
		server.service.CreateTask("Task 2")

		request := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		responseTasks := getTasksFromResponse(t, response.Body)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, len(server.service.GetTasks()), len(responseTasks))
	})
}

func TestCreateTask(t *testing.T) {
	service := NewInMemoryTodoService()
	server := NewTodoServer(service)
	router := server.setupRouter()

	t.Run("create task with valid title", func(t *testing.T) {
		requestBody := strings.NewReader(`{"title": "New task"}`)
		request := httptest.NewRequest(http.MethodPost, "/api/tasks", requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.Equal(t, 1, len(server.service.GetTasks()))
	})

	t.Run("create task with empty JSON body", func(t *testing.T) {
		requestBody := strings.NewReader(`{}`)
		request := httptest.NewRequest(http.MethodPost, "/api/tasks", requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("create task with invalid title length", func(t *testing.T) {
		requestBody := strings.NewReader(`{"title": "ab"}`)
		request := httptest.NewRequest(http.MethodPost, "/api/tasks", requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestGetTask(t *testing.T) {
	service := NewInMemoryTodoService()
	server := NewTodoServer(service)
	router := server.setupRouter()

	server.service.CreateTask("Do something")
	server.service.CreateTask("Task 2")
	server.service.UpdateTask(2, "Do something else", true)

	t.Run("get task with id 1", func(t *testing.T) {
		expectedResponse := `{"id": 1, "title": "Do something", "completed": false}`

		request := newGetTaskRequest(1)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())
	})

	t.Run("get task with id 2", func(t *testing.T) {
		expectedResponse := `{"completed": true, "title": "Do something else", "id": 2}`

		request := newGetTaskRequest(2)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, expectedResponse, response.Body.String())
	})

	t.Run("get non existing task", func(t *testing.T) {
		request := newGetTaskRequest(32)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})

	t.Run("get task with invalid id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/tasks/invalid1532", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestDeleteTask(t *testing.T) {
	service := NewInMemoryTodoService()
	server := NewTodoServer(service)
	router := server.setupRouter()

	t.Run("delete task with id 1", func(t *testing.T) {
		server.service.CreateTask("Do something")

		request := newDeleteTaskRequest(1)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNoContent, response.Code)
		assert.Equal(t, 0, len(server.service.GetTasks()))
	})

	t.Run("delete non existing task", func(t *testing.T) {
		request := newDeleteTaskRequest(22)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNoContent, response.Code)
	})

	t.Run("delete task with invalid id", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/api/tasks/aewa23vea", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestUpdateTask(t *testing.T) {
	service := NewInMemoryTodoService()
	server := NewTodoServer(service)
	router := server.setupRouter()

	server.service.CreateTask("Do something")

	t.Run("update existing task with id 1", func(t *testing.T) {
		expectedTask := Task{
			ID:        1,
			Title:     "Do something else",
			Completed: true,
		}
		requestBody := strings.NewReader(`{"id": 1, "title": "Do something else", "completed": true}`)

		request := newUpdateTaskRequest(1, requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNoContent, response.Code)

		task, err := server.service.GetTask(1)
		if err != nil {
			t.Fatalf("unable to get task, '%v'", err)
		}

		assert.Equal(t, expectedTask, task)
	})

	t.Run("update existing task with id 1 without id field", func(t *testing.T) {
		expectedTask := Task{
			ID:        1,
			Title:     "Jump",
			Completed: true,
		}
		requestBody := strings.NewReader(`{"title": "Jump", "completed": true}`)

		request := newUpdateTaskRequest(1, requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNoContent, response.Code)

		task, err := server.service.GetTask(1)
		if err != nil {
			t.Fatalf("unable to get task, '%v'", err)
		}

		assert.Equal(t, expectedTask, task)
	})

	t.Run("update task with mismatching id in parameter and body", func(t *testing.T) {
		requestBody := strings.NewReader(`{"id": 24, "title": "Jump", "completed": false}`)

		request := newUpdateTaskRequest(1, requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("update task with invalid parameter id", func(t *testing.T) {
		requestBody := strings.NewReader(`{"title": "Jump", "completed": false}`)

		request := httptest.NewRequest(http.MethodPut, "/api/tasks/4242abbbb", requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("update existing task id 1 with invalid title length", func(t *testing.T) {
		requestBody := strings.NewReader(`{"id": 1, "title": "ea", "completed": false}`)

		request := newUpdateTaskRequest(1, requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("update existing task id 1 with empty body", func(t *testing.T) {
		requestBody := strings.NewReader(`{}`)

		request := newUpdateTaskRequest(1, requestBody)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func newGetTaskRequest(id int) *http.Request {
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d", id), nil)
	return request
}

func newDeleteTaskRequest(id int) *http.Request {
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/tasks/%d", id), nil)
	return request
}

func newUpdateTaskRequest(id int, body *strings.Reader) *http.Request {
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/tasks/%d", id), body)
	return request
}

func getTasksFromResponse(t testing.TB, body io.Reader) (task []Task) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&task)

	if err != nil {
		t.Fatalf("unable to parse response %q into a slice of Task, '%v'", body, err)
	}

	return
}
