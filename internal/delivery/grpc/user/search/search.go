package user_grpc

import (
	"context"
	"pinstack-user-service/internal/model"
)

type UserSearcher interface {
	Search(ctx context.Context, query string, page, limit int) ([]*model.User, int, error)
}

func SearchGRPC() {
	// TODO: реализовать обработку gRPC запроса на поиск пользователей
}
