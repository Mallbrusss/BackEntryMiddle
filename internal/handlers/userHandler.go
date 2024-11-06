package handlers

import (
	"net/http"
	"os"

	"github.com/Mallbrusss/BackEntryMiddle/internal/service"

	"github.com/labstack/echo/v4"

	"github.com/Mallbrusss/BackEntryMiddle/models"
)

type UserHandler struct {
	UserService service.UserServiceInterface
	errRes      *models.ErrorResponse
}

func NewUserHandlers(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
		errRes:      models.NewErrorResponse()}
}

func (uh UserHandler) Register(c echo.Context) error {
	var req models.User
	isAdmin := false
	adminToken := os.Getenv("ADMIN_TOKEN")
	if err := c.Bind(&req); err != nil {

		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{"error": uh.errRes.GetErrorResponse(http.StatusBadRequest)},
		})
	}

	if req.Token == adminToken {
		isAdmin = true
	}

	user, err := uh.UserService.Register(req.Login, req.Password, isAdmin)
	if err != nil {

		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": uh.errRes.GetErrorResponse(http.StatusBadRequest),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"response": map[string]string{
			"login": user.Login,
		}})
}

func (uh UserHandler) Authenticate(c echo.Context) error {
	var req models.User

	if err := c.Bind(&req); err != nil {

		return c.JSON(http.StatusBadRequest, echo.Map{"error": uh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	user, err := uh.UserService.Authenticate(req.Login, req.Password)
	if err != nil {

		return c.JSON(http.StatusUnauthorized, echo.Map{"error": uh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	c.Response().Header().Set("Authorization", user.Token)

	return c.JSON(http.StatusOK, echo.Map{
		"response": map[string]string{
			"token": user.Token,
		}})
}

func (uh *UserHandler) Logout(c echo.Context) error {

	token := c.Param("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": uh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	err := uh.UserService.DeleteToken(token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": uh.errRes.GetErrorResponse(http.StatusBadRequest)})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"response": map[string]string{
			token: "true",
		}})
}
