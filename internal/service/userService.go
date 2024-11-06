package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"github.com/Mallbrusss/BackEntryMiddle/internal/repository"
	"github.com/Mallbrusss/BackEntryMiddle/models"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Register(login, password string, isAdmin bool) (*models.User, error)
	Authenticate(login, password string) (*models.User, error)
	DeleteToken(token string) error
}

type UserService struct {
	userRepository  repository.UserRepositoryInterface
	isValidLogin    func(string) bool
	isValidPassword func(string) bool
}

func NewUserService(userRepository repository.UserRepositoryInterface, isValidLogin func(string) bool, isValidPassword func(string) bool) *UserService {
	return &UserService{
		userRepository:  userRepository,
		isValidLogin:    isValidLogin,
		isValidPassword: isValidPassword,
	}
}

func (us *UserService) Register(login, password string, isAdmin bool) (*models.User, error) {
	// TODO: Перенести в middleware
	if ok := us.isValidLogin(login); !ok {
		log.Println("Invalid login")
		return nil, errors.New("invalid login")
	}
	if ok := us.isValidPassword(password); !ok {
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
	//TODO: Кешировать токены
	err = us.userRepository.UpdateUser(user)
	if err != nil {
		return nil, errors.New("error update user (token)")
	}

	return user, nil
}

func (us *UserService) DeleteToken(token string) error {
	//TODO: Убирать токен из кеша
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
