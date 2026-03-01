package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Order represents a single order
type Order struct {
	ID       int    `json:"id"`
	Item     string `json:"item" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gte=1"`
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	// Structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("Starting Order Service")

	// Build DB connection string
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("POSTGRES_USER", "postgres"),     // default user
		getEnv("POSTGRES_PASSWORD", "postgres"), // default password
		getEnv("POSTGRES_HOST", "localhost"),    // default host
		getEnv("POSTGRES_PORT", "5432"),         // default port
		getEnv("ORDER_DB", "orderdb"),           // default database
	)

	// Retry logic to wait for DB
	var db *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil && db.Ping() == nil {
			break
		}
		logger.Warn("DB not ready, retrying...", "attempt", i+1, "error", err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		logger.Error("Failed to connect DB", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create orders table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS orders(
			id SERIAL PRIMARY KEY,
			item TEXT NOT NULL,
			quantity INT NOT NULL CHECK (quantity > 0)
		)
	`)
	if err != nil {
		logger.Error("Table creation failed", "error", err)
		os.Exit(1)
	}

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "DOWN"})
			return
		}
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Get all orders
	r.GET("/orders", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, item, quantity FROM orders")
		if err != nil {
			logger.Error("Query failed", "error", err)
			c.JSON(500, gin.H{"error": "DB Error"})
			return
		}
		defer rows.Close()

		var orders []Order
		for rows.Next() {
			var o Order
			if err := rows.Scan(&o.ID, &o.Item, &o.Quantity); err != nil {
				logger.Error("Row scan failed", "error", err)
				c.JSON(500, gin.H{"error": "DB Error"})
				return
			}
			orders = append(orders, o)
		}
		logger.Info("Fetched orders", "count", len(orders))
		c.JSON(200, orders)
	})

	// Create a new order
	r.POST("/orders", func(c *gin.Context) {
		var o Order
		if err := c.ShouldBindJSON(&o); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}

		err := db.QueryRow(
			"INSERT INTO orders(item, quantity) VALUES($1, $2) RETURNING id",
			o.Item, o.Quantity,
		).Scan(&o.ID)
		if err != nil {
			logger.Error("Insert failed", "error", err)
			c.JSON(500, gin.H{"error": "Insert failed"})
			return
		}

		logger.Info("Order created", "id", o.ID)
		c.JSON(201, o)
	})

	logger.Info("Order Service running on port 8080")
	if err := r.Run(":8080"); err != nil {
		logger.Error("Failed to run server", "error", err)
	}
}
