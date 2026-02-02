package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	modes := []string{
		"sslmode=verify-full",
		"sslmode=verify-ca",
		"sslmode=disable",
	}

	certPaths := []string{
		"/etc/ssl/certs/ca-certificates.crt",
		"/etc/ssl/certs/ca-bundle.crt",
		"/etc/pki/tls/certs/ca-bundle.crt",
		"",
	}

	for _, mode := range modes {
		for _, certPath := range certPaths {
			var connStr string
			if certPath != "" && (mode == "sslmode=verify-full" || mode == "sslmode=verify-ca") {
				if _, err := os.Stat(certPath); err == nil {
					connStr = fmt.Sprintf("postgresql://postgres.btuyxzhfurruethkzqfu:Biawak123%%21%%40%%23@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?%s&sslrootcert=%s", mode, certPath)
				} else {
					continue
				}
			} else {
				connStr = fmt.Sprintf("postgresql://postgres.btuyxzhfurruethkzqfu:Biawak123%%21%%40%%23@aws-1-ap-southeast-1.pooler.supabase.com:6543/postgres?%s", mode)
			}

			desc := mode
			if certPath != "" {
				desc += fmt.Sprintf(" (cert: %s)", certPath)
			}
			log.Printf("Testing: %s\n", desc)

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

			fmt.Printf("  ✓ SUCCESS!\n")
			db.Close()
			return
		}
	}

	fmt.Println("No connection method worked")
}
