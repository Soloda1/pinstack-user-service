package user_grpc

import (
	"context"
	user_service "pinstack-user-service/internal/service/user"
)

type UserUpdater interface {
	Update(ctx context.Context, user *user_service.Service) error // примерная сигнатура, уточни по необходимости
}

func UpdateGRPC() {
	// TODO: реализовать обработку gRPC запроса на обновление пользователя
}
