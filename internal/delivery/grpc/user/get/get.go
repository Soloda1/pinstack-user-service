package user_grpc

import (
	"context"
	"pinstack-user-service/internal/model"
)

type UserProvider interface {
	Get(ctx context.Context, id int64) (*model.User, error)
}

func GetGRPC() {
	// TODO: реализовать обработку gRPC запроса на получение пользователя
}
