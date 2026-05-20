package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefreshToken handles regenerating access and refresh tokens using a valid refresh token
// @Summary      Refresh access token
// @Description  Provide a valid refresh token to get a new access token and a rotated refresh token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh Token Info"
// @Success      200  {object}  map[string]interface{} "Successful token refresh"
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      401  {object}  map[string]interface{} "Invalid or expired refresh token"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/refresh [post]
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, refreshToken, err := ctrl.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}
