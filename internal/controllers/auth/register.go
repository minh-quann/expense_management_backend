package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register handles user registration and linking if Google login exists
// @Summary      Register a new user
// @Description  Create a new user account with email, password, and display name, or link to an existing Google account.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration Info"
// @Success      201  {object}  map[string]interface{} "Successful registration"
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, refreshToken, err := ctrl.authService.Register(req.Email, req.Password, req.DisplayName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "User registered successfully",
		"user":          user,
		"token":         token,
		"refresh_token": refreshToken,
	})
}
