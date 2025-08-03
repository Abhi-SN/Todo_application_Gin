package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Todo struct remains the same
type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// Global variable for the database connection pool
var dbpool *pgxpool.Pool

func main() {
	// Initialize the database connection and handle potential errors
	var err error
	dbpool, err = initDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	// Ensure the connection is closed when the application exits
	defer dbpool.Close()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes
	router.GET("/health", healthCheckHandler)
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodoByID)
	router.POST("/todos", createTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)

	log.Println("Starting server on port 8080")
	router.Run(":8080")
}

// Function now returns the connection pool and an error
func initDB() (*pgxpool.Pool, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// Create a new connection pool
	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping the database to ensure the connection is live
	err = pool.Ping(context.Background())
	if err != nil {
		// Close the pool if ping fails
		pool.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Successfully connected to the database.")

	// Create the 'todos' table if it doesn't already exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE
	);`
	_, err = pool.Exec(context.Background(), createTableSQL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to create table: %w", err)
	}
	log.Println("Todos table is ready.")

	// Return the initialized pool
	return pool, nil
}

// --- Handler Functions ---

func healthCheckHandler(c *gin.Context) {
	// We can even make the health check ping the DB to be more thorough
	if err := dbpool.Ping(context.Background()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "database not responding"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func getTodos(c *gin.Context) {
	rows, err := dbpool.Query(context.Background(), "SELECT id, text, completed FROM todos ORDER BY id ASC")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Text, &t.Completed); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan todo"})
			return
		}
		todos = append(todos, t)
	}

	c.IndentedJSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	var newTodoData struct {
		Text string `json:"text"`
	}
	if err := c.BindJSON(&newTodoData); err != nil {
		return
	}

	var newTodo Todo
	err := dbpool.QueryRow(context.Background(),
		"INSERT INTO todos (text, completed) VALUES ($1, $2) RETURNING id, text, completed",
		newTodoData.Text, false).Scan(&newTodo.ID, &newTodo.Text, &newTodo.Completed)

	if err != nil {
		log.Printf("Error creating todo: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodoByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var t Todo
	err = dbpool.QueryRow(context.Background(),
		"SELECT id, text, completed FROM todos WHERE id=$1", id).Scan(&t.ID, &t.Text, &t.Completed)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, t)
}

func toggleTodoStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var t Todo
	err = dbpool.QueryRow(context.Background(),
		"UPDATE todos SET completed = NOT completed WHERE id = $1 RETURNING id, text, completed", id).Scan(&t.ID, &t.Text, &t.Completed)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, t)
}

func updateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedFields struct {
		Text string `json:"text"`
	}
	if err := c.BindJSON(&updatedFields); err != nil {
		return
	}

	var t Todo
	err = dbpool.QueryRow(context.Background(),
		"UPDATE todos SET text = $1 WHERE id = $2 RETURNING id, text, completed", updatedFields.Text, id).Scan(&t.ID, &t.Text, &t.Completed)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, t)
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	res, err := dbpool.Exec(context.Background(), "DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	if res.RowsAffected() == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}
