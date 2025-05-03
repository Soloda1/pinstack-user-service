package user_grpc

import (
	"context"
)

type UserDeleter interface {
	Delete(ctx context.Context, id int64) error
}

func DeleteGRPC() {
	// TODO: реализовать обработку gRPC запроса на удаление пользователя
}
