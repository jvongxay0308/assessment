package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAuth(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	next := func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"code":    http.StatusOK,
			"message": "OK",
		})
	}

	err := Auth(next)(c)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	req.Header.Set("Authorization", "January 02, 2006")
	rec = httptest.NewRecorder()
	err2 := Auth(next)(c)
	if err2 != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}
}
