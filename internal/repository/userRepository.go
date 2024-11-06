package repository

import (
	"github.com/Mallbrusss/BackEntryMiddle/models"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	FindUser(login string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteToken(token string) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *UserRepository) FindUser(login string) (*models.User, error) {
	var user models.User

	if err := ur.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) UpdateUser(user *models.User) error {
	return ur.db.Save(user).Error
}

func (ur *UserRepository) DeleteUser(user *models.User) error {
	return ur.db.Delete(user).Error
}

func (ur *UserRepository) DeleteToken(token string) error {
	return ur.db.Model(&models.User{}).Where("token = ?", token).Update("token", nil).Error
}
