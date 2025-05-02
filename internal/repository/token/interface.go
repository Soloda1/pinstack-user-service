package token

import (
	"context"
	"time"

	"pinstack-user-service/internal/model"
)

type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	GetRefreshTokenByJTI(ctx context.Context, jti string) (*model.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByJTI(ctx context.Context, jti string) error
	DeleteUserRefreshTokens(ctx context.Context, userID int64) error
	DeleteExpiredTokens(ctx context.Context, before time.Time) error
}
