package user_grpc

import (
	"context"
	"pinstack-user-service/internal/model"
)

type UserByUsernameProvider interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
}

func GetByUsernameGRPC() {
	// TODO: реализовать обработку gRPC запроса на получение пользователя по username
}
