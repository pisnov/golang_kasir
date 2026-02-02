package main

import (
	"fmt"
	"net/url"
)

func main() {
	password := "Biawak123!@#"
	encoded := url.QueryEscape(password)
	fmt.Println("Original:", password)
	fmt.Println("Encoded:", encoded)

	// Try different encoding
	connStr := fmt.Sprintf("postgresql://postgres.btuyxzhfurruethkzqfu:%s@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=disable", encoded)
	fmt.Println("Connection String:", connStr)
}
