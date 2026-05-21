package user

import (
	"net/http"
	"os"
	"path/filepath"

	"expense_management_backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetProfile retrieves the currently logged in user's profile
// @Summary      Get user profile
// @Description  Get detailed profile information of the current logged-in user.
// @Tags         profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]interface{} "User not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /profile [get]
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := ctrl.userService.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the current logged-in user's profile information
// @Summary      Update user profile
// @Description  Update profile details like display name and preferred currency code.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateProfileRequest true "Profile Info"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /profile [put]
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userService.UpdateProfile(userID, req.DisplayName, req.CurrencyCode, req.PhoneNumber, req.Address, req.Gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UploadAvatar handles uploading and setting a new profile photo for the logged-in user
// @Summary      Upload profile avatar
// @Description  Upload an image file to set as the user's profile avatar.
// @Tags         profile
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        avatar formData file true "Avatar Image File"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /profile/avatar [post]
func (ctrl *UserController) UploadAvatar(c *gin.Context) {
	userID := c.GetString("user_id")

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read avatar file from form"})
		return
	}

	// Ensure upload directory exists
	uploadDir := filepath.Join("uploads", "avatars")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate a unique filename using UUID
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// Save the file locally
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	// Save relative URL path to database
	photoURL := "/uploads/avatars/" + filename
	user, err := ctrl.userService.UpdateAvatar(userID, photoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// SetPIN handles configuring PIN code protection
// @Summary      Set PIN code
// @Description  Configure or update the user's PIN code and security question/answer.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body SetPINRequest true "PIN Settings"
// @Success      200  {object}  map[string]interface{} "PIN set successfully"
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /profile/pin [post]
func (ctrl *UserController) SetPIN(c *gin.Context) {
	userID := c.GetString("user_id")

	var req SetPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.userService.SetPIN(userID, req.Pin, req.SecurityQuestion, req.SecurityAnswer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN and security question set successfully"})
}

// VerifyPIN handles verifying the PIN code
// @Summary      Verify PIN code
// @Description  Verify if the provided PIN code matches the stored one.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body VerifyPINRequest true "PIN Code"
// @Success      200  {object}  map[string]interface{} "PIN verified"
// @Failure      400  {object}  map[string]interface{} "Bad request or incorrect PIN"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Router       /profile/pin/verify [post]
func (ctrl *UserController) VerifyPIN(c *gin.Context) {
	userID := c.GetString("user_id")

	var req VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verified, err := ctrl.userService.VerifyPIN(userID, req.Pin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !verified {
		c.JSON(http.StatusBadRequest, gin.H{"verified": false, "error": "Incorrect PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"verified": true, "message": "PIN verified successfully"})
}

// GetSecurityQuestion retrieves the security question for recovery
// @Summary      Get security question
// @Description  Get the user's configured security question for PIN recovery.
// @Tags         profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{} "Security question"
// @Failure      400  {object}  map[string]interface{} "Not configured"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Router       /profile/pin/security-question [get]
func (ctrl *UserController) GetSecurityQuestion(c *gin.Context) {
	userID := c.GetString("user_id")

	question, err := ctrl.userService.GetSecurityQuestion(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"security_question": question})
}

// ResetPIN handles resetting the PIN using the security answer
// @Summary      Reset PIN code
// @Description  Reset the PIN code using the correct security answer.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body ResetPINRequest true "Reset PIN Info"
// @Success      200  {object}  map[string]interface{} "PIN reset successfully"
// @Failure      400  {object}  map[string]interface{} "Incorrect answer or invalid PIN"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Router       /profile/pin/reset [post]
func (ctrl *UserController) ResetPIN(c *gin.Context) {
	userID := c.GetString("user_id")

	var req ResetPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.userService.ResetPIN(userID, req.SecurityAnswer, req.NewPin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN reset successfully"})
}

// DisablePIN handles disabling PIN protection
// @Summary      Disable PIN code
// @Description  Disable PIN code protection by providing the correct PIN.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body DisablePINRequest true "Disable PIN Info"
// @Success      200  {object}  map[string]interface{} "PIN disabled successfully"
// @Failure      400  {object}  map[string]interface{} "Incorrect PIN"
// @Failure      401  {object}  map[string]interface{} "Unauthorized"
// @Router       /profile/pin [delete]
func (ctrl *UserController) DisablePIN(c *gin.Context) {
	userID := c.GetString("user_id")

	var req DisablePINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.userService.DisablePIN(userID, req.Pin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN protection disabled successfully"})
}
