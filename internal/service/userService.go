package service

import (
	"internal/models"
	"internal/repository"
	"github.com/Mallbrusss/BackEntryMiddle/pkg/utils"
)

type UserService struct {
	userRepository *repository.UserRepository
	authToken      string
}

func NewUserService(userRepository *repository.UserRepository, authToken string) *UserService {
	return &UserService{
		userRepository: userRepository,
		authToken:      authToken,
	}
}

func (us *UserService) Register(login, password string) (*models.User, error) {

}

func (us *UserService) Authenticate() {

}

func (us *UserService) GetUserByUsername() {

}

func (us *UserService) ValidateToken() {

}
