package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Price        int     `json:"price"`
	Stock        int     `json:"stock"`
	CategoryID   *int    `json:"category_id"`
	CategoryName *string `json:"category_name"`
}

func main() {
	// viper setup
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	dbConn := viper.GetString("DB_CONN")
	if dbConn == "" {
		log.Fatal("DB_CONN tidak ditemukan dalam environment variables")
	}

	db, err := sql.Open("pgx", dbConn)
	if err != nil {
		log.Fatal("Gagal connect ke database:", err)
	}
	defer db.Close()

	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("Gagal query products:", err)
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var p Product
		var catID sql.NullInt64
		var catName sql.NullString

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID, &catName)
		if err != nil {
			log.Fatal("Gagal scan row:", err)
		}

		if catID.Valid {
			v := int(catID.Int64)
			p.CategoryID = &v
		}
		if catName.Valid {
			p.CategoryName = &catName.String
		}

		products = append(products, p)
	}

	if len(products) == 0 {
		fmt.Println("⚠️  Database tidak memiliki produk. Silakan tambahkan produk terlebih dahulu.")
		return
	}

	fmt.Println("\n📦 Produk yang ada di database:\n")
	for _, p := range products {
		catName := "No Category"
		if p.CategoryName != nil {
			catName = *p.CategoryName
		}
		fmt.Printf("  [ID: %d] %s - Rp %d (Stock: %d) - Kategori: %s\n",
			p.ID, p.Name, p.Price, p.Stock, catName)
	}

	// Output as JSON for easy copy
	fmt.Println("\n📋 JSON format (untuk testing):\n")
	jsonData, _ := json.MarshalIndent(products, "", "  ")
	fmt.Println(string(jsonData))
}
