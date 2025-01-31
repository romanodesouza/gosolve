package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/romanodesouza/gosolve/internal/api"
	"github.com/romanodesouza/gosolve/internal/search"
	"github.com/spf13/viper"
)

func main() {
	// Set up config
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Set up config defaults
	viper.SetDefault("LOG_LEVEL", "DEBUG")
	viper.SetDefault("SERVICE_PORT", "8080")

	// Set up logger
	opts := &slog.HandlerOptions{
		Level: func() slog.Level {
			level := viper.GetString("LOG_LEVEL")
			switch strings.ToUpper(level) {
			case "DEBUG":
				return slog.LevelDebug
			case "INFO":
				return slog.LevelInfo
			case "WARN":
				return slog.LevelWarn
			case "ERROR":
				return slog.LevelError
			}

			panic(fmt.Errorf(`invalid log level: "%s", please choose between DEBUG, INFO, WARN or ERROR`, level))
		}(),
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))

	// Set up search index
	filePath := "input.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("could not read input file", slog.Any("err", err), slog.String("file", filePath))
		return
	}
	idx, err := search.NewIndex(f, *logger)
	if err != nil {
		logger.Error("could not ingest input file", slog.Any("err", err))
		return
	}

	// Set up web server
	mux := http.NewServeMux()
	api.NewSearchHandler(*logger, idx).AssignRoutes(mux)

	addr := fmt.Sprintf(":%s", viper.GetString("SERVICE_PORT"))
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Info("starting web server", slog.String("service_port", addr))
		}
	}()
	<-done
	logger.Info("server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", slog.Any("err", err))
	}
}
