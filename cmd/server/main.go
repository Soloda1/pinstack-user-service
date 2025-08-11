package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	user_service "pinstack-user-service/internal/application/service"
	"pinstack-user-service/internal/infrastructure/config"
	infra_grpc "pinstack-user-service/internal/infrastructure/inbound/grpc"
	user_grpc "pinstack-user-service/internal/infrastructure/inbound/grpc/user"
	"pinstack-user-service/internal/infrastructure/logger"
	user_repository "pinstack-user-service/internal/infrastructure/outbound/repository/postgres"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	userRepo := user_repository.NewUserRepository(pool, log)
	userService := user_service.NewUserService(userRepo, log)

	userGRPCApi := user_grpc.NewUserGRPCService(userService, log)
	grpcServer := infra_grpc.NewServer(userGRPCApi, cfg.GRPCServer.Address, cfg.GRPCServer.Port, log)

	metricsAddr := fmt.Sprintf("%s:%d", cfg.Prometheus.Address, cfg.Prometheus.Port)
	metricsServer := &http.Server{
		Addr:    metricsAddr,
		Handler: nil,
	}

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

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Info("Starting Prometheus metrics server", slog.String("address", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Prometheus metrics server error", slog.String("error", err.Error()))
		}
		metricsDone <- true
	}()

	<-quit
	log.Info("Shutting down servers...")

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
