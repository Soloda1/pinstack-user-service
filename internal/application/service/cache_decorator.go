package service

import (
	"context"
	"errors"
	"log/slog"

	"pinstack-user-service/internal/domain/models"
	input "pinstack-user-service/internal/domain/ports/input"
	output "pinstack-user-service/internal/domain/ports/output"
	"pinstack-user-service/internal/domain/ports/output/cache"

	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
)

type UserServiceCacheDecorator struct {
	service   input.UserService
	userCache cache.UserCache
	log       output.Logger
	metrics   output.MetricsProvider
}

func NewUserServiceCacheDecorator(
	service input.UserService,
	userCache cache.UserCache,
	log output.Logger,
	metrics output.MetricsProvider,
) input.UserService {
	return &UserServiceCacheDecorator{
		service:   service,
		userCache: userCache,
		log:       log,
		metrics:   metrics,
	}
}

func (d *UserServiceCacheDecorator) Create(ctx context.Context, user *models.User) (*models.User, error) {
	d.log.Debug("Creating user with cache decorator",
		slog.String("username", user.Username),
		slog.String("email", user.Email))

	result, err := d.service.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := d.userCache.SetUser(ctx, result); err != nil {
		d.log.Warn("Failed to cache created user",
			slog.Int64("user_id", result.ID),
			slog.String("error", err.Error()))
	}

	return result, nil
}

func (d *UserServiceCacheDecorator) Get(ctx context.Context, id int64) (*models.User, error) {
	d.log.Debug("Getting user by ID with cache decorator", slog.Int64("user_id", id))

	cachedUser, err := d.userCache.GetUserByID(ctx, id)
	if err == nil {
		d.log.Debug("User found in cache", slog.Int64("user_id", id))
		d.metrics.IncrementCacheHits()
		return cachedUser, nil
	}

	if !errors.Is(err, custom_errors.ErrCacheMiss) {
		d.log.Warn("Failed to get user from cache",
			slog.Int64("user_id", id),
			slog.String("error", err.Error()))
	} else {
		d.metrics.IncrementCacheMisses()
	}

	d.log.Debug("User cache miss, fetching from service", slog.Int64("user_id", id))
	user, err := d.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := d.userCache.SetUser(ctx, user); err != nil {
		d.log.Warn("Failed to cache user",
			slog.Int64("user_id", id),
			slog.String("error", err.Error()))
	}

	return user, nil
}

func (d *UserServiceCacheDecorator) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	d.log.Debug("Getting user by username with cache decorator", slog.String("username", username))

	cachedUser, err := d.userCache.GetUserByUsername(ctx, username)
	if err == nil {
		d.log.Debug("User found in cache by username", slog.String("username", username))
		d.metrics.IncrementCacheHits()
		return cachedUser, nil
	}

	if !errors.Is(err, custom_errors.ErrCacheMiss) {
		d.log.Warn("Failed to get user by username from cache",
			slog.String("username", username),
			slog.String("error", err.Error()))
	} else {
		d.metrics.IncrementCacheMisses()
	}

	d.log.Debug("User username cache miss, fetching from service", slog.String("username", username))
	user, err := d.service.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if err := d.userCache.SetUser(ctx, user); err != nil {
		d.log.Warn("Failed to cache user by username",
			slog.String("username", username),
			slog.String("error", err.Error()))
	}

	return user, nil
}

func (d *UserServiceCacheDecorator) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	d.log.Debug("Getting user by email with cache decorator", slog.String("email", email))

	cachedUser, err := d.userCache.GetUserByEmail(ctx, email)
	if err == nil {
		d.log.Debug("User found in cache by email", slog.String("email", email))
		d.metrics.IncrementCacheHits()
		return cachedUser, nil
	}

	if !errors.Is(err, custom_errors.ErrCacheMiss) {
		d.log.Warn("Failed to get user by email from cache",
			slog.String("email", email),
			slog.String("error", err.Error()))
	} else {
		d.metrics.IncrementCacheMisses()
	}

	d.log.Debug("User email cache miss, fetching from service", slog.String("email", email))
	user, err := d.service.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := d.userCache.SetUser(ctx, user); err != nil {
		d.log.Warn("Failed to cache user by email",
			slog.String("email", email),
			slog.String("error", err.Error()))
	}

	return user, nil
}

func (d *UserServiceCacheDecorator) Update(ctx context.Context, user *models.User) (*models.User, error) {
	d.log.Debug("Updating user with cache decorator",
		slog.Int64("user_id", user.ID),
		slog.String("username", user.Username))

	oldUser, err := d.service.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	updatedUser, err := d.service.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if err := d.userCache.DeleteUser(ctx, oldUser); err != nil {
		d.log.Warn("Failed to invalidate old user cache after update",
			slog.Int64("user_id", oldUser.ID),
			slog.String("old_username", oldUser.Username),
			slog.String("old_email", oldUser.Email),
			slog.String("error", err.Error()))
	}

	if err := d.userCache.SetUser(ctx, updatedUser); err != nil {
		d.log.Warn("Failed to cache updated user",
			slog.Int64("user_id", updatedUser.ID),
			slog.String("error", err.Error()))
	}

	return updatedUser, nil
}

func (d *UserServiceCacheDecorator) Delete(ctx context.Context, id int64) error {
	d.log.Debug("Deleting user with cache decorator", slog.Int64("user_id", id))

	user, err := d.service.Get(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrUserNotFound) {
			if cacheErr := d.userCache.DeleteUserByID(ctx, id); cacheErr != nil {
				d.log.Warn("Failed to invalidate user cache by ID after deletion attempt",
					slog.Int64("user_id", id),
					slog.String("error", cacheErr.Error()))
			}
		}
		return err
	}

	err = d.service.Delete(ctx, id)
	if err != nil {
		return err
	}

	if err := d.userCache.DeleteUser(ctx, user); err != nil {
		d.log.Warn("Failed to invalidate user cache after deletion",
			slog.Int64("user_id", id),
			slog.String("error", err.Error()))
	}

	return nil
}

func (d *UserServiceCacheDecorator) Search(ctx context.Context, query string, page, limit int) ([]*models.User, int, error) {
	d.log.Debug("Searching users with cache decorator",
		slog.String("query", query),
		slog.Int("page", page),
		slog.Int("limit", limit))

	users, count, err := d.service.Search(ctx, query, page, limit)
	if err != nil {
		return nil, 0, err
	}

	for _, user := range users {
		if err := d.userCache.SetUser(ctx, user); err != nil {
			d.log.Warn("Failed to cache user from search results",
				slog.Int64("user_id", user.ID),
				slog.String("username", user.Username),
				slog.String("error", err.Error()))
		}
	}

	return users, count, nil
}

func (d *UserServiceCacheDecorator) UpdatePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	d.log.Debug("Updating user password with cache decorator", slog.Int64("user_id", id))

	err := d.service.UpdatePassword(ctx, id, oldPassword, newPassword)
	if err != nil {
		return err
	}

	if err := d.userCache.DeleteUserByID(ctx, id); err != nil {
		d.log.Warn("Failed to invalidate user cache after password update",
			slog.Int64("user_id", id),
			slog.String("error", err.Error()))
	}

	return nil
}

func (d *UserServiceCacheDecorator) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	d.log.Debug("Updating user avatar with cache decorator",
		slog.Int64("user_id", id),
		slog.String("avatar_url", avatarURL))

	err := d.service.UpdateAvatar(ctx, id, avatarURL)
	if err != nil {
		return err
	}

	if err := d.userCache.DeleteUserByID(ctx, id); err != nil {
		d.log.Warn("Failed to invalidate user cache after avatar update",
			slog.Int64("user_id", id),
			slog.String("error", err.Error()))
	}

	return nil
}
