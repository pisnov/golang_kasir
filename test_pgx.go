package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	connStr := "postgresql://postgres.btuyxzhfurruethkzqfu:Biawak123%21%40%23@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=require"

	log.Println("Connecting with pgx driver...")

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close(context.Background())

	log.Println("✓ Connected successfully!")

	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatal("Query failed:", err)
	}

	fmt.Println("PostgreSQL version:", version)
}
