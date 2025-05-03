package user_service

import (
	"context"
	"pinstack-user-service/internal/model"
)

type UserService interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, query string, page, limit int) ([]*model.User, int, error)
	UpdatePassword(ctx context.Context, id int64, password string) error
	UpdateAvatar(ctx context.Context, id int64, avatarURL string) error
}
