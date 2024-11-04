package handlers

import (
	"github.com/Mallbrusss/BackEntryMiddle/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Mallbrusss/BackEntryMiddle/models"
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
			Code: 123,
			Text: "So sad",
		}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": errResp,
		})
	}

	user, err := uh.UserService.Register(req.Login, req.Password, req.Token)
	if err != nil {
		errResp := models.ErrorResponce{
			Code: 123,
			Text: "So sad",
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

	if err := c.Bind(&req); err != nil {
		errResp := models.ErrorResponce{
			Code: 123,
			Text: "So sad",
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": errResp})
	}

	user, err := uh.UserService.Authenticate(req.Login, req.Password)
	if err != nil {
		errResp := models.ErrorResponce{
			Code: 123,
			Text: "So sad",
		}
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": errResp})
	}

	c.Response().Header().Set("Authorization", user.Token)

	return c.JSON(http.StatusOK, echo.Map{
		"response": user.Token,
	})
}
