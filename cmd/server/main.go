package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"pinstack-user-service/config"
	"pinstack-user-service/internal/logger"
	"pinstack-user-service/internal/repository"
	user_repository "pinstack-user-service/internal/repository/user"
	user_service "pinstack-user-service/internal/service/user"
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
	_ = userService
}
