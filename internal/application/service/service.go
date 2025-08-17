package service

import (
	"context"
	"errors"
	"log/slog"
	"pinstack-user-service/internal/domain/models"
	input "pinstack-user-service/internal/domain/ports/input"
	output "pinstack-user-service/internal/domain/ports/output"

	"github.com/jackc/pgx/v5"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
)

type Service struct {
	repo    output.UserRepository
	log     output.Logger
	metrics output.MetricsProvider
}

func NewUserService(repo output.UserRepository, log output.Logger, metrics output.MetricsProvider) input.UserService {
	return &Service{repo: repo, log: log, metrics: metrics}
}

func (s *Service) Create(ctx context.Context, user *models.User) (*models.User, error) {
	s.log.Debug("Creating user",
		slog.String("username", user.Username),
		slog.String("email", user.Email))

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		s.metrics.IncrementUserOperations("create", false)
		switch {
		case errors.Is(err, custom_errors.ErrUsernameExists):
			s.log.Debug("Username already exists",
				slog.String("username", user.Username))
			return nil, custom_errors.ErrUsernameExists
		case errors.Is(err, custom_errors.ErrEmailExists):
			s.log.Debug("Email already exists",
				slog.String("email", user.Email))
			return nil, custom_errors.ErrEmailExists
		default:
			s.log.Error("Failed to create user",
				slog.String("error", err.Error()),
				slog.String("username", user.Username),
				slog.String("email", user.Email))
			return nil, custom_errors.ErrDatabaseQuery
		}
	}

	s.metrics.IncrementUserOperations("create", true)
	s.log.Debug("User created successfully",
		slog.Int64("id", createdUser.ID),
		slog.String("username", createdUser.Username))
	return createdUser, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*models.User, error) {
	s.log.Debug("Getting user by ID", slog.Int64("id", id))

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.metrics.IncrementUserOperations("get", false)
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.Int64("id", id))
			return nil, custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed to get user by id",
				slog.String("error", err.Error()),
				slog.Int64("id", id),
			)
			return nil, custom_errors.ErrDatabaseQuery
		}
	}
	s.metrics.IncrementUserOperations("get", true)
	s.log.Debug("User retrieved successfully",
		slog.Int64("id", user.ID),
		slog.String("username", user.Username))
	return user, nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	s.log.Debug("Getting user by username", slog.String("username", username))

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		s.metrics.IncrementUserOperations("get_by_username", false)
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.String("username", username))
			return nil, custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed to get user by username",
				slog.String("error", err.Error()),
				slog.String("username", username),
			)
			return nil, custom_errors.ErrDatabaseQuery
		}
	}
	s.metrics.IncrementUserOperations("get_by_username", true)
	s.log.Debug("User retrieved by username successfully",
		slog.Int64("id", user.ID),
		slog.String("username", user.Username))
	return user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	s.log.Debug("Getting user by email", slog.String("email", email))

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		s.metrics.IncrementUserOperations("get_by_email", false)
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.String("email", email))
			return nil, custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed to get user by username",
				slog.String("error", err.Error()),
				slog.String("email", email),
			)
			return nil, custom_errors.ErrDatabaseQuery
		}
	}
	s.metrics.IncrementUserOperations("get_by_email", true)
	s.log.Debug("User retrieved by email successfully",
		slog.Int64("id", user.ID),
		slog.String("email", user.Email))
	return user, nil
}

func (s *Service) Update(ctx context.Context, user *models.User) (*models.User, error) {
	s.log.Debug("Updating user",
		slog.Int64("id", user.ID),
		slog.String("username", user.Username))

	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		s.metrics.IncrementUserOperations("update", false)
		switch {
		case errors.Is(err, custom_errors.ErrUsernameExists):
			s.log.Debug("Username already exists", slog.String("username", user.Username))
			return nil, custom_errors.ErrUsernameExists
		case errors.Is(err, custom_errors.ErrEmailExists):
			s.log.Debug("Email already exists", slog.String("email", user.Email))
			return nil, custom_errors.ErrEmailExists
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.Int64("id", user.ID))
			return nil, custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed update user",
				slog.String("error", err.Error()),
				slog.Int64("id", user.ID),
			)
			return nil, custom_errors.ErrDatabaseQuery
		}
	}
	s.metrics.IncrementUserOperations("update", true)
	s.log.Debug("User updated successfully",
		slog.Int64("id", updatedUser.ID),
		slog.String("username", updatedUser.Username))
	return updatedUser, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	s.log.Debug("Deleting user", slog.Int64("id", id))

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.metrics.IncrementUserOperations("delete", false)
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.Int64("id", id))
			return custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed to delete user",
				slog.String("error", err.Error()),
				slog.Int64("id", id),
			)
			return custom_errors.ErrDatabaseQuery
		}
	}
	s.metrics.IncrementUserOperations("delete", true)
	s.log.Debug("User deleted successfully", slog.Int64("id", id))
	return nil
}

func (s *Service) Search(ctx context.Context, query string, page, limit int) ([]*models.User, int, error) {
	s.log.Debug("Searching users",
		slog.String("query", query),
		slog.Int("page", page),
		slog.Int("limit", limit))

	offset := (page - 1) * limit
	users, count, err := s.repo.Search(ctx, query, offset, limit)
	if err != nil {
		s.metrics.IncrementUserOperations("search", false)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Debug("No users found for search query",
				slog.String("query", query),
				slog.Int("page", page),
				slog.Int("limit", limit))
			return []*models.User{}, 0, nil
		default:
			s.log.Error("Failed to search users",
				slog.String("error", err.Error()),
				slog.String("query", query),
				slog.Int("page", page),
				slog.Int("limit", limit))
			return nil, 0, custom_errors.ErrDatabaseQuery
		}
	}
	s.metrics.IncrementUserOperations("search", true)
	s.log.Debug("Search completed successfully",
		slog.String("query", query),
		slog.Int("count", count))
	return users, count, nil
}

func (s *Service) UpdatePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	s.log.Debug("Updating user password", slog.Int64("id", id))

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.metrics.IncrementUserOperations("update_password", false)
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.Int64("id", id))
			return custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed to get user",
				slog.String("error", err.Error()),
				slog.Int64("id", id))
			return custom_errors.ErrDatabaseQuery
		}
	}

	err = s.repo.UpdatePassword(ctx, id, newPassword)
	if err != nil {
		s.metrics.IncrementUserOperations("update_password", false)
		s.log.Error("Failed update user",
			slog.String("error", err.Error()),
			slog.Int64("id", id))
		return custom_errors.ErrDatabaseQuery
	}
	s.metrics.IncrementUserOperations("update_password", true)
	s.log.Debug("User password updated successfully", slog.Int64("id", id))
	return nil
}

func (s *Service) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	s.log.Debug("Updating user avatar",
		slog.Int64("id", id),
		slog.String("avatarURL", avatarURL))

	err := s.repo.UpdateAvatar(ctx, id, avatarURL)
	if err != nil {
		s.metrics.IncrementUserOperations("update_avatar", false)
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.Int64("id", id))
			return custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed to update avatar user",
				slog.String("error", err.Error()),
				slog.Int64("id", id),
			)
			return custom_errors.ErrDatabaseQuery
		}
	}

	s.metrics.IncrementUserOperations("update_avatar", true)
	s.log.Debug("User avatar updated successfully", slog.Int64("id", id))
	return nil
}
