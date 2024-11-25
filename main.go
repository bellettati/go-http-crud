package main

import (
	"go-http-crud/api"
	"go-http-crud/database"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err)
		os.Exit(1)
	}

	slog.Info("all systems offline")
}

func run() error {
	application := database.New()
	handler := api.NewHandler(application)

	s := http.Server{
		Addr: ":8080",
		Handler: handler,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	slog.Info("application running on port 8080")

	return nil
}