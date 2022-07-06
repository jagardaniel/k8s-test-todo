package main

/*
import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTask(t *testing.T) {
	service := NewInMemoryTodoService()
	server := NewTodoServer(service)
	router := server.setupRouter()

	t.Run("returns task with id 1", func(t *testing.T) {
		request := newGetTaskRequest(1)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		// assert.JSONEq()
	})
}

func newGetTaskRequest(id int) *http.Request {
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d", id), nil)
	return request
}
*/
