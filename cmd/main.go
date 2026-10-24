package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/N1ktarchik/habit-graph/internal/core/logger"
	storage "github.com/N1ktarchik/habit-graph/internal/core/repository/postgres"
	cache "github.com/N1ktarchik/habit-graph/internal/core/repository/redis"
	"github.com/N1ktarchik/habit-graph/internal/core/transport/middleware"
	"github.com/N1ktarchik/habit-graph/internal/core/transport/server"
	"github.com/N1ktarchik/habit-graph/internal/habits/repository/postgres"
	"github.com/N1ktarchik/habit-graph/internal/habits/repository/redis"
	"github.com/N1ktarchik/habit-graph/internal/habits/service"
	transport "github.com/N1ktarchik/habit-graph/internal/habits/transportHTTP"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("godotenv: .env file not found, reading from system env")
	}

	lvlStr := os.Getenv("LOG_LEVEL")
	var lvl slog.Level
	switch lvlStr {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	log := logger.Setup(lvl)
	slog.SetDefault(log)
	log.Info("starting habbit-graph app", slog.String("log_level", lvl.String()))

	initCtx, initCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer initCancel()

	connStr := storage.GetPostgresValues()
	pgConfig := storage.NewPostgresConfig(connStr, 25, 5, 30*time.Minute, 5*time.Minute, 1*time.Minute)

	pool, err := storage.CreatePool(initCtx, pgConfig, log)
	if err != nil {
		log.Error("failed to initialize database pool", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		log.Info("closing postgres connection pool")
		pool.Close()
	}()
	log.Info("postgres connection pool established")

	rHost := os.Getenv("REDIS_HOST")
	if rHost == "" {
		rHost = "localhost"
	}
	rPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil || rPort <= 0 {
		rPort = 6379
	}
	rDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil || rDB < 0 {
		rDB = 0
	}
	rPass := os.Getenv("REDIS_PASSWORD")

	redisCfg := cache.NewRedisConfig(rHost, rPort, rPass, rDB)
	rClient, err := cache.New(redisCfg, log)
	if err != nil {
		log.Error("failed to initialize redis client", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		log.Info("closing redis client")
		if err := rClient.Close(); err != nil {
			log.Error("failed to close redis client", slog.Any("error", err))
		}
	}()
	log.Info("redis connection established")

	habitsRepo := postgres.NewRepository(log, pool)
	habitsCache := redis.NewCache(log, rClient)
	habitsService := service.NewService(log, habitsRepo, habitsCache)
	habitsHandler := transport.NewTransportHTTP(log, habitsService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/habits", habitsHandler.CreateHabit)
	mux.HandleFunc("GET /api/habits", habitsHandler.GetHabits)
	mux.HandleFunc("PATCH /api/habits/{id}", habitsHandler.UpdateHabit)
	mux.HandleFunc("DELETE /api/habits/{id}", habitsHandler.DeleteHabit)

	mux.HandleFunc("POST /api/habits/{id}/activities", habitsHandler.TrackActivity)
	mux.HandleFunc("POST /api/habits/{id}/track", habitsHandler.TrackActivity)
	mux.HandleFunc("GET /api/habits/{id}/graph", habitsHandler.GetGraph)

	allowedOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

	handlerWithCORS := middleware.CorsMiddleware(mux, allowedOrigin)

	srv := server.NewServer(handlerWithCORS, log)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.Start(addr); err != nil {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Info("received shutdown signal", slog.String("signal", sig.String()))
	case err := <-serverErr:
		log.Error("server crashed", slog.Any("error", err))
		os.Exit(1)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", slog.Any("error", err))
	}

	log.Info("server exited cleanly")
}
