package user_service

import (
	"context"
	"log/slog"
	"pinstack-user-service/internal/custom_errors"
	"pinstack-user-service/internal/logger"
	"pinstack-user-service/internal/model"
	user_repository "pinstack-user-service/internal/repository/user"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type Service struct {
	repo user_repository.UserRepository
	log  *logger.Logger
}

func NewUserService(repo user_repository.UserRepository, log *logger.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) Create(ctx context.Context, user *model.User) (*model.User, error) {
	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
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

	return createdUser, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
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
	user.Password = ""
	return user, nil
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
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
	user.Password = ""
	return user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
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
	user.Password = ""
	return user, nil
}

func (s *Service) Update(ctx context.Context, user *model.User) (*model.User, error) {
	user, err := s.repo.Update(ctx, user)
	if err != nil {
		switch {
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
	return user, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
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
	return nil
}

func (s *Service) Search(ctx context.Context, query string, page, limit int) ([]*model.User, int, error) {
	offset := (page - 1) * limit
	users, count, err := s.repo.Search(ctx, query, offset, limit)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Debug("No users found for search query",
				slog.String("query", query),
				slog.Int("page", page),
				slog.Int("limit", limit))
			return []*model.User{}, 0, nil
		default:
			s.log.Error("Failed to search users",
				slog.String("error", err.Error()),
				slog.String("query", query),
				slog.Int("page", page),
				slog.Int("limit", limit))
			return nil, 0, custom_errors.ErrDatabaseQuery
		}
	}
	return users, count, nil
}

func (s *Service) UpdatePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.Int64("id", id))
			return custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed to get user for password update",
				slog.String("error", err.Error()),
				slog.Int64("id", id))
			return custom_errors.ErrDatabaseQuery
		}
	}

	if user.Password != oldPassword {
		s.log.Debug("Invalid old password", slog.Int64("id", id))
		return custom_errors.ErrInvalidPassword
	}

	err = s.repo.UpdatePassword(ctx, id, newPassword)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			s.log.Debug("User not found", slog.Int64("id", id))
			return custom_errors.ErrUserNotFound
		default:
			s.log.Error("Failed to update password user",
				slog.String("error", err.Error()),
				slog.Int64("id", id))
			return custom_errors.ErrDatabaseQuery
		}
	}
	return nil
}

func (s *Service) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	err := s.repo.UpdateAvatar(ctx, id, avatarURL)
	if err != nil {
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

	return nil
}
