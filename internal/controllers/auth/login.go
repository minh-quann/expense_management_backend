package auth

import (
	"net/http"

	"expense_management_backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// Login handles standard email/password authentication
// @Summary      Login user
// @Description  Authenticate user with email and password and return a JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Info"
// @Success      200  {object}  map[string]interface{} "Successful login"
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/login [post]
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "AUTH_INVALID_REQUEST", err.Error())
		return
	}

	user, token, refreshToken, err := ctrl.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.RespondWithCustomError(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"user":          user,
		"token":         token,
		"refresh_token": refreshToken,
	})
}

// GoogleLogin handles verification of Google ID token and accounts linking
// @Summary      Google Login
// @Description  Verify Google ID Token and register or login user, linking to email/password if email matches.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body GoogleLoginRequest true "Google ID Token Info"
// @Success      200  {object}  map[string]interface{} "Successful login"
// @Success      211  {object}  map[string]interface{} "Successful registration via Google"
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/google [post]
func (ctrl *AuthController) GoogleLogin(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "AUTH_INVALID_REQUEST", err.Error())
		return
	}

	user, token, refreshToken, isNewUser, err := ctrl.authService.GoogleLogin(req.IDToken)
	if err != nil {
		utils.RespondWithCustomError(c, http.StatusUnauthorized, err)
		return
	}

	statusCode := http.StatusOK
	message := "Google login successful (linked account)"
	if isNewUser {
		statusCode = http.StatusCreated
		message = "Google registration successful"
	}

	c.JSON(statusCode, gin.H{
		"message":       message,
		"user":          user,
		"token":         token,
		"refresh_token": refreshToken,
	})
}
