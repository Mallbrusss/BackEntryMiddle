package handlers

import (
	"internal/service"
	"net/http"

	"github.com/labstack/echo/v4"

	"internal/models"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandlers(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (uh UserHandler) Register(c echo.Context) error {
	var req models.User

	if err := c.Bind(&req); err != nil {

		errResp := models.ErrorResponce{
			Code: http.StatusBadRequest,
			Text: "Invalid request",
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": errResp,
		})
	}

	user, err := uh.UserService.Register(req.Login, req.Password, req.Token)
	if err != nil{
		errResp := models.ErrorResponce{
			Code: http.StatusBadRequest,
			Text: "Invalid request",
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": errResp,
		})
	}

	//FIXME: поправить вывод
	return c.JSON(http.StatusOK, echo.Map{
		"responce": user.Login,
	})
}

func (uh UserHandler) Authenticate(c echo.Context) error {
	var req models.User

	if err := c.Bind(&req); err != nil{
		errResp := models.ErrorResponce{
			Code: http.StatusBadRequest,
			Text: "Invalid request",
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": errResp})
	}

	user, err := uh.UserService.Authenticate(req.Login, req.Password)
	if err != nil {
		errResp := models.ErrorResponce{
			Code: http.StatusUnauthorized,
			Text: "Authentication failed",
		}
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": errResp})
    }

	c.Request().Header.Set("Authorization", user.Token)

	return c.JSON(http.StatusOK, echo.Map{
		"responce": user.Token,
	})
}
