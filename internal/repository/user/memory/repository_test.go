package memory_test

import (
	"context"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pinstack-user-service/internal/logger"
	"pinstack-user-service/internal/model"
	"pinstack-user-service/internal/repository/user/memory"
)

func setupTest(t *testing.T) (*memory.Repository, func()) {
	log := logger.New("test")
	repo := memory.NewUserRepository(log)
	return repo, func() {}
}

func TestRepository_Create(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password",
		FullName:  stringPtr("Test User"),
		Bio:       stringPtr("Test Bio"),
		AvatarURL: stringPtr("https://example.com/avatar.jpg"),
	}

	tests := []struct {
		name    string
		user    *model.User
		want    *model.User
		wantErr error
	}{
		{
			name: "successful creation",
			user: user,
			want: &model.User{
				Username:  user.Username,
				Email:     user.Email,
				FullName:  user.FullName,
				Bio:       user.Bio,
				AvatarURL: user.AvatarURL,
			},
			wantErr: nil,
		},
		{
			name:    "duplicate username",
			user:    user,
			want:    nil,
			wantErr: custom_errors.ErrUsernameExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Create(context.Background(), tt.user)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.NotZero(t, got.ID)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.FullName, got.FullName)
				assert.Equal(t, tt.want.Bio, got.Bio)
				assert.Equal(t, tt.want.AvatarURL, got.AvatarURL)
				assert.False(t, got.CreatedAt.IsZero())
				assert.False(t, got.UpdatedAt.IsZero())
			}
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестового пользователя
	user := &model.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password",
		FullName:  stringPtr("Test User"),
		Bio:       stringPtr("Test Bio"),
		AvatarURL: stringPtr("https://example.com/avatar.jpg"),
	}
	createdUser, err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	tests := []struct {
		name    string
		id      int64
		want    *model.User
		wantErr error
	}{
		{
			name:    "successful get",
			id:      createdUser.ID,
			want:    createdUser,
			wantErr: nil,
		},
		{
			name:    "user not found",
			id:      999,
			want:    nil,
			wantErr: custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.FullName, got.FullName)
				assert.Equal(t, tt.want.Bio, got.Bio)
				assert.Equal(t, tt.want.AvatarURL, got.AvatarURL)
			}
		})
	}
}

func TestRepository_Search(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестовых пользователей
	users := []*model.User{
		{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "password",
			FullName: stringPtr("John Doe"),
			Bio:      stringPtr("Software Developer"),
		},
		{
			Username: "jane_smith",
			Email:    "jane@example.com",
			Password: "password",
			FullName: stringPtr("Jane Smith"),
			Bio:      stringPtr("Designer"),
		},
	}

	for _, user := range users {
		_, err := repo.Create(context.Background(), user)
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		query     string
		offset    int
		pageSize  int
		wantCount int
		wantTotal int
		wantErr   error
	}{
		{
			name:      "search by username",
			query:     "john",
			offset:    0,
			pageSize:  10,
			wantCount: 1,
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name:      "search by email",
			query:     "jane@example.com",
			offset:    0,
			pageSize:  10,
			wantCount: 1,
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name:      "search by full name",
			query:     "Smith",
			offset:    0,
			pageSize:  10,
			wantCount: 1,
			wantTotal: 1,
			wantErr:   nil,
		},
		{
			name:      "search with pagination",
			query:     "",
			offset:    1,
			pageSize:  1,
			wantCount: 1,
			wantTotal: 2,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, total, err := repo.Search(context.Background(), tt.query, tt.offset, tt.pageSize)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantCount, len(got))
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}

// Вспомогательная функция для создания указателя на строку
func stringPtr(s string) *string {
	return &s
}
