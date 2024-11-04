package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/Mallbrusss/BackEntryMiddle/models"
)

func AuthMiddleWare(authToken string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")

			if token != authToken {
				errResp := models.ErrorResponce{
					Code: http.StatusUnauthorized,
					Text: "Invalid token",
				}
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": errResp,
				})
			}
			return next(c)
		}
	}
}
