package service

import (
	"errors"
	"internal/models"
	"internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
	authToken      string
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) Register(login, password, token string) (*models.User, error) {
	//TODO: МБ надо  будет убрать
	// if ok := utils.IsValidLogin(login); !ok{
	// 	log.Println("Invalid login")
	// 	return nil, errors.New("invalid login")
	// }
	// if ok := utils.IsValidPassword(password); !ok{
	// 	log.Println("Invalid password")
	// 	return nil, errors.New("invalid password")
	// }

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Login:    login,
		Password: string(hashPwd),
		Token:    token,
	}

	if err := us.userRepository.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) Authenticate(login, password string) (*models.User, error) {
	// if token != us.authToken {
	// 	return nil, errors.New("invalid token")
	// }

	user, err := us.userRepository.FindUser(login)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}

func (us *UserService) DeleteToken() {

}
