package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logout handles logging out the user (revoking the session/token)
// @Summary      Logout user
// @Description  Revoke/invalidate the refresh token, logging out the user from the current session.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest false "Refresh Token to invalidate"
// @Success      200  {object}  map[string]interface{} "Successfully logged out"
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /auth/logout [post]
func (ctrl *AuthController) Logout(c *gin.Context) {
	// For now, since JWT is stateless, we return a success response immediately.
	// This endpoint is fully prepared for future blacklisting or session revoking logic.
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}
