package main

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

var session *gocql.Session

func init() {
	cluster := gocql.NewCluster("scylla")
	cluster.Keyspace = "todo"
	cluster.Consistency = gocql.Quorum
	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("unable to connect to ScyllaDB: %v", err)
	}
}

type TodoItem struct {
	ID          gocql.UUID `json:"id"`
	UserID      gocql.UUID `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

func createTodoHandler(c *gin.Context) {
	var todo TodoItem
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	todo.ID = gocql.TimeUUID()
	todo.Created = time.Now()
	todo.Updated = time.Now()

	if err := session.Query(`
        INSERT INTO items (id, user_id, title, description, status, created, updated)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.Created, todo.Updated).Exec(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, todo)
}

func readTodoHandler(c *gin.Context) {
	id := c.Param("id")
	var todo TodoItem
	if err := session.Query(`
        SELECT id, user_id, title, description, status, created, updated FROM items WHERE id = ?`,
		id).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated); err != nil {
		c.JSON(404, gin.H{"error": "TODO item not found"})
		return
	}
	c.JSON(200, todo)
}

func updateTodoHandler(c *gin.Context) {
	id := c.Param("id")
	var todo TodoItem
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	todo.Updated = time.Now()

	if err := session.Query(`
        UPDATE items SET title = ?, description = ?, status = ?, updated = ? WHERE id = ?`,
		todo.Title, todo.Description, todo.Status, todo.Updated, id).Exec(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "TODO item updated"})
}

func deleteTodoHandler(c *gin.Context) {
	id := c.Param("id")
	if err := session.Query(`DELETE FROM items WHERE id = ?`, id).Exec(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "TODO item deleted"})
}

func listTodoHandler(c *gin.Context) {
	userID := c.Query("user_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var todos []TodoItem
	var todo TodoItem
	var query string
	var totalTodos int

	if status != "" {
		query = "SELECT COUNT(*) FROM items WHERE user_id = ? AND status = ? ALLOW FILTERING"
		if err := session.Query(query, userID, status).Scan(&totalTodos); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		query = "SELECT id, user_id, title, description, status, created, updated FROM items WHERE user_id = ? AND status = ? ALLOW FILTERING"
		iter := session.Query(query, userID, status).Iter()
		defer iter.Close()

		for iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
			todos = append(todos, todo)
		}
	} else {
		query = "SELECT COUNT(*) FROM items WHERE user_id = ? ALLOW FILTERING"
		if err := session.Query(query, userID).Scan(&totalTodos); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		query = "SELECT id, user_id, title, description, status, created, updated FROM items WHERE user_id = ? ALLOW FILTERING"
		iter := session.Query(query, userID).Iter()
		defer iter.Close()

		for iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
			todos = append(todos, todo)
		}
	}

	totalPages := int(math.Ceil(float64(totalTodos) / float64(pageSize)))

	var nextPage, prevPage int
	if page < totalPages {
		nextPage = page + 1
	} else {
		nextPage = -1
	}

	if page > 1 {
		prevPage = page - 1
	} else {
		prevPage = -1
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if endIndex > totalTodos {
		endIndex = totalTodos
	}

	currentPageTodos := todos[startIndex:endIndex]

	c.JSON(200, gin.H{
		"todos":      currentPageTodos,
		"page":       page,
		"totalPages": totalPages,
		"nextPage":   nextPage,
		"prevPage":   prevPage,
	})
}

func main() {
	defer session.Close()

	router := gin.Default()

	router.POST("/todos", createTodoHandler)
	router.GET("/todos/:id", readTodoHandler)
	router.PUT("/todos/:id", updateTodoHandler)
	router.DELETE("/todos/:id", deleteTodoHandler)
	router.GET("/todos", listTodoHandler)

	router.Run(":8080")
}
