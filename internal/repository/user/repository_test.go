package user_repository_test

import (
	"context"
	"pinstack-user-service/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pinstack-user-service/internal/custom_errors"
	"pinstack-user-service/internal/model"
	user_repository "pinstack-user-service/internal/repository/user"
	"pinstack-user-service/internal/repository/user/memory"
)

func setupTest(t *testing.T) (user_repository.UserRepository, func()) {
	log := logger.New("test")
	repo := memory.NewUserRepository(log)
	return repo, func() {}
}

func TestUserRepository_Create(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		user    *model.User
		want    *model.User
		wantErr error
	}{
		{
			name: "successful creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			want: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: nil,
		},
		{
			name: "duplicate username",
			user: &model.User{
				Username: "testuser",
				Email:    "another@example.com",
				Password: "password123",
			},
			want:    nil,
			wantErr: custom_errors.ErrUsernameExists,
		},
		{
			name: "duplicate email",
			user: &model.User{
				Username: "anotheruser",
				Email:    "test@example.com",
				Password: "password123",
			},
			want:    nil,
			wantErr: custom_errors.ErrEmailExists,
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
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Password, got.Password)
				assert.NotZero(t, got.ID)
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестового пользователя
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, created)

	tests := []struct {
		name    string
		id      int64
		want    *model.User
		wantErr error
	}{
		{
			name:    "successful get",
			id:      created.ID,
			want:    created,
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
				assert.Equal(t, tt.want.Password, got.Password)
			}
		})
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестового пользователя
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, created)

	tests := []struct {
		name     string
		username string
		want     *model.User
		wantErr  error
	}{
		{
			name:     "successful get",
			username: "testuser",
			want:     created,
			wantErr:  nil,
		},
		{
			name:     "user not found",
			username: "nonexistent",
			want:     nil,
			wantErr:  custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByUsername(context.Background(), tt.username)

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
				assert.Equal(t, tt.want.Password, got.Password)
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестового пользователя
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, created)

	tests := []struct {
		name    string
		email   string
		want    *model.User
		wantErr error
	}{
		{
			name:    "successful get",
			email:   "test@example.com",
			want:    created,
			wantErr: nil,
		},
		{
			name:    "user not found",
			email:   "nonexistent@example.com",
			want:    nil,
			wantErr: custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByEmail(context.Background(), tt.email)

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
				assert.Equal(t, tt.want.Password, got.Password)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем первого тестового пользователя
	user1 := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	created1, err := repo.Create(context.Background(), user1)
	require.NoError(t, err)
	require.NotNil(t, created1)

	// Создаем второго тестового пользователя
	user2 := &model.User{
		Username: "anotheruser",
		Email:    "another@example.com",
		Password: "password123",
	}
	created2, err := repo.Create(context.Background(), user2)
	require.NoError(t, err)
	require.NotNil(t, created2)

	tests := []struct {
		name    string
		user    *model.User
		want    *model.User
		wantErr error
	}{
		{
			name: "successful update",
			user: &model.User{
				ID:       created1.ID,
				Username: "updateduser",
				Email:    "updated@example.com",
				Password: "newpassword",
			},
			want: &model.User{
				ID:       created1.ID,
				Username: "updateduser",
				Email:    "updated@example.com",
				Password: "newpassword",
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			user: &model.User{
				ID:       999,
				Username: "nonexistent",
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			want:    nil,
			wantErr: custom_errors.ErrUserNotFound,
		},
		{
			name: "duplicate username",
			user: &model.User{
				ID:       created1.ID,
				Username: created2.Username, // Используем username второго пользователя
				Email:    "updated@example.com",
				Password: "newpassword",
			},
			want:    nil,
			wantErr: custom_errors.ErrUsernameExists,
		},
		{
			name: "duplicate email",
			user: &model.User{
				ID:       created1.ID,
				Username: "updateduser",
				Email:    created2.Email, // Используем email второго пользователя
				Password: "newpassword",
			},
			want:    nil,
			wantErr: custom_errors.ErrEmailExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Update(context.Background(), tt.user)

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
				assert.Equal(t, tt.want.Password, got.Password)
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестового пользователя
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, created)

	tests := []struct {
		name    string
		id      int64
		wantErr error
	}{
		{
			name:    "successful delete",
			id:      created.ID,
			wantErr: nil,
		},
		{
			name:    "user not found",
			id:      999,
			wantErr: custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				// Проверяем, что пользователь действительно удален
				_, err := repo.GetByID(context.Background(), tt.id)
				assert.Error(t, err)
				assert.Equal(t, custom_errors.ErrUserNotFound, err)
			}
		})
	}
}

func TestUserRepository_Search(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестовых пользователей
	users := []*model.User{
		{
			Username: "testuser1",
			Email:    "test1@example.com",
			Password: "password123",
		},
		{
			Username: "testuser2",
			Email:    "test2@example.com",
			Password: "password123",
		},
		{
			Username: "anotheruser",
			Email:    "another@example.com",
			Password: "password123",
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
		wantUsers []*model.User
		wantCount int
		wantErr   error
	}{
		{
			name:     "search by username prefix",
			query:    "test",
			offset:   0,
			pageSize: 10,
			wantUsers: []*model.User{
				{
					Username: "testuser1",
					Email:    "test1@example.com",
				},
				{
					Username: "testuser2",
					Email:    "test2@example.com",
				},
			},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name:     "search by email",
			query:    "another@example.com",
			offset:   0,
			pageSize: 10,
			wantUsers: []*model.User{
				{
					Username: "anotheruser",
					Email:    "another@example.com",
				},
			},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name:      "no results",
			query:     "nonexistent",
			offset:    0,
			pageSize:  10,
			wantUsers: []*model.User{},
			wantCount: 0,
			wantErr:   nil,
		},
		{
			name:     "pagination",
			query:    "test",
			offset:   1,
			pageSize: 1,
			wantUsers: []*model.User{
				{
					Username: "testuser2",
					Email:    "test2@example.com",
				},
			},
			wantCount: 2,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUsers, gotCount, err := repo.Search(context.Background(), tt.query, tt.offset, tt.pageSize)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, gotUsers)
				assert.Equal(t, 0, gotCount)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotUsers)
				assert.Equal(t, tt.wantCount, gotCount)
				assert.Equal(t, len(tt.wantUsers), len(gotUsers))
				for i, wantUser := range tt.wantUsers {
					assert.Equal(t, wantUser.Username, gotUsers[i].Username)
					assert.Equal(t, wantUser.Email, gotUsers[i].Email)
				}
			}
		})
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестового пользователя
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "oldpassword",
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, created)

	tests := []struct {
		name     string
		id       int64
		password string
		wantErr  error
	}{
		{
			name:     "successful update",
			id:       created.ID,
			password: "newpassword",
			wantErr:  nil,
		},
		{
			name:     "user not found",
			id:       999,
			password: "newpassword",
			wantErr:  custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdatePassword(context.Background(), tt.id, tt.password)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				// Проверяем, что пароль действительно обновлен
				got, err := repo.GetByID(context.Background(), tt.id)
				assert.NoError(t, err)
				assert.Equal(t, tt.password, got.Password)
			}
		})
	}
}

func TestUserRepository_UpdateAvatar(t *testing.T) {
	repo, cleanup := setupTest(t)
	defer cleanup()

	// Создаем тестового пользователя
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, created)

	tests := []struct {
		name      string
		id        int64
		avatarURL string
		wantErr   error
	}{
		{
			name:      "successful update",
			id:        created.ID,
			avatarURL: "https://example.com/avatar.jpg",
			wantErr:   nil,
		},
		{
			name:      "user not found",
			id:        999,
			avatarURL: "https://example.com/avatar.jpg",
			wantErr:   custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateAvatar(context.Background(), tt.id, tt.avatarURL)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				// Проверяем, что аватар действительно обновлен
				got, err := repo.GetByID(context.Background(), tt.id)
				assert.NoError(t, err)
				assert.Equal(t, tt.avatarURL, *got.AvatarURL)
			}
		})
	}
}
