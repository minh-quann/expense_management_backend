package transaction

import (
	"net/http"
	"strconv"

	"expense_management_backend/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	txService *services.TransactionService
}

func NewTransactionController(txService *services.TransactionService) *TransactionController {
	return &TransactionController{
		txService: txService,
	}
}

// GetTransactions retrieves the user's transactions (optionally filtered by wallet)
// @Summary      Get all transactions
// @Description  Get a list of all transactions belonging to the logged-in user, optionally filtered by wallet ID.
// @Tags         transactions
// @Produce      json
// @Security     BearerAuth
// @Param        wallet_id query string false "Filter by wallet ID"
// @Success      200  {array}   models.Transaction
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /transactions [get]
func (ctrl *TransactionController) GetTransactions(c *gin.Context) {
	userID := c.GetString("user_id")
	walletID := c.Query("wallet_id")

	transactions, err := ctrl.txService.GetTransactions(userID, walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// CreateTransaction creates a new transaction and adjusts wallet balances inside a DB transaction
// @Summary      Create a transaction
// @Description  Create a transaction (EXPENSE, INCOME, or TRANSFER) and automatically update wallet balance(s).
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body TransactionRequest true "Transaction Info"
// @Success      201  {object}  models.Transaction
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /transactions [post]
func (ctrl *TransactionController) CreateTransaction(c *gin.Context) {
	userID := c.GetString("user_id")

	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := ctrl.txService.CreateTransaction(
		userID,
		req.Amount,
		req.Type,
		req.CategoryID,
		req.WalletID,
		req.ToWalletID,
		req.Date,
		req.Note,
		req.ImageURL,
		req.RecurringID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// UpdateTransaction updates an existing transaction details and re-calculates wallet balance adjustments inside a DB transaction
// @Summary      Update a transaction
// @Description  Update transaction details. Reverts the old balance adjustments and applies new adjustments in a single transaction.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Transaction ID"
// @Param        request body TransactionRequest true "Updated Transaction Info"
// @Success      200  {object}  models.Transaction
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      404  {object}  map[string]interface{} "Transaction not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /transactions/{id} [put]
func (ctrl *TransactionController) UpdateTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	transactionID := c.Param("id")

	var req TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := ctrl.txService.UpdateTransaction(
		userID,
		transactionID,
		req.Amount,
		req.Type,
		req.CategoryID,
		req.WalletID,
		req.ToWalletID,
		req.Date,
		req.Note,
		req.ImageURL,
		req.RecurringID,
	)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "Transaction not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction deletes an existing transaction and reverts its wallet balance adjustment(s) inside a DB transaction
// @Summary      Delete a transaction
// @Description  Delete a transaction and automatically revert the wallet balance adjustment.
// @Tags         transactions
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Transaction ID"
// @Success      200  {object}  map[string]interface{} "Transaction deleted successfully"
// @Failure      404  {object}  map[string]interface{} "Transaction not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /transactions/{id} [delete]
func (ctrl *TransactionController) DeleteTransaction(c *gin.Context) {
	userID := c.GetString("user_id")
	transactionID := c.Param("id")

	err := ctrl.txService.DeleteTransaction(userID, transactionID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "Transaction not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

// GetStatistics retrieves statistics of income and expenses grouped by categories
// @Summary      Get transaction statistics
// @Description  Get sum of income and expenses grouped by categories for a given month and year
// @Tags         transactions
// @Produce      json
// @Security     BearerAuth
// @Param        month query int false "Month (1-12)"
// @Param        year query int false "Year (e.g. 2026)"
// @Param        wallet_id query string false "Filter by wallet ID"
// @Success      200  {object}  services.StatisticsResponse
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /transactions/statistics [get]
func (ctrl *TransactionController) GetStatistics(c *gin.Context) {
	userID := c.GetString("user_id")
	walletID := c.Query("wallet_id")

	var month, year int
	var err error

	monthStr := c.Query("month")
	if monthStr != "" {
		month, err = strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month parameter. Must be 1-12"})
			return
		}
	}

	yearStr := c.Query("year")
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil || year < 1900 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
			return
		}
	}

	stats, err := ctrl.txService.GetStatistics(userID, walletID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
