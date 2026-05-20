package main

import (
	"fmt"
	"log"

	"expense_management_backend/config"
	"expense_management_backend/internal/controllers/auth"
	"expense_management_backend/internal/controllers/category"
	"expense_management_backend/internal/controllers/transaction"
	"expense_management_backend/internal/controllers/wallet"
	"expense_management_backend/internal/database"
	"expense_management_backend/internal/models"
	"expense_management_backend/internal/repositories"
	"expense_management_backend/internal/router"
	"expense_management_backend/internal/services"

	"github.com/joho/godotenv"
)

// @title           Expense Management API
// @version         1.0
// @description     This is a Go API server for Expense Management.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by your JWT token.
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database.Connect(&cfg.Database)

	// Run migrations
	database.Migrate(
		&models.User{},
		&models.Wallet{},
		&models.Category{},
		&models.Transaction{},
		&models.Budget{},
		&models.RecurringTemplate{},
		&models.SavingGoal{},
		&models.RefreshToken{},
		&models.PasswordReset{},
	)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	walletRepo := repositories.NewWalletRepository(database.DB)
	categoryRepo := repositories.NewCategoryRepository(database.DB)
	txRepo := repositories.NewTransactionRepository(database.DB)

	// Initialize services
	authService := services.NewAuthService(database.DB, userRepo)
	walletService := services.NewWalletService(walletRepo, txRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	txService := services.NewTransactionService(database.DB, txRepo, walletRepo, categoryRepo)

	// Initialize controllers
	authCtrl := auth.NewAuthController(authService)
	walletCtrl := wallet.NewWalletController(walletService)
	categoryCtrl := category.NewCategoryController(categoryService)
	transactionCtrl := transaction.NewTransactionController(txService)

	// Setup router
	r := router.SetupRouter(authCtrl, walletCtrl, categoryCtrl, transactionCtrl)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("🚀 Server starting on http://localhost%s\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
