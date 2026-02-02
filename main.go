package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pisnov/golang_kasir/database"
	"github.com/pisnov/golang_kasir/handlers"
	"github.com/pisnov/golang_kasir/repositories"
	"github.com/pisnov/golang_kasir/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// viper setup: read environment and .env if present
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	// Initialize database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize repositories, services, and handlers
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// categories
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	http.HandleFunc("/api/kategori", categoryHandler.HandleCategories)
	http.HandleFunc("/api/kategori/", categoryHandler.HandleCategoryByID)

	// Start server
	addr := "0.0.0.0:" + config.Port
	log.Println("Server running di", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("gagal running server", err)
	}
}
