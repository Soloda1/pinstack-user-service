package output

import (
	"context"
	"pinstack-user-service/internal/domain/models"
)

//go:generate mockery --name UserRepository --dir . --output ../../../../mocks --outpkg mocks --with-expecter
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, searchQuery string, offset, pageSize int) ([]*models.User, int, error)
	UpdatePassword(ctx context.Context, id int64, newPassword string) error
	UpdateAvatar(ctx context.Context, id int64, avatarURL string) error
}
