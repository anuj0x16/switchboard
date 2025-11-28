package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/anuj0x16/switchboard/internal/env"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	httpPort    int
	databaseUrl string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	cfg := config{
		httpPort:    env.GetInt("HTTP_PORT", 4000),
		databaseUrl: env.GetString("DATABASE_URL", ""),
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

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

	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), cfg.databaseUrl)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
