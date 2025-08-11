package postgres

import (
	"context"
	"errors"
	"log/slog"
	ports "pinstack-user-service/internal/domain/ports/output"
	"pinstack-user-service/internal/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"pinstack-user-service/internal/domain/models"

	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
)

type Repository struct {
	pool *pgxpool.Pool
	log  ports.Logger
}

func NewUserRepository(pool *pgxpool.Pool, log ports.Logger) *Repository {
	return &Repository{pool, log}
}

func (r *Repository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	r.log.Debug("Creating user in database",
		slog.String("username", user.Username),
		slog.String("email", user.Email))

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

	var createdUser models.User
	err := r.pool.QueryRow(ctx, query, args).Scan(
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
				r.log.Debug("Username constraint violation",
					slog.String("username", user.Username),
					slog.String("error", err.Error()))
				return nil, custom_errors.ErrUsernameExists
			}
			if pgErr.ConstraintName == "users_email_key" {
				r.log.Debug("Email constraint violation",
					slog.String("email", user.Email),
					slog.String("error", err.Error()))
				return nil, custom_errors.ErrEmailExists
			}
		}
		r.log.Error("Error creating user in database",
			slog.String("error", err.Error()),
			slog.String("username", user.Username))
		return nil, err
	}

	r.log.Debug("User created successfully in database",
		slog.Int64("id", createdUser.ID),
		slog.String("username", createdUser.Username))
	return &createdUser, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	r.log.Debug("Getting user by ID from database", slog.Int64("id", id))

	args := pgx.NamedArgs{"id": id}
	query := `SELECT id, username, password, email, full_name, bio, avatar_url, created_at, updated_at
                FROM users WHERE id = @id`
	row := r.pool.QueryRow(ctx, query, args)
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.FullName,
		&user.Bio,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug("User not found by id",
				slog.Int64("id", id),
				slog.String("error", err.Error()))
			return nil, custom_errors.ErrUserNotFound
		}
		r.log.Error("Error getting user by id", slog.String("error", err.Error()))
		return nil, err
	}
	r.log.Debug("User retrieved by ID successfully from database",
		slog.Int64("id", user.ID),
		slog.String("username", user.Username))
	return user, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.log.Debug("Getting user by username from database", slog.String("username", username))

	args := pgx.NamedArgs{"username": username}
	query := `SELECT id, username, password, email, full_name, bio, avatar_url, created_at, updated_at
                FROM users WHERE username = @username`
	row := r.pool.QueryRow(ctx, query, args)
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.FullName,
		&user.Bio,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug("User not found by username",
				slog.String("username", username),
				slog.String("error", err.Error()))
			return nil, custom_errors.ErrUserNotFound
		}
		r.log.Error("Error getting user by username", slog.String("error", err.Error()))
		return nil, err
	}
	r.log.Debug("User retrieved by username successfully from database",
		slog.Int64("id", user.ID),
		slog.String("username", user.Username))
	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.log.Debug("Getting user by email from database", slog.String("email", email))

	args := pgx.NamedArgs{"email": email}
	query := `SELECT id, username, password, email, full_name, bio, avatar_url, created_at, updated_at
                FROM users WHERE email = @email`
	row := r.pool.QueryRow(ctx, query, args)
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.FullName,
		&user.Bio,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug("User not found by email",
				slog.String("email", email),
				slog.String("error", err.Error()))
			return nil, custom_errors.ErrUserNotFound
		}
		r.log.Error("Error getting user by email", slog.String("error", err.Error()))
		return nil, err
	}
	r.log.Debug("User retrieved by email successfully from database",
		slog.Int64("id", user.ID),
		slog.String("email", user.Email))
	return user, nil
}

func (r *Repository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	r.log.Debug("Updating user in database",
		slog.Int64("id", user.ID),
		slog.String("username", user.Username))

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
		args["full_name"] = utils.StrPtrToStr(user.FullName)
	}
	if user.Bio != nil {
		query += ", bio = @bio"
		args["bio"] = utils.StrPtrToStr(user.Bio)
	}

	query += ` WHERE id = @id 
        RETURNING id, username, email, full_name, bio, avatar_url, created_at, updated_at`

	var updatedUser models.User
	err := r.pool.QueryRow(ctx, query, args).Scan(
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
			r.log.Debug("User not found for update",
				slog.Int64("id", user.ID),
				slog.String("error", err.Error()))
			return nil, custom_errors.ErrUserNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_username_key" {
				r.log.Debug("Username constraint violation during update",
					slog.String("username", user.Username),
					slog.String("error", err.Error()))
				return nil, custom_errors.ErrUsernameExists
			}
			if pgErr.ConstraintName == "users_email_key" {
				r.log.Debug("Email constraint violation during update",
					slog.String("email", user.Email),
					slog.String("error", err.Error()))
				return nil, custom_errors.ErrEmailExists
			}
		}
		r.log.Debug("Database error updating user",
			slog.Int64("id", user.ID),
			slog.String("error", err.Error()))
		return nil, err
	}

	r.log.Debug("User updated successfully in database",
		slog.Int64("id", updatedUser.ID),
		slog.String("username", updatedUser.Username))
	return &updatedUser, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	r.log.Debug("Deleting user from database", slog.Int64("id", id))

	args := pgx.NamedArgs{"id": id}
	query := `DELETE FROM users WHERE id = @id`
	result, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		r.log.Error("Error deleting user", slog.String("error", err.Error()))
		return err
	}
	if result.RowsAffected() == 0 {
		return custom_errors.ErrUserNotFound
	}
	r.log.Debug("User deleted successfully from database", slog.Int64("id", id))
	return nil
}

func (r *Repository) Search(ctx context.Context, searchQuery string, offset, limit int) ([]*models.User, int, error) {
	r.log.Debug("Searching users in database",
		slog.String("query", searchQuery),
		slog.Int("offset", offset),
		slog.Int("limit", limit))

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

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		r.log.Error("Error searching users", slog.String("error", err.Error()))
		return nil, 0, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
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

	r.log.Debug("Search completed successfully in database",
		slog.String("query", searchQuery),
		slog.Int("count", len(users)))
	return users, len(users), nil
}

func (r *Repository) UpdatePassword(ctx context.Context, id int64, newPassword string) error {
	r.log.Debug("Updating user password in database", slog.Int64("id", id))

	updatedAt := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	args := pgx.NamedArgs{
		"id":         id,
		"password":   newPassword,
		"updated_at": updatedAt,
	}

	query := `
        UPDATE users 
        SET password = @password,
            updated_at = @updated_at
        WHERE id = @id
        RETURNING id`

	var userID int64
	err := r.pool.QueryRow(ctx, query, args).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug("User not found for password update",
				slog.Int64("id", id),
				slog.String("error", err.Error()))
			return custom_errors.ErrUserNotFound
		}
		r.log.Error("Error updating password", slog.String("error", err.Error()))
		return err
	}

	r.log.Debug("User password updated successfully in database", slog.Int64("id", id))
	return nil
}

func (r *Repository) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	r.log.Debug("Updating user avatar in database",
		slog.Int64("id", id),
		slog.String("avatarURL", avatarURL))

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
	err := r.pool.QueryRow(ctx, query, args).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug("User not found for avatar update",
				slog.Int64("id", id),
				slog.String("error", err.Error()))
			return custom_errors.ErrUserNotFound
		}
		r.log.Error("Error updating avatar", slog.String("error", err.Error()))
		return err
	}

	r.log.Debug("User avatar updated successfully in database", slog.Int64("id", id))
	return nil
}
