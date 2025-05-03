package user_grpc

import (
	"context"
	"pinstack-user-service/internal/model"
)

type UserCreator interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
}

func CreateGRPC() {
	// TODO: реализовать обработку gRPC запроса на создание пользователя
}
