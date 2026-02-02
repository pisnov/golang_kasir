package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgresql://postgres.btuyxzhfurruethkzqfu:Biawak123%21%40%23@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=require"

	log.Println("Attempting connection to:", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open:", err)
	}
	defer db.Close()

	fmt.Println("Testing ping...")
	if err := db.Ping(); err != nil {
		log.Fatal("Ping failed:", err)
	}

	fmt.Println("Connection successful!")

	// Try a simple query
	rows, err := db.Query("SELECT 1")
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	fmt.Println("Query successful!")
}
