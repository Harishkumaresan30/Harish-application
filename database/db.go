package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection and ensures necessary tables are created.
func InitDB() error {
	connStr := "postgres://postgres:postgres@localhost:5432/serviceweaver?sslmode=disable"
	var err error

	// Connect to the database
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Ping the database to ensure it's reachable
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Ensure tables are created
	if err := ensureTables(); err != nil {
		return fmt.Errorf("failed to ensure tables: %v", err)
	}

	return nil
}

// ensureTables checks and creates the necessary tables if they don't exist.
func ensureTables() error {
	queries := []string{
		// Create products table
		`CREATE TABLE IF NOT EXISTS public.products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			stock INTEGER NOT NULL,
			price NUMERIC(10, 2) NOT NULL
		);`,
		// Create orders table
		`CREATE TABLE IF NOT EXISTS public.orders (
			id SERIAL PRIMARY KEY,
			product_id VARCHAR(255) NOT NULL,
			quantity INTEGER NOT NULL,
			total NUMERIC(10, 2) NOT NULL,
			status VARCHAR(50) DEFAULT 'Pending' NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		// Create metrics table
		`CREATE TABLE IF NOT EXISTS public.metrics (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			value NUMERIC(10, 2) NOT NULL,
			time TIMESTAMP DEFAULT now() NOT NULL
		);`,
	}

	// Execute each query
	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %v", query, err)
		}
	}

	return nil
}
