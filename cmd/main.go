package main

import (
	"context"
	"errors"
	"fmt"
	"gravitum/internal/infra/postgres"
	"gravitum/internal/presentation/restapi"
	"gravitum/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	background := context.Background()

	process, shutdown := signal.NotifyContext(background, syscall.SIGINT, syscall.SIGTERM)
	defer shutdown()

	pool, err := pgxpool.New(process, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepository(pool)
	srv := service.New(userRepo)
	handler := restapi.NewHandler(logger, srv)

	httpServer := http.Server{
		Addr:              ":8080",
		Handler:           handler.Routes(),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second, // protection against Slowloris DDoS attack
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       time.Minute,
	}

	go func() {
		logger.LogAttrs(background, slog.LevelInfo, "http.Server", slog.String("status", "running"), slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.LogAttrs(background, slog.LevelError, "http.Server", slog.String("status", "error"), slog.String("error", err.Error()))
			shutdown()
		}
	}()

	defer func() {
		logger.LogAttrs(background, slog.LevelInfo, "http.Server", slog.String("status", "shutdown"), slog.String("addr", httpServer.Addr))
		await, timeout := context.WithTimeout(background, 10*time.Second)
		defer timeout()
		if err := httpServer.Shutdown(await); err != nil {
			logger.LogAttrs(background, slog.LevelError, "http.Server.Shutdown", slog.String("status", "error"), slog.String("error", err.Error()))
		}
	}()

	<-process.Done()
	logger.LogAttrs(background, slog.LevelInfo, "application", slog.String("status", "interrupt"))
}
