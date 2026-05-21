package services

import (
	"errors"
	"strings"

	"expense_management_backend/internal/models"
	"expense_management_backend/internal/repositories"
	"expense_management_backend/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile retrieves the profile details of the user by ID
func (s *UserService) GetProfile(userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return user, nil
}

// UpdateProfile updates the profile details of the user
func (s *UserService) UpdateProfile(userID, displayName, currencyCode, phoneNumber, address, gender string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	if displayName != "" {
		user.DisplayName = displayName
	}
	if currencyCode != "" {
		user.CurrencyCode = currencyCode
	}
	user.PhoneNumber = phoneNumber
	user.Address = address
	user.Gender = gender

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateAvatar updates the photo URL of the user
func (s *UserService) UpdateAvatar(userID, photoURL string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	user.PhotoURL = photoURL

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// SetPIN configures or updates the PIN code and security question/answer for the user
func (s *UserService) SetPIN(userID, pin, question, answer string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("User not found")
	}

	if len(pin) < 4 || len(pin) > 6 {
		return errors.New("PIN must be between 4 and 6 characters")
	}

	if question == "" || answer == "" {
		return errors.New("Security question and answer are required")
	}

	pinHash, err := utils.HashPassword(pin)
	if err != nil {
		return err
	}

	cleanAnswer := strings.ToLower(strings.TrimSpace(answer))
	answerHash, err := utils.HashPassword(cleanAnswer)
	if err != nil {
		return err
	}

	user.PinHash = &pinHash
	user.SecurityQuestion = &question
	user.SecurityAnswerHash = &answerHash

	return s.userRepo.Save(user)
}

// VerifyPIN verifies the provided PIN code against the stored hash
func (s *UserService) VerifyPIN(userID, pin string) (bool, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, errors.New("User not found")
	}

	if user.PinHash == nil {
		return false, errors.New("PIN is not configured for this account")
	}

	return utils.CheckPasswordHash(pin, *user.PinHash), nil
}

// GetSecurityQuestion retrieves the security question of the user
func (s *UserService) GetSecurityQuestion(userID string) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", errors.New("User not found")
	}

	if user.SecurityQuestion == nil || *user.SecurityQuestion == "" {
		return "", errors.New("Security question is not configured for this account")
	}

	return *user.SecurityQuestion, nil
}

// ResetPIN allows resetting the PIN using the security answer
func (s *UserService) ResetPIN(userID, answer, newPin string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("User not found")
	}

	if user.SecurityAnswerHash == nil {
		return errors.New("Security question/answer is not configured for this account")
	}

	cleanAnswer := strings.ToLower(strings.TrimSpace(answer))
	if !utils.CheckPasswordHash(cleanAnswer, *user.SecurityAnswerHash) {
		return errors.New("Incorrect security answer")
	}

	if len(newPin) < 4 || len(newPin) > 6 {
		return errors.New("PIN must be between 4 and 6 characters")
	}

	newPinHash, err := utils.HashPassword(newPin)
	if err != nil {
		return err
	}

	user.PinHash = &newPinHash
	return s.userRepo.Save(user)
}

// DisablePIN disables the PIN code protection by removing the stored PIN hash
func (s *UserService) DisablePIN(userID, pin string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("User not found")
	}

	if user.PinHash == nil {
		return errors.New("PIN is not configured for this account")
	}

	if !utils.CheckPasswordHash(pin, *user.PinHash) {
		return errors.New("Incorrect PIN")
	}

	user.PinHash = nil
	user.SecurityQuestion = nil
	user.SecurityAnswerHash = nil

	return s.userRepo.Save(user)
}
