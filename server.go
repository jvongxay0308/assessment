package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	zlog := zap.Must(zap.NewDevelopment())
	defer zlog.Sync()

	// Replace the global logger with the one we just created
	// so that we can use the zap.L() function to get the logger
	// anywhere in the code.
	zap.ReplaceGlobals(zlog)

	// parse the database URL from the environment
	dbURL, err := url.Parse(GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost/postgres?sslmode=disable"))
	if err != nil {
		zlog.Fatal("parse the database URL", zap.Error(err))
	}

	// try to migrate the database to the latest version
	// before starting the server
	if err := tryToMigrate(dbURL.String()); err != nil {
		zlog.Fatal("migrate the database",
			zap.String("dbURL", dbURL.Redacted()), // Redacted password from the URL
			zap.Error(err),
		)
	}

	dbConn, err := sql.Open("postgres", dbURL.String())
	if err != nil {
		zlog.Fatal("open the database connection", zap.Error(err))
	}
	defer dbConn.Close()

	e := echo.New()
	e.Use(Auth)

	db := New(dbConn)
	handler := NewHandler(db)
	handler.Install(e)

	// Start the server
	errC := make(chan error, 1)
	go func() {
		errC <- e.Start(":" + GetEnv("PORT", "2565"))
	}()

	// Wait for an interrupt or kill signal
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	// Shutdown the server gracefully
	select {
	case err := <-errC:
		if err != nil && err != http.ErrServerClosed {
			zlog.Fatal("Server failed to start", zap.Error(err))
		}
		zlog.Info("Server gracefully stopped")

	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		zlog.Info("Shutting down the server")
		if err := e.Shutdown(ctx); err != nil {
			zlog.Fatal("Server failed to shutdown", zap.Error(err))
		}
		zlog.Info("Server gracefully stopped")
	}
}

// GetEnv looks up the given key from the environment, returning its value if
// it exists, and otherwise returning the given fallback value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
