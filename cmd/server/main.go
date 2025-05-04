package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"pinstack-user-service/config"
	"pinstack-user-service/internal/delivery/grpc"
	user_grpc "pinstack-user-service/internal/delivery/grpc/user"
	"pinstack-user-service/internal/logger"
	user_repository "pinstack-user-service/internal/repository/user/postgres"
	user_service "pinstack-user-service/internal/service/user"
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

	userRepo := user_repository.NewUserRepository(pool, log)
	userService := user_service.NewUserService(userRepo, log)

	userGRPCApi := user_grpc.NewUserGRPCService(userService, log)
	grpcServer := grpc.NewServer(userGRPCApi, cfg.GRPCServer.Address, cfg.GRPCServer.Port, log)

	done := make(chan bool)
	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Error("gRPC server error", slog.String("error", err.Error()))
		}
		done <- true
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down gRPC server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := grpcServer.Shutdown(); err != nil {
		log.Error("gRPC server shutdown error", slog.String("error", err.Error()))
	}
	<-done
	log.Info("Server exiting")
}
