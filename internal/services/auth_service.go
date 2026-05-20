package services

import (
	"errors"
	"log"
	"strings"
	"time"

	"expense_management_backend/internal/database"
	"expense_management_backend/internal/models"
	"expense_management_backend/internal/repositories"
	"expense_management_backend/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db       *gorm.DB
	userRepo *repositories.UserRepository
}

func NewAuthService(db *gorm.DB, userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		db:       db,
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(email, password, displayName string) (*models.User, string, string, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	txUserRepo := s.userRepo.WithTx(tx)
	existingUser, err := txUserRepo.FindByEmail(email)

	if err == nil {
		// Email already exists
		if existingUser.PasswordHash != nil {
			tx.Rollback()
			return nil, "", "", errors.New("Email already registered")
		}

		// User existed via Google but has no password yet. Link password to this account.
		hash, err := utils.HashPassword(password)
		if err != nil {
			tx.Rollback()
			return nil, "", "", errors.New("Failed to process password")
		}

		existingUser.PasswordHash = &hash
		if displayName != "" {
			existingUser.DisplayName = displayName
		}

		if err := txUserRepo.Save(existingUser); err != nil {
			tx.Rollback()
			return nil, "", "", errors.New("Failed to update user details")
		}

		if err := tx.Commit().Error; err != nil {
			return nil, "", "", errors.New("Failed to commit transaction")
		}

		token, err := utils.GenerateToken(existingUser.ID, existingUser.Email)
		if err != nil {
			return nil, "", "", errors.New("Failed to generate token")
		}

		refreshToken, err := s.generateAndSaveRefreshToken(tx, existingUser.ID)
		if err != nil {
			return nil, "", "", errors.New("Failed to generate refresh token")
		}

		return existingUser, token, refreshToken, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, "", "", errors.New("Database error")
	}

	// Create new user
	hash, err := utils.HashPassword(password)
	if err != nil {
		tx.Rollback()
		return nil, "", "", errors.New("Failed to hash password")
	}

	newUser := models.User{
		Email:        email,
		PasswordHash: &hash,
		DisplayName:  displayName,
		CurrencyCode: "VND",
	}

	if err := txUserRepo.Create(&newUser); err != nil {
		tx.Rollback()
		return nil, "", "", errors.New("Failed to create user")
	}

	// Seed default wallet and categories
	if err := database.SeedDefaultData(tx, newUser.ID); err != nil {
		tx.Rollback()
		return nil, "", "", errors.New("Failed to seed default data")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", "", errors.New("Failed to commit transaction")
	}

	token, err := utils.GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return nil, "", "", errors.New("Failed to generate token")
	}

	refreshToken, err := s.generateAndSaveRefreshToken(s.db, newUser.ID)
	if err != nil {
		return nil, "", "", errors.New("Failed to generate refresh token")
	}

	return &newUser, token, refreshToken, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("Invalid email or password")
		}
		return nil, "", "", errors.New("Database error")
	}

	if user.PasswordHash == nil {
		return nil, "", "", errors.New("This account is only registered with Google Login. Please use Google to sign in.")
	}

	if !utils.CheckPasswordHash(password, *user.PasswordHash) {
		return nil, "", "", errors.New("Invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", "", errors.New("Failed to generate token")
	}

	refreshToken, err := s.generateAndSaveRefreshToken(s.db, user.ID)
	if err != nil {
		return nil, "", "", errors.New("Failed to generate refresh token")
	}

	return user, token, refreshToken, nil
}

func (s *AuthService) GoogleLogin(idToken string) (*models.User, string, string, bool, error) {
	// Verify ID token with Google API
	googleInfo, err := utils.VerifyGoogleToken(idToken)
	if err != nil {
		return nil, "", "", false, errors.New("Invalid Google token: " + err.Error())
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	txUserRepo := s.userRepo.WithTx(tx)
	user, err := txUserRepo.FindByEmail(googleInfo.Email)

	if err == nil {
		// User with this email already exists
		// Link Google ID if not linked yet
		if user.GoogleID == nil || *user.GoogleID == "" {
			user.GoogleID = &googleInfo.Sub
			if user.PhotoURL == "" {
				user.PhotoURL = googleInfo.Picture
			}
			if err := txUserRepo.Save(user); err != nil {
				tx.Rollback()
				return nil, "", "", false, errors.New("Failed to link Google account")
			}
		}

		if err := tx.Commit().Error; err != nil {
			return nil, "", "", false, errors.New("Failed to commit transaction")
		}

		token, err := utils.GenerateToken(user.ID, user.Email)
		if err != nil {
			return nil, "", "", false, errors.New("Failed to generate token")
		}

		refreshToken, err := s.generateAndSaveRefreshToken(s.db, user.ID)
		if err != nil {
			return nil, "", "", false, errors.New("Failed to generate refresh token")
		}

		return user, token, refreshToken, false, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, "", "", false, errors.New("Database error")
	}

	// User does not exist, register new Google user
	newUser := models.User{
		Email:        googleInfo.Email,
		GoogleID:     &googleInfo.Sub,
		DisplayName:  googleInfo.Name,
		PhotoURL:     googleInfo.Picture,
		CurrencyCode: "VND",
	}

	if err := txUserRepo.Create(&newUser); err != nil {
		tx.Rollback()
		return nil, "", "", false, errors.New("Failed to create user")
	}

	// Seed default wallet and categories
	if err := database.SeedDefaultData(tx, newUser.ID); err != nil {
		tx.Rollback()
		return nil, "", "", false, errors.New("Failed to seed default data")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", "", false, errors.New("Failed to commit transaction")
	}

	token, err := utils.GenerateToken(newUser.ID, newUser.Email)
	if err != nil {
		return nil, "", "", false, errors.New("Failed to generate token")
	}

	refreshToken, err := s.generateAndSaveRefreshToken(s.db, newUser.ID)
	if err != nil {
		return nil, "", "", false, errors.New("Failed to generate refresh token")
	}

	return &newUser, token, refreshToken, true, nil
}

func (s *AuthService) RefreshToken(refreshTokenStr string) (string, string, error) {
	storedToken, err := s.userRepo.FindRefreshToken(refreshTokenStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("Invalid refresh token")
		}
		return "", "", errors.New("Database error")
	}

	// Check if expired
	if time.Now().After(storedToken.ExpiresAt) {
		s.userRepo.DeleteRefreshToken(storedToken)
		return "", "", errors.New("Refresh token expired. Please login again.")
	}

	// Rotate refresh token: Delete the old one
	s.userRepo.DeleteRefreshToken(storedToken)

	// Generate new access token
	accessToken, err := utils.GenerateToken(storedToken.UserID, storedToken.User.Email)
	if err != nil {
		return "", "", errors.New("Failed to generate access token")
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateAndSaveRefreshToken(s.db, storedToken.UserID)
	if err != nil {
		return "", "", errors.New("Failed to generate refresh token")
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) ForgotPassword(email string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("No account found with this email")
		}
		return "", errors.New("Database error")
	}

	if user.PasswordHash == nil {
		return "", errors.New("This account is registered with Google. Password reset is not available.")
	}

	// Generate numeric/alphanumeric 6-digit reset code
	resetToken, err := utils.GenerateRandomString(3) // 6 hex digits
	if err != nil {
		return "", errors.New("Failed to generate reset code")
	}
	resetToken = strings.ToUpper(resetToken)

	// Clean up any old reset tokens for this email
	s.userRepo.DeletePasswordResetsByEmail(email)

	// Save reset token (expires in 15 minutes)
	passwordReset := models.PasswordReset{
		Email:     email,
		Token:     resetToken,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.userRepo.CreatePasswordReset(&passwordReset); err != nil {
		return "", errors.New("Failed to store reset token")
	}

	// Print reset code to terminal/logs
	log.Printf("\n🔑 PASSWORD RESET CODE for %s: %s\n", email, resetToken)

	return resetToken, nil
}

func (s *AuthService) ResetPassword(email, token, newPassword string) error {
	passwordReset, err := s.userRepo.FindPasswordReset(email, strings.ToUpper(token))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Invalid email or reset code")
		}
		return errors.New("Database error")
	}

	// Check expiration
	if time.Now().After(passwordReset.ExpiresAt) {
		s.userRepo.DeletePasswordReset(passwordReset)
		return errors.New("Reset code expired")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("Failed to find user")
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("Failed to process new password")
	}

	user.PasswordHash = &hash
	if err := s.userRepo.Save(user); err != nil {
		return errors.New("Failed to update password")
	}

	// Delete used token
	s.userRepo.DeletePasswordReset(passwordReset)

	return nil
}

func (s *AuthService) generateAndSaveRefreshToken(db *gorm.DB, userID string) (string, error) {
	tokenStr, err := utils.GenerateRandomString(32) // 64 hex characters
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	refreshToken := models.RefreshToken{
		UserID:    userID,
		Token:     tokenStr,
		ExpiresAt: expiresAt,
	}

	txUserRepo := s.userRepo.WithTx(db)
	if err := txUserRepo.CreateRefreshToken(&refreshToken); err != nil {
		return "", err
	}

	return tokenStr, nil
}
