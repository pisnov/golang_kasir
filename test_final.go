package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgresql://postgres.btuyxzhfurruethkzqfu:Biawak123%21%40%23@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=require&application_name=golang_kasir"

	log.Println("Connecting to Supabase pooler...")
	log.Printf("Connection string (sanitized): postgresql://postgres.btuyxzhfurruethkzqfu:***@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres\n")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open connection:", err)
	}
	defer db.Close()

	// Set connection pool parameters
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	log.Println("Opening connection...")
	if err := db.Ping(); err != nil {
		log.Fatal("Ping failed:", err)
	}

	log.Println("✓ Connected successfully!")

	var version string
	if err := db.QueryRow("SELECT version()").Scan(&version); err != nil {
		log.Fatal("Query failed:", err)
	}

	fmt.Println("PostgreSQL version:", version)
}
