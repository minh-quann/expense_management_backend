package services

import (
	"errors"

	"expense_management_backend/internal/models"
	"expense_management_backend/internal/repositories"

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
