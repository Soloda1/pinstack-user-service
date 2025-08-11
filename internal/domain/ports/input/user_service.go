package input

import (
	"context"
	"pinstack-user-service/internal/domain/models"
)

//go:generate mockery --name UserService --dir . --output ../../../../mocks --outpkg mocks --with-expecter
type UserService interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Get(ctx context.Context, id int64) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, query string, page, limit int) ([]*models.User, int, error)
	UpdatePassword(ctx context.Context, id int64, oldPassword, newPassword string) error
	UpdateAvatar(ctx context.Context, id int64, avatarURL string) error
}
