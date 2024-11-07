package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mallbrusss/BackEntryMiddle/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockErrorResponse struct {
	mock.Mock
}

func (m *MockErrorResponse) GetErrorResponse(code int) models.ErrorResponse {
	args := m.Called(code)
	return args.Get(0).(models.ErrorResponse)
}

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(login, password string, isAdmin bool) (*models.User, error) {
	args := m.Called(login, password, isAdmin)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) Authenticate(login, password string) (*models.User, error) {
	args := m.Called(login, password)

	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserService) DeleteToken(token string) error {
	args := m.Called(token)

	return args.Error(1)
}

func TestRegister(t *testing.T) {
	mockUserService := new(MockUserService)
	mockErrorResponse := new(MockErrorResponse)

	mockErrorResponse.On("GetErrorResponse", http.StatusBadRequest).
		Return(models.ErrorResponse{
			Code: 123, 
			Text: "So sad",
		})

	handler := UserHandler{
		UserService: mockUserService, 
		errRes: mockErrorResponse,
	}

	t.Run("Bad Request", func(t *testing.T) {
		req := models.User{
			Login:    "",
			Password: "passwordqwerty",
			Token:    "",
		}

		e := echo.New()
		reqBody, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(reqBody))
		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)

		if assert.NoError(t, handler.Register(c)) {
			assert.Equal(t, http.StatusBadRequest, recorder.Code)

			var response models.ErrorResponse

			_ = json.Unmarshal(recorder.Body.Bytes(), &response.Text)

			assert.Equal(t, "So sad", response.Text)
			assert.Equal(t, 123, response.Code)

		}
	})
}
