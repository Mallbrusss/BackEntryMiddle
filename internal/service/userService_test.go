package service

import (
	"errors"
	"testing"

	"github.com/Mallbrusss/BackEntryMiddle/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUser(login string) (*models.User, error) {
	args := m.Called(login)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func TestUserService_Register(t *testing.T) {
	t.Run("should return error if login is invalid", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		isValidLogin := func(login string) bool {
			return false
		}
		isValidPassword := func(password string) bool {
			return true
		}

		us := NewUserService(mockRepo, isValidLogin, isValidPassword)

		_, err := us.Register("invalid_login", "password", false)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid login", err.Error())
	})

	t.Run("should return error if password is invalid", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		isValidLogin := func(login string) bool {
			return true
		}
		isValidPassword := func(password string) bool {
			return false
		}

		us := NewUserService(mockRepo, isValidLogin, isValidPassword)

		_, err := us.Register("valid_login", "short", false)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid password", err.Error())
	})
	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		isValidLogin := func(login string) bool {
			return true
		}
		isValidPassword := func(password string) bool {
			return true
		}

		mockRepo.On("CreateUser", mock.Anything).Return(nil).Once()

		us := NewUserService(mockRepo, isValidLogin, isValidPassword)

		user, err := us.Register("valid_login", "strongpassword", false)
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "valid_login", user.Login)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Authenticate(t *testing.T) {
	t.Run("should return error if user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		isValidLogin := func(login string) bool {
			return true
		}
		isValidPassword := func(password string) bool {
			return true
		}

		us := NewUserService(mockRepo, isValidLogin, isValidPassword)

		mockRepo.On("FindUser", "nonexistent").Return(nil, errors.New("user not found")).Once()

		_, err := us.Authenticate("nonexistent", "password")
		assert.NotNil(t, err)
		assert.Equal(t, "user not found", err.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if password is incorrect", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		isValidLogin := func(login string) bool {
			return true
		}
		isValidPassword := func(password string) bool {
			return true
		}

		us := NewUserService(mockRepo, isValidLogin, isValidPassword)

		user := &models.User{Login: "valid_login", Password: "$2a$10$somehashedpassword"}
		mockRepo.On("FindUser", "valid_login").Return(user, nil).Once()

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("wrongpassword"))
		assert.NotNil(t, err)

		_, errAuth := us.Authenticate("valid_login", "wrongpassword")
		assert.NotNil(t, errAuth)
		assert.Equal(t, "invalid password", errAuth.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return user with token after successful authentication", func(t *testing.T) {
		mockRepo := new(MockUserRepository)

		password := "correctpassword"

		hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			t.Fatal(err)
		}

		user := &models.User{
			Login:    "valid_login",
			Password: string(hashPwd),
		}

		mockRepo.On("FindUser", "valid_login").Return(user, nil)

		mockRepo.On("UpdateUser", mock.Anything).Return(nil)

		isValidLogin := func(login string) bool {
			return true
		}
		isValidPassword := func(password string) bool {
			return true
		}

		us := NewUserService(mockRepo, isValidLogin, isValidPassword)

		authUser, err := us.Authenticate("valid_login", password)
		assert.Nil(t, err)
		assert.NotNil(t, authUser)
		assert.Equal(t, "valid_login", authUser.Login)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteToken(t *testing.T) {
	t.Run("should delete token successfully", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		isValidLogin := func(login string) bool {
			return true
		}
		isValidPassword := func(password string) bool {
			return true
		}
		us := NewUserService(mockRepo, isValidLogin, isValidPassword)

		mockRepo.On("DeleteToken", "some_token").Return(nil).Once()

		err := us.DeleteToken("some_token")
		assert.Nil(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if unable to delete token", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		isValidLogin := func(login string) bool {
			return true
		}
		isValidPassword := func(password string) bool {
			return true
		}
		us := NewUserService(mockRepo, isValidLogin, isValidPassword)

		mockRepo.On("DeleteToken", "some_token").Return(errors.New("failed to delete token")).Once()

		err := us.DeleteToken("some_token")
		assert.NotNil(t, err)
		assert.Equal(t, "failed to delete token", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
