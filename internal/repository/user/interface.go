package user_repository

import (
	"context"

	"pinstack-user-service/internal/model"
)

//go:generate mockery --name UserRepository --dir . --output ../../../mocks --outpkg mocks --with-expecter
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, searchQuery string, offset, pageSize int) ([]*model.User, int, error)
	UpdatePassword(ctx context.Context, id int64, password string) error
	UpdateAvatar(ctx context.Context, id int64, avatarURL string) error
}
