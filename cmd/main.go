package main

import (
	"fmt"
	"log"

	"expense_management_backend/config"
	"expense_management_backend/database"
	"expense_management_backend/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database.Connect(&cfg.Database)

	// Run migrations (add models here later)
	// database.Migrate(&models.User{}, &models.Expense{})

	// Setup router
	r := router.SetupRouter()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("🚀 Server starting on http://localhost%s\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
