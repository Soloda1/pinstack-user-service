package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"pinstack-user-service/config"
	user_grpc "pinstack-user-service/internal/delivery/grpc/user"
	"pinstack-user-service/internal/logger"
	"pinstack-user-service/internal/repository"
	user_repository "pinstack-user-service/internal/repository/user"
	user_service "pinstack-user-service/internal/service/user"
	"syscall"
	"time"
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

	storage, err := repository.NewStorage(ctx, dsn)
	if err != nil {
		log.Debug("Failed to create storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	userRepo := user_repository.NewUserRepository(storage, log)
	userService := user_service.NewUserService(userRepo, log)

	grpcApi := user_grpc.NewGRPCServer(userService, log)
	grpcServer := user_grpc.NewServer(grpcApi, cfg.GRPCServer.Address, cfg.GRPCServer.Port, log)

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
	storage.Close()
	<-done
	log.Info("Server exiting")
}
