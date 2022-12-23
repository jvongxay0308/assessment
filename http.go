package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *DB
}

func NewHandler(db *DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Install(e *echo.Echo) {
	e.POST("/expenses", h.Create)
}

func (h *Handler) Create(c echo.Context) error {
	expense := &Expense{}
	if err := c.Bind(expense); err != nil {
		return err
	}

	ctx := c.Request().Context()
	expense, err := h.db.Create(ctx, expense)

	switch {
	case err == nil:
		return c.JSON(http.StatusCreated, expense)

	case errors.Is(err, ErrClosed):
		return echo.NewHTTPError(http.StatusServiceUnavailable, echo.Map{
			"code":    http.StatusServiceUnavailable,
			"message": err.Error(),
		})

	case errors.Is(err, ErrInvalidExpense):
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})

	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}
}
