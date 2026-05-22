package router

import (
	"net/http"

	"expense_management_backend/internal/controllers/auth"
	"expense_management_backend/internal/controllers/category"
	"expense_management_backend/internal/controllers/transaction"
	"expense_management_backend/internal/controllers/user"
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
	userCtrl *user.UserController,
) *gin.Engine {
	r := gin.Default()

	// Enable CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static files in uploads folder
	r.Static("/uploads", "./uploads")

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
			authGroup.POST("/logout", authCtrl.Logout)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Profile routes
			profile := protected.Group("/profile")
			{
				profile.GET("", userCtrl.GetProfile)
				profile.PUT("", userCtrl.UpdateProfile)
				profile.POST("/avatar", userCtrl.UploadAvatar)

				// PIN and Security routes
				profile.POST("/pin", userCtrl.SetPIN)
				profile.POST("/pin/verify", userCtrl.VerifyPIN)
				profile.GET("/pin/security-question", userCtrl.GetSecurityQuestion)
				profile.POST("/pin/reset", userCtrl.ResetPIN)
				profile.DELETE("/pin", userCtrl.DisablePIN)
			}

			// Wallet routes
			wallets := protected.Group("/wallets")
			{
				wallets.GET("", walletCtrl.GetWallets)
				wallets.POST("", walletCtrl.CreateWallet)
				wallets.PUT("/:id", walletCtrl.UpdateWallet)
				wallets.DELETE("/:id", walletCtrl.DeleteWallet)
				wallets.PATCH("/:id/favorite", walletCtrl.ToggleFavorite)
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
				transactions.GET("/statistics", transactionCtrl.GetStatistics)
				transactions.POST("", transactionCtrl.CreateTransaction)
				transactions.PUT("/:id", transactionCtrl.UpdateTransaction)
				transactions.DELETE("/:id", transactionCtrl.DeleteTransaction)
			}
		}
	}

	return r
}
