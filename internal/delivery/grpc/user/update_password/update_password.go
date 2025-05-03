package user_grpc

import (
	"context"
)

type UserPasswordUpdater interface {
	UpdatePassword(ctx context.Context, id int64, password string) error
}

func UpdatePasswordGRPC() {
	// TODO: реализовать обработку gRPC запроса на обновление пароля пользователя
}
