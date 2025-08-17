package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	user_service "pinstack-user-service/internal/application/service"
	"pinstack-user-service/internal/infrastructure/config"
	infra_grpc "pinstack-user-service/internal/infrastructure/inbound/grpc"
	user_grpc "pinstack-user-service/internal/infrastructure/inbound/grpc/user"
	metrics_server "pinstack-user-service/internal/infrastructure/inbound/metrics"
	"pinstack-user-service/internal/infrastructure/logger"
	redis_cache "pinstack-user-service/internal/infrastructure/outbound/cache/redis"
	prometheus_metrics "pinstack-user-service/internal/infrastructure/outbound/metrics/prometheus"
	user_repository "pinstack-user-service/internal/infrastructure/outbound/repository/postgres"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.MustLoad()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DbName)
	ctx := context.Background()
	log := logger.New(cfg.Env)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error("Failed to parse postgres poolConfig", slog.String("error", err.Error()))
		os.Exit(1)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("Failed to create postgres pool", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	log.Info("Connecting to Redis",
		slog.String("address", cfg.Redis.Address),
		slog.Int("port", cfg.Redis.Port),
		slog.Int("db", cfg.Redis.DB))
	redisClient, err := redis_cache.NewClient(cfg.Redis, log)
	if err != nil {
		log.Error("Failed to create Redis client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Error("Failed to close Redis connection", slog.String("error", err.Error()))
		}
	}()

	metrics := prometheus_metrics.NewPrometheusMetricsProvider()

	metrics.SetServiceHealth(true)

	userCache := redis_cache.NewUserCache(redisClient, log, metrics)

	userRepo := user_repository.NewUserRepository(pool, log, metrics)
	originalUserService := user_service.NewUserService(userRepo, log, metrics)

	userService := user_service.NewUserServiceCacheDecorator(
		originalUserService,
		userCache,
		log,
		metrics,
	)

	userGRPCApi := user_grpc.NewUserGRPCService(userService, log)
	grpcServer := infra_grpc.NewServer(userGRPCApi, cfg.GRPCServer.Address, cfg.GRPCServer.Port, log, metrics)

	metricsServer := metrics_server.NewMetricsServer(cfg.Prometheus.Address, cfg.Prometheus.Port, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	done := make(chan bool, 1)
	metricsDone := make(chan bool, 1)

	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Error("gRPC server error", slog.String("error", err.Error()))
		}
		done <- true
	}()

	go func() {
		if err := metricsServer.Run(); err != nil {
			log.Error("Metrics server error", slog.String("error", err.Error()))
		}
		metricsDone <- true
	}()

	<-quit
	log.Info("Shutting down servers...")

	metrics.SetServiceHealth(false)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := grpcServer.Shutdown(); err != nil {
		log.Error("gRPC server shutdown error", slog.String("error", err.Error()))
	}

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Error("Metrics server shutdown error", slog.String("error", err.Error()))
	}

	<-done
	<-metricsDone

	log.Info("Server exited")
}
