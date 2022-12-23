package main

import (
	"errors"
	"net/http"
	"strconv"

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
	e.GET("/expenses", h.List)
	e.GET("/expenses/:id", h.Get)
	e.PUT("/expenses/:id", h.Update)
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

func (h *Handler) Get(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	ctx := c.Request().Context()
	expense, err := h.db.Get(ctx, id)

	switch {
	case err == nil:
		return c.JSON(http.StatusOK, expense)

	case errors.Is(err, ErrClosed):
		return echo.NewHTTPError(http.StatusServiceUnavailable, echo.Map{
			"code":    http.StatusServiceUnavailable,
			"message": err.Error(),
		})

	case errors.Is(err, ErrNoExpense):
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{
			"code":    http.StatusNotFound,
			"message": err.Error(),
		})

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}
}

func (h *Handler) List(c echo.Context) error {
	ctx := c.Request().Context()
	expenses, err := h.db.List(ctx)

	switch {
	case err == nil:
		return c.JSON(http.StatusOK, expenses)

	case errors.Is(err, ErrClosed):
		return echo.NewHTTPError(http.StatusServiceUnavailable, echo.Map{
			"code":    http.StatusServiceUnavailable,
			"message": err.Error(),
		})

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}
}

func (h *Handler) Update(c echo.Context) error {
	expense := &Expense{}
	if err := c.Bind(expense); err != nil {
		return err
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	expense.ID = id

	ctx := c.Request().Context()
	expense, err = h.db.Update(ctx, expense)

	switch {
	case err == nil:
		return c.JSON(http.StatusOK, expense)

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

	case errors.Is(err, ErrNoExpense):
		return echo.NewHTTPError(http.StatusNotFound, echo.Map{
			"code":    http.StatusNotFound,
			"message": err.Error(),
		})

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
	}
}
