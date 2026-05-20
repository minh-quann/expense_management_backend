package router

import (
	"net/http"

	"expense_management_backend/internal/controllers/auth"
	"expense_management_backend/internal/controllers/category"
	"expense_management_backend/internal/controllers/transaction"
	"expense_management_backend/internal/controllers/wallet"
	"expense_management_backend/internal/middleware"
	_ "expense_management_backend/internal/docs" // Import swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter initializes the Gin router with all routes
func SetupRouter(
	authCtrl *auth.AuthController,
	walletCtrl *wallet.WalletController,
	categoryCtrl *category.CategoryController,
	transactionCtrl *transaction.TransactionController,
) *gin.Engine {
	r := gin.Default()

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Expense Management API is running",
		})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Authentication routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authCtrl.Register)
			authGroup.POST("/login", authCtrl.Login)
			authGroup.POST("/google", authCtrl.GoogleLogin)
			authGroup.POST("/refresh", authCtrl.RefreshToken)
			authGroup.POST("/forgot-password", authCtrl.ForgotPassword)
			authGroup.POST("/reset-password", authCtrl.ResetPassword)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", func(c *gin.Context) {
				userID := c.GetString("user_id")
				email := c.GetString("email")
				c.JSON(http.StatusOK, gin.H{
					"user_id": userID,
					"email":   email,
				})
			})

			// Wallet routes
			wallets := protected.Group("/wallets")
			{
				wallets.GET("", walletCtrl.GetWallets)
				wallets.POST("", walletCtrl.CreateWallet)
				wallets.PUT("/:id", walletCtrl.UpdateWallet)
				wallets.DELETE("/:id", walletCtrl.DeleteWallet)
			}

			// Category routes
			categories := protected.Group("/categories")
			{
				categories.GET("", categoryCtrl.GetCategories)
				categories.POST("", categoryCtrl.CreateCategory)
				categories.PUT("/:id", categoryCtrl.UpdateCategory)
				categories.DELETE("/:id", categoryCtrl.DeleteCategory)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.GET("", transactionCtrl.GetTransactions)
				transactions.POST("", transactionCtrl.CreateTransaction)
				transactions.PUT("/:id", transactionCtrl.UpdateTransaction)
				transactions.DELETE("/:id", transactionCtrl.DeleteTransaction)
			}
		}
	}

	return r
}
