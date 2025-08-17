package cache

import (
	"context"
	"pinstack-user-service/internal/domain/models"
)

//go:generate mockery --name UserCache --dir . --output ../../../../mocks/cache --outpkg mocks --with-expecter --filename UserCache.go
type UserCache interface {
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	SetUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, user *models.User) error
	DeleteUserByID(ctx context.Context, userID int64) error
}
