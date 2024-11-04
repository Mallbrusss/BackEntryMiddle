package repository

import (
	"github.com/Mallbrusss/BackEntryMiddle/models"

	"gorm.io/gorm"
)

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

func (ur *UserRepository) FindUser(username string) (*models.User, error) {
	var user models.User

	if err := ur.db.Where("login = ?", username).First(&user).Error; err != nil {
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
