package repositories

import (
	"expense_management_backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Save(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) CreateRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *UserRepository) FindRefreshToken(tokenStr string) (*models.RefreshToken, error) {
	var storedToken models.RefreshToken
	if err := r.db.Preload("User").Where("token = ?", tokenStr).First(&storedToken).Error; err != nil {
		return nil, err
	}
	return &storedToken, nil
}

func (r *UserRepository) DeleteRefreshToken(token *models.RefreshToken) error {
	return r.db.Delete(token).Error
}

func (r *UserRepository) CreatePasswordReset(reset *models.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *UserRepository) FindPasswordReset(email, token string) (*models.PasswordReset, error) {
	var pr models.PasswordReset
	if err := r.db.Where("email = ? AND token = ?", email, token).First(&pr).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *UserRepository) DeletePasswordReset(reset *models.PasswordReset) error {
	return r.db.Delete(reset).Error
}

func (r *UserRepository) DeletePasswordResetsByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.PasswordReset{}).Error
}
