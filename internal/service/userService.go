package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"github.com/Mallbrusss/BackEntryMiddle/internal/repository"
	"github.com/Mallbrusss/BackEntryMiddle/models"
	"github.com/Mallbrusss/BackEntryMiddle/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
	// authToken      string
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) Register(login, password string, isAdmin bool) (*models.User, error) {
	// TODO: Перенести в middleware

	if ok := utils.IsValidLogin(login); !ok {
		log.Println("Invalid login")
		return nil, errors.New("invalid login")
	}
	if ok := utils.IsValidPassword(password); !ok {
		log.Println("Invalid password")
		return nil, errors.New("invalid password")
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Login:    login,
		Password: string(hashPwd),
		IsAdmin:  isAdmin,
	}

	if err := us.userRepository.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Authenticate(login, password string) (*models.User, error) {
	user, err := us.userRepository.FindUser(login)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	newToken, err := us.generateToken(8)
	if err != nil {
		return nil, errors.New("error generate token")
	}

	user.Token = newToken
	err = us.userRepository.UpdateUser(user)
	if err != nil {
		return nil, errors.New("error update user (token)")
	}

	return user, nil
}

func (us *UserService) DeleteToken(token string) error {
	return us.userRepository.DeleteToken(token)
}

func (us *UserService) generateToken(length int) (string, error) {
	token := make([]byte, length)

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}
