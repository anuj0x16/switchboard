package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/anuj0x16/switchboard/internal/env"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	httpPort int
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	cfg := config{
		httpPort: env.GetInt("HTTP_PORT", 4000),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", cfg.httpPort),
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting http server", "port", cfg.httpPort)

	err := srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error(), "port", cfg.httpPort)
		os.Exit(1)
	}
}
