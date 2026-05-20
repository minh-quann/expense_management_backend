package wallet

import (
	"net/http"

	"expense_management_backend/internal/services"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	walletService *services.WalletService
}

func NewWalletController(walletService *services.WalletService) *WalletController {
	return &WalletController{
		walletService: walletService,
	}
}

// GetWallets retrieves all active wallets for the logged-in user
// @Summary      Get all wallets
// @Description  Get a list of all active wallets belonging to the logged-in user.
// @Tags         wallets
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Wallet
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /wallets [get]
func (ctrl *WalletController) GetWallets(c *gin.Context) {
	userID := c.GetString("user_id")

	wallets, err := ctrl.walletService.GetWallets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallets)
}

// CreateWallet creates a new wallet for the logged-in user
// @Summary      Create a wallet
// @Description  Create a new wallet for the logged-in user.
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body WalletRequest true "Wallet Info"
// @Success      201  {object}  models.Wallet
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /wallets [post]
func (ctrl *WalletController) CreateWallet(c *gin.Context) {
	userID := c.GetString("user_id")

	var req WalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := ctrl.walletService.CreateWallet(
		userID,
		req.Name,
		req.Type,
		req.Balance,
		req.Currency,
		req.Icon,
		req.Color,
		req.ExcludeFromTotal,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// UpdateWallet updates an existing wallet's details
// @Summary      Update a wallet
// @Description  Update the details of an existing wallet.
// @Tags         wallets
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Wallet ID"
// @Param        request body WalletRequest true "Updated Wallet Info"
// @Success      200  {object}  models.Wallet
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      404  {object}  map[string]interface{} "Wallet not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /wallets/{id} [put]
func (ctrl *WalletController) UpdateWallet(c *gin.Context) {
	userID := c.GetString("user_id")
	walletID := c.Param("id")

	var req WalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := ctrl.walletService.UpdateWallet(
		userID,
		walletID,
		req.Name,
		req.Type,
		req.Balance,
		req.Currency,
		req.Icon,
		req.Color,
		req.ExcludeFromTotal,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// DeleteWallet deletes a wallet if it does not have any associated transactions
// @Summary      Delete a wallet
// @Description  Soft delete a wallet if there are no transactions associated with it.
// @Tags         wallets
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Wallet ID"
// @Success      200  {object}  map[string]interface{} "Wallet deleted successfully"
// @Failure      404  {object}  map[string]interface{} "Wallet not found"
// @Failure      409  {object}  map[string]interface{} "Conflict - Wallet has associated transactions"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /wallets/{id} [delete]
func (ctrl *WalletController) DeleteWallet(c *gin.Context) {
	userID := c.GetString("user_id")
	walletID := c.Param("id")

	err := ctrl.walletService.DeleteWallet(userID, walletID)
	if err != nil {
		// Can return appropriate errors. If containing 'transactions', return 409 Conflict.
		// For simplicity, we match the error string or return 400.
		statusCode := http.StatusBadRequest
		if err.Error() == "Wallet not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "Cannot delete wallet because it has associated transactions. Please delete or reassign those transactions first." {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
}
