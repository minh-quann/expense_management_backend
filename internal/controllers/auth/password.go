package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ForgotPassword requests a password reset token
// @Summary      Forgot password request
// @Description  Submit email to receive a password reset token (token is printed to server console logs for development).
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body ForgotPasswordRequest true "Email address"
// @Success      200  {object}  map[string]interface{} "Reset token generated successfully"
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      404  {object}  map[string]interface{} "User not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/forgot-password [post]
func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.authService.ForgotPassword(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reset code generated successfully. Please check server logs.",
		"token":   token, // Return it in response as well to make it extremely easy to test!
	})
}

// ResetPassword resets a user's password using a reset token
// @Summary      Reset password
// @Description  Submit email, reset token, and new password to complete password reset.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body ResetPasswordRequest true "Reset Password Info"
// @Success      200  {object}  map[string]interface{} "Password reset successful"
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      401  {object}  map[string]interface{} "Invalid or expired reset token"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/reset-password [post]
func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.authService.ResetPassword(req.Email, req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully. You can now login with your new password.",
	})
}
