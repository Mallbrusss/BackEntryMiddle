package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		errRes:      mockErrorResponse,
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

			var response map[string]any
			_ = json.Unmarshal(recorder.Body.Bytes(), &response)

			errorResp := response["error"].(map[string]any)
			assert.Equal(t, 123.0, errorResp["code"].(float64))
			assert.Equal(t, "So sad", errorResp["text"].(string))

		}
	})

	t.Run("Success Request", func(t *testing.T) {

		mockUserService.On("Register", "EgorBogachev", "Testp123!", true).
		Return(&models.User{Login: "EgorBogachev", Password: "Testp123!"}, nil)

		req := models.User{
			Login:    "EgorBogachev",
			Password: "Testp123!",
			Token:    "",
		}

		e := echo.New()
		reqBody, _ := json.Marshal(req)
		request := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(reqBody))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		c := e.NewContext(request, recorder)

		if assert.NoError(t, handler.Register(c)) {
			assert.Equal(t, http.StatusOK, recorder.Code)

			if recorder.Body.Len() == 0 {
				t.Fatalf("Response body is empty")
			}

			fmt.Println("Response body:", recorder.Body.String())

			var response map[string]any
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("Error unmarshalling response body: %v", err)
			}

			if errorField, ok := response["error"]; ok {
				t.Fatalf("Expected 'response' but got error: %v", errorField)
			}

			responseData, ok := response["response"].(map[string]any)
			if !ok {
				t.Fatalf("Expected 'response' to be a map, got %T", response["response"])
			}

			assert.Equal(t, "EgorBogachev", responseData["login"].(string))
		}
	})
}

// func TestAuthenticate(t *testing.T){
// 	mockUserService := new(MockUserService)
// 	mockErrorResponse := new(MockErrorResponse)

// }
