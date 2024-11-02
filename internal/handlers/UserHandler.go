package handlers

import (
	"internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandlers(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h UserHandler) Register(c echo.Context) error {

	return c.JSON(http.StatusOK, echo.Map{
		"responce": "user",
	})
}

func (h UserHandler) Authenticate (c echo.Context) error {

	return c.JSON(http.StatusOK, echo.Map{
		"token": "token",
	})
}
