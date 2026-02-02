package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	ports := []string{"6543", "5432"}
	modes := []string{"sslmode=disable", "sslmode=require"}

	for _, port := range ports {
		for _, mode := range modes {
			connStr := fmt.Sprintf("postgresql://postgres.btuyxzhfurruethkzqfu:Biawak123%%21%%40%%23@aws-1-ap-southeast-1.pooler.supabase.com:%s/postgres?%s", port, mode)
			log.Printf("Testing port %s with %s\n", port, mode)

			db, err := sql.Open("postgres", connStr)
			if err != nil {
				log.Printf("  Open failed: %v\n", err)
				continue
			}

			if err := db.Ping(); err != nil {
				log.Printf("  Ping failed: %v\n", err)
				db.Close()
				continue
			}

			fmt.Printf("  ✓ SUCCESS with port %s!\n", port)
			db.Close()
			return
		}
	}
}
