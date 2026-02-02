package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Test with simple connection string format
	tests := []struct {
		name    string
		connStr string
	}{
		{
			"psql format - password with special chars raw",
			"host=aws-1-ap-southeast-1.pooler.supabase.com port=6543 user=postgres.btuyxzhfurruethkzqfu password='Biawak123!@#' dbname=postgres sslmode=disable",
		},
		{
			"psql format - password escaped",
			"host=aws-1-ap-southeast-1.pooler.supabase.com port=6543 user=postgres.btuyxzhfurruethkzqfu password=Biawak123%21%40%23 dbname=postgres sslmode=disable",
		},
		{
			"url format - password percent encoded",
			"postgresql://postgres.btuyxzhfurruethkzqfu:Biawak123%21%40%23@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=disable",
		},
	}

	for _, test := range tests {
		log.Printf("Testing: %s\n", test.name)
		db, err := sql.Open("postgres", test.connStr)
		if err != nil {
			log.Printf("  Open failed: %v\n", err)
			continue
		}

		if err := db.Ping(); err != nil {
			log.Printf("  Ping failed: %v\n", err)
			db.Close()
			continue
		}

		fmt.Printf("  ✓ SUCCESS!\n")
		db.Close()
		return
	}
}
