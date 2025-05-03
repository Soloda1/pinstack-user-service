package user_grpc

import (
	"context"
)

type UserAvatarUpdater interface {
	UpdateAvatar(ctx context.Context, id int64, avatarURL string) error
}

func UpdateAvatarGRPC() {
	// TODO: реализовать обработку gRPC запроса на обновление аватара пользователя
}
