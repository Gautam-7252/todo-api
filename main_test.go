package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/todos", createTodoHandler)
	r.GET("/todos/:id", readTodoHandler)
	r.PUT("/todos/:id", updateTodoHandler)
	r.DELETE("/todos/:id", deleteTodoHandler)
	r.GET("/todos", listTodoHandler)
	return r
}

func TestCreateTodoHandler(t *testing.T) {
	r := SetupRouter()

	todo := TodoItem{
		UserID:      gocql.TimeUUID(),
		Title:       "Test TODO",
		Description: "Test Description",
		Status:      "pending",
	}
	jsonValue, _ := json.Marshal(todo)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestReadTodoHandler(t *testing.T) {
	r := SetupRouter()

	id := gocql.TimeUUID()
	todo := TodoItem{
		ID:          id,
		UserID:      gocql.TimeUUID(),
		Title:       "Test TODO",
		Description: "Test Description",
		Status:      "pending",
		Created:     time.Now(),
		Updated:     time.Now(),
	}
	session.Query(`
        INSERT INTO items (id, user_id, title, description, status, created, updated)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated).Exec()

	req, _ := http.NewRequest("GET", "/todos/"+id.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestUpdateTodoHandler(t *testing.T) {
	r := SetupRouter()

	id := gocql.TimeUUID()
	todo := TodoItem{
		ID:          id,
		UserID:      gocql.TimeUUID(),
		Title:       "Test TODO",
		Description: "Test Description",
		Status:      "pending",
		Created:     time.Now(),
		Updated:     time.Now(),
	}
	session.Query(`
        INSERT INTO items (id, user_id, title, description, status, created, updated)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated).Exec()

	updatedTodo := TodoItem{
		Title:       "Updated Test TODO",
		Description: "Updated Description",
		Status:      "completed",
	}
	jsonValue, _ := json.Marshal(updatedTodo)

	req, _ := http.NewRequest("PUT", "/todos/"+id.String(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestDeleteTodoHandler(t *testing.T) {
	r := SetupRouter()

	id := gocql.TimeUUID()
	todo := TodoItem{
		ID:          id,
		UserID:      gocql.TimeUUID(),
		Title:       "Test TODO",
		Description: "Test Description",
		Status:      "pending",
		Created:     time.Now(),
		Updated:     time.Now(),
	}
	session.Query(`
        INSERT INTO items (id, user_id, title, description, status, created, updated)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated).Exec()

	req, _ := http.NewRequest("DELETE", "/todos/"+id.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestListTodoHandler(t *testing.T) {
	r := SetupRouter()

	for i := 0; i < 5; i++ {
		todo := TodoItem{
			ID:          gocql.TimeUUID(),
			UserID:      gocql.TimeUUID(),
			Title:       "Test TODO " + strconv.Itoa(i),
			Description: "Test Description " + strconv.Itoa(i),
			Status:      "pending",
			Created:     time.Now(),
			Updated:     time.Now(),
		}
		session.Query(`
            INSERT INTO items (id, user_id, title, description, status, created, updated)
            VALUES (?, ?, ?, ?, ?, ?, ?)`,
			todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated).Exec()
	}

	req, _ := http.NewRequest("GET", "/todos?user_id=b4ff8577-0b4f-4033-8ab0-3d1b2e4a7f25&page=1&limit=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
