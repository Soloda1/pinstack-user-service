package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"pinstack-user-service/internal/domain/models"
	ports "pinstack-user-service/internal/domain/ports/output"

	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
)

const (
	userCacheKeyPrefix         = "user:"
	userEmailCacheKeyPrefix    = "user:email:"
	userUsernameCacheKeyPrefix = "user:username:"
	userCacheTTL               = 30 * time.Minute
)

type UserCache struct {
	client  *Client
	log     ports.Logger
	metrics ports.MetricsProvider
}

func NewUserCache(client *Client, log ports.Logger, metrics ports.MetricsProvider) *UserCache {
	return &UserCache{
		client:  client,
		log:     log,
		metrics: metrics,
	}
}

func (u *UserCache) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	start := time.Now()
	key := u.getUserKey(userID)

	var user models.User
	err := u.client.Get(ctx, key, &user)

	duration := time.Since(start)
	u.metrics.RecordCacheOperationDuration("get", duration)

	if err != nil {
		if errors.Is(err, custom_errors.ErrCacheMiss) {
			u.log.Debug("User cache miss", slog.Int64("user_id", userID))
			return nil, custom_errors.ErrCacheMiss
		}
		u.log.Error("Failed to get user from cache",
			slog.Int64("user_id", userID),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user from cache: %w", err)
	}

	u.log.Debug("User cache hit", slog.Int64("user_id", userID))
	return &user, nil
}

func (u *UserCache) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	start := time.Now()
	key := u.getUserEmailKey(email)

	var user models.User
	err := u.client.Get(ctx, key, &user)

	duration := time.Since(start)
	u.metrics.RecordCacheOperationDuration("get", duration)

	if err != nil {
		if errors.Is(err, custom_errors.ErrCacheMiss) {
			u.log.Debug("User email cache miss", slog.String("email", email))
			return nil, custom_errors.ErrCacheMiss
		}
		u.log.Error("Failed to get user by email from cache",
			slog.String("email", email),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user by email from cache: %w", err)
	}

	u.log.Debug("User email cache hit", slog.String("email", email))
	return &user, nil
}

func (u *UserCache) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	start := time.Now()
	key := u.getUserUsernameKey(username)

	var user models.User
	err := u.client.Get(ctx, key, &user)

	duration := time.Since(start)
	u.metrics.RecordCacheOperationDuration("get", duration)

	if err != nil {
		if errors.Is(err, custom_errors.ErrCacheMiss) {
			u.log.Debug("User username cache miss", slog.String("username", username))
			return nil, custom_errors.ErrCacheMiss
		}
		u.log.Error("Failed to get user by username from cache",
			slog.String("username", username),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get user by username from cache: %w", err)
	}

	u.log.Debug("User username cache hit", slog.String("username", username))
	return &user, nil
}

func (u *UserCache) SetUser(ctx context.Context, user *models.User) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		u.metrics.RecordCacheOperationDuration("set", duration)
	}()

	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	idKey := u.getUserKey(user.ID)
	if err := u.client.Set(ctx, idKey, user, userCacheTTL); err != nil {
		u.log.Error("Failed to set user cache by ID",
			slog.Int64("user_id", user.ID),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to set user cache by ID: %w", err)
	}

	emailKey := u.getUserEmailKey(user.Email)
	if err := u.client.Set(ctx, emailKey, user, userCacheTTL); err != nil {
		u.log.Error("Failed to set user cache by email",
			slog.String("email", user.Email),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to set user cache by email: %w", err)
	}

	usernameKey := u.getUserUsernameKey(user.Username)
	if err := u.client.Set(ctx, usernameKey, user, userCacheTTL); err != nil {
		u.log.Error("Failed to set user cache by username",
			slog.String("username", user.Username),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to set user cache by username: %w", err)
	}

	u.log.Debug("User cached successfully",
		slog.Int64("user_id", user.ID),
		slog.String("username", user.Username),
		slog.String("email", user.Email),
		slog.Duration("ttl", userCacheTTL))
	return nil
}

func (u *UserCache) DeleteUser(ctx context.Context, user *models.User) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		u.metrics.RecordCacheOperationDuration("delete", duration)
	}()

	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	idKey := u.getUserKey(user.ID)
	if err := u.client.Delete(ctx, idKey); err != nil {
		u.log.Error("Failed to delete user from cache by ID",
			slog.Int64("user_id", user.ID),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete user from cache by ID: %w", err)
	}

	emailKey := u.getUserEmailKey(user.Email)
	if err := u.client.Delete(ctx, emailKey); err != nil {
		u.log.Error("Failed to delete user from cache by email",
			slog.String("email", user.Email),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete user from cache by email: %w", err)
	}

	usernameKey := u.getUserUsernameKey(user.Username)
	if err := u.client.Delete(ctx, usernameKey); err != nil {
		u.log.Error("Failed to delete user from cache by username",
			slog.String("username", user.Username),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete user from cache by username: %w", err)
	}

	u.log.Debug("User deleted from cache",
		slog.Int64("user_id", user.ID),
		slog.String("username", user.Username),
		slog.String("email", user.Email))
	return nil
}

func (u *UserCache) DeleteUserByID(ctx context.Context, userID int64) error {
	user, err := u.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, custom_errors.ErrCacheMiss) {
			u.log.Debug("User not in cache, nothing to delete", slog.Int64("user_id", userID))
			return nil
		}
		u.log.Warn("Failed to get user for full cache deletion, deleting by ID only",
			slog.Int64("user_id", userID),
			slog.String("error", err.Error()))
	}

	if user != nil {
		return u.DeleteUser(ctx, user)
	}

	start := time.Now()
	defer func() {
		duration := time.Since(start)
		u.metrics.RecordCacheOperationDuration("delete", duration)
	}()

	idKey := u.getUserKey(userID)
	if err := u.client.Delete(ctx, idKey); err != nil {
		u.log.Error("Failed to delete user from cache by ID",
			slog.Int64("user_id", userID),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete user from cache by ID: %w", err)
	}

	u.log.Debug("User deleted from cache by ID", slog.Int64("user_id", userID))
	return nil
}

func (u *UserCache) getUserKey(userID int64) string {
	return userCacheKeyPrefix + strconv.FormatInt(userID, 10)
}

func (u *UserCache) getUserEmailKey(email string) string {
	return userEmailCacheKeyPrefix + email
}

func (u *UserCache) getUserUsernameKey(username string) string {
	return userUsernameCacheKeyPrefix + username
}
