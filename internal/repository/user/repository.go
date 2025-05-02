package user_repository

import (
	"context"
	"errors"
	"log/slog"
	"pinstack-user-service/internal/custom_errors"
	"pinstack-user-service/internal/logger"
	"pinstack-user-service/internal/repository"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"pinstack-user-service/internal/model"
)

type Repository struct {
	storage *repository.Storage
	log     *logger.Logger
}

func NewUserRepository(storage *repository.Storage, log *logger.Logger) *Repository {
	return &Repository{storage, log}
}

func (r *Repository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	createdAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}

	args := pgx.NamedArgs{
		"username":   user.Username,
		"password":   user.Password,
		"email":      user.Email,
		"full_name":  user.FullName,
		"bio":        user.Bio,
		"avatar_url": user.AvatarURL,
		"created_at": createdAt,
		"updated_at": createdAt,
	}

	query := `
		INSERT INTO users (username, password, email, full_name, bio, avatar_url, created_at, updated_at)
		VALUES (@username, @password, @email, @full_name, @bio, @avatar_url, @created_at, @updated_at)
		RETURNING id, username, email, full_name, bio, avatar_url, created_at, updated_at`

	var createdUser model.User
	err := r.storage.Pool.QueryRow(ctx, query, args).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.FullName,
		&createdUser.Bio,
		&createdUser.AvatarURL,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_username_key" {
				return nil, custom_errors.ErrUsernameExists
			}
			if pgErr.ConstraintName == "users_email_key" {
				return nil, custom_errors.ErrEmailExists
			}
		}
		return nil, err
	}

	return &createdUser, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := pgx.NamedArgs{"id": id}
	query := `SELECT id, username, email, full_name, avatar_url, created_at, updated_at
				FROM users WHERE id = @id`
	row := r.storage.Pool.QueryRow(ctx, query, args)
	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Bio,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug("User not found by id", slog.String("error", err.Error()))
			return nil, custom_errors.ErrUserNotFound
		}
		r.log.Error("Error getting user by id", slog.String("error", err.Error()))
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	args := pgx.NamedArgs{"username": username}
	query := `SELECT id, username, email, full_name, avatar_url, created_at, updated_at
				FROM users WHERE username = @username`
	row := r.storage.Pool.QueryRow(ctx, query, args)
	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Bio,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug("User not found by username", slog.String("error", err.Error()))
			return nil, custom_errors.ErrUserNotFound
		}
		r.log.Error("Error getting user by username", slog.String("error", err.Error()))
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := pgx.NamedArgs{"email": email}
	query := `SELECT id, username, email, full_name, avatar_url, created_at, updated_at
				FROM users WHERE email = @email`
	row := r.storage.Pool.QueryRow(ctx, query, args)
	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Bio,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug("User not found by email", slog.String("error", err.Error()))
			return nil, custom_errors.ErrUserNotFound
		}
		r.log.Error("Error getting user by email", slog.String("error", err.Error()))
		return nil, err
	}
	return user, nil
}

func (r *Repository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	updatedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}

	args := pgx.NamedArgs{
		"id":         user.ID,
		"updated_at": updatedAt,
	}

	query := `UPDATE users SET updated_at = @updated_at`

	if user.Username != "" {
		query += ", username = @username"
		args["username"] = user.Username
	}
	if user.Email != "" {
		query += ", email = @email"
		args["email"] = user.Email
	}
	if user.FullName != nil {
		query += ", full_name = @full_name"
		args["full_name"] = *user.FullName
	}
	if user.Bio != nil {
		query += ", bio = @bio"
		args["bio"] = *user.Bio
	}

	query += ` WHERE id = @id 
		RETURNING id, username, email, full_name, bio, avatar_url, created_at, updated_at`

	var updatedUser model.User
	err := r.storage.Pool.QueryRow(ctx, query, args).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.FullName,
		&updatedUser.Bio,
		&updatedUser.AvatarURL,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, custom_errors.ErrUserNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_username_key" {
				return nil, custom_errors.ErrUsernameExists
			}
			if pgErr.ConstraintName == "users_email_key" {
				return nil, custom_errors.ErrEmailExists
			}
		}
		return nil, err
	}

	return &updatedUser, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	args := pgx.NamedArgs{"id": id}
	query := `DELETE FROM users WHERE id = @id`
	result, err := r.storage.Pool.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("Error deleting user", slog.String("error", err.Error()))
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return custom_errors.ErrUserNotFound
	}
	return nil
}

func (r *Repository) Search(ctx context.Context, searchQuery string, offset, limit int) ([]*model.User, int, error) {
	args := pgx.NamedArgs{
		"query":  searchQuery,
		"offset": offset,
		"limit":  limit,
	}

	query := `
			SELECT id, username, email, full_name, bio, avatar_url, created_at, updated_at FROM users 
			WHERE 
				username ILIKE '%' || @query || '%' OR
				email ILIKE '%' || @query || '%' OR
				full_name ILIKE '%' || @query || '%'
			LIMIT @limit OFFSET @offset
			`

	rows, err := r.storage.Pool.Query(ctx, query, args)
	if err != nil {
		r.log.Error("Error searching users", slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FullName,
			&user.Bio,
			&user.AvatarURL,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			r.log.Error("Error getting user", slog.String("error", err.Error()))
			return nil, 0, err
		}
		users = append(users, &user)
	}

	return users, len(users), nil
}

func (r *Repository) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	updatedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	args := pgx.NamedArgs{
		"id":         id,
		"password":   hashedPassword,
		"updated_at": updatedAt,
	}

	query := `
		UPDATE users 
		SET password = @password,
			updated_at = @updated_at
		WHERE id = @id
		RETURNING id`

	var userID int64
	err := r.storage.Pool.QueryRow(ctx, query, args).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return custom_errors.ErrUserNotFound
		}
		r.log.Error("Error updating password", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *Repository) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	updatedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}

	args := pgx.NamedArgs{
		"id":         id,
		"avatar_url": avatarURL,
		"updated_at": updatedAt,
	}

	query := `
		UPDATE users 
		SET avatar_url = @avatar_url,
			updated_at = @updated_at
		WHERE id = @id
		RETURNING id`

	var userID int64
	err := r.storage.Pool.QueryRow(ctx, query, args).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return custom_errors.ErrUserNotFound
		}
		r.log.Error("Error updating avatar", slog.String("error", err.Error()))
		return err
	}

	return nil
}
