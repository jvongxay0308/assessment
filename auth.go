package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Auth is a middleware that checks if the user is authenticated.
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")

		if _, err := time.Parse("January 02, 2006", authorization); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"code":    http.StatusUnauthorized,
				"message": "Unauthorized",
			})
		}

		return next(c)
	}
}
