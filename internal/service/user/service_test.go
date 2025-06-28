package user_service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"pinstack-user-service/internal/custom_errors"
	"pinstack-user-service/internal/logger"
	"pinstack-user-service/internal/model"
	user_service "pinstack-user-service/internal/service/user"
	"pinstack-user-service/mocks"
)

func setupTest(t *testing.T) (user_service.UserService, *mocks.UserRepository, func()) {
	mockRepo := mocks.NewUserRepository(t)
	log := logger.New("test")
	service := user_service.NewUserService(mockRepo, log)
	return service, mockRepo, func() {}
}

func TestUserService_Create(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		user      *model.User
		mockSetup func()
		want      *model.User
		wantErr   error
	}{
		{
			name: "successful creation",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(
					&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
					}, nil).Once()
			},
			want: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "username already exists",
			user: &model.User{
				Username: "existinguser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(
					nil, custom_errors.ErrUsernameExists).Once()
			},
			want:    nil,
			wantErr: custom_errors.ErrUsernameExists,
		},
		{
			name: "email already exists",
			user: &model.User{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(
					nil, custom_errors.ErrEmailExists).Once()
			},
			want:    nil,
			wantErr: custom_errors.ErrEmailExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := service.Create(context.Background(), tt.user)

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
			}
		})
	}
}

func TestUserService_Get(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		id        int64
		mockSetup func()
		want      *model.User
		wantErr   error
	}{
		{
			name: "successful get",
			id:   1,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, int64(1)).Return(
					&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
					}, nil).Once()
			},
			want: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			id:   999,
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, custom_errors.ErrUserNotFound).Once()
			},
			want:    nil,
			wantErr: custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := service.Get(context.Background(), tt.id)

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
			}
		})
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		username  string
		mockSetup func()
		want      *model.User
		wantErr   error
	}{
		{
			name:     "successful get",
			username: "testuser",
			mockSetup: func() {
				mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(
					&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
					}, nil).Once()
			},
			want: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: nil,
		},
		{
			name:     "user not found",
			username: "nonexistent",
			mockSetup: func() {
				mockRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, custom_errors.ErrUserNotFound).Once()
			},
			want:    nil,
			wantErr: custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := service.GetByUsername(context.Background(), tt.username)

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
			}
		})
	}
}

func TestUserService_GetByEmail(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		email     string
		mockSetup func()
		want      *model.User
		wantErr   error
	}{
		{
			name:  "successful get",
			email: "test@example.com",
			mockSetup: func() {
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(
					&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
					}, nil).Once()
			},
			want: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			wantErr: nil,
		},
		{
			name:  "user not found",
			email: "nonexistent@example.com",
			mockSetup: func() {
				mockRepo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, custom_errors.ErrUserNotFound).Once()
			},
			want:    nil,
			wantErr: custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := service.GetByEmail(context.Background(), tt.email)

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
			}
		})
	}
}

func TestUserService_Update(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		user      *model.User
		mockSetup func()
		want      *model.User
		wantErr   error
	}{
		{
			name: "successful update",
			user: &model.User{
				ID:       1,
				Username: "updateduser",
				Email:    "updated@example.com",
			},
			mockSetup: func() {
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(
					&model.User{
						ID:       1,
						Username: "updateduser",
						Email:    "updated@example.com",
					}, nil).Once()
			},
			want: &model.User{
				ID:       1,
				Username: "updateduser",
				Email:    "updated@example.com",
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			user: &model.User{
				ID:       999,
				Username: "nonexistent",
				Email:    "nonexistent@example.com",
			},
			mockSetup: func() {
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(
					&model.User{
						ID:       999,
						Username: "nonexistent",
						Email:    "nonexistent@example.com",
					}, custom_errors.ErrUserNotFound).Once()
			},
			want:    nil,
			wantErr: custom_errors.ErrUserNotFound,
		},
		{
			name: "database error",
			user: &model.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			mockSetup: func() {
				mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(
					&model.User{
						ID:       1,
						Username: "testuser",
						Email:    "test@example.com",
					}, custom_errors.ErrDatabaseQuery).Once()
			},
			want:    nil,
			wantErr: custom_errors.ErrDatabaseQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := service.Update(context.Background(), tt.user)

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
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		id        int64
		mockSetup func()
		wantErr   error
	}{
		{
			name: "successful delete",
			id:   1,
			mockSetup: func() {
				mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			id:   999,
			mockSetup: func() {
				mockRepo.On("Delete", mock.Anything, int64(999)).Return(custom_errors.ErrUserNotFound).Once()
			},
			wantErr: custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.Delete(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_Search(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		query     string
		page      int
		limit     int
		mockSetup func()
		wantUsers []*model.User
		wantCount int
		wantErr   error
	}{
		{
			name:  "successful search",
			query: "test",
			page:  1,
			limit: 10,
			mockSetup: func() {
				mockRepo.On("Search", mock.Anything, "test", 0, 10).Return(
					[]*model.User{
						{
							ID:       1,
							Username: "testuser1",
							Email:    "test1@example.com",
						},
						{
							ID:       2,
							Username: "testuser2",
							Email:    "test2@example.com",
						},
					}, 2, nil).Once()
			},
			wantUsers: []*model.User{
				{
					ID:       1,
					Username: "testuser1",
					Email:    "test1@example.com",
				},
				{
					ID:       2,
					Username: "testuser2",
					Email:    "test2@example.com",
				},
			},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name:  "no results",
			query: "nonexistent",
			page:  1,
			limit: 10,
			mockSetup: func() {
				mockRepo.On("Search", mock.Anything, "nonexistent", 0, 10).Return(
					[]*model.User{}, 0, nil).Once()
			},
			wantUsers: []*model.User{},
			wantCount: 0,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotUsers, gotCount, err := service.Search(context.Background(), tt.query, tt.page, tt.limit)

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
					assert.Equal(t, wantUser.ID, gotUsers[i].ID)
					assert.Equal(t, wantUser.Username, gotUsers[i].Username)
					assert.Equal(t, wantUser.Email, gotUsers[i].Email)
				}
			}
		})
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name          string
		id            int64
		oldPassword   string
		newPassword   string
		mockSetup     func()
		expectedError error
	}{
		{
			name:        "successful password update",
			id:          1,
			oldPassword: "oldpass",
			newPassword: "newpass",
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, int64(1)).Return(
					&model.User{
						ID:       1,
						Password: "oldpass",
					}, nil).Once()
				mockRepo.On("UpdatePassword", mock.Anything, int64(1), "newpass").Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:        "user not found",
			id:          999,
			oldPassword: "oldpass",
			newPassword: "newpass",
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, custom_errors.ErrUserNotFound).Once()
			},
			expectedError: custom_errors.ErrUserNotFound,
		},
		{
			name:        "database error on get user",
			id:          1,
			oldPassword: "oldpass",
			newPassword: "newpass",
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, int64(1)).Return(nil, assert.AnError).Once()
			},
			expectedError: custom_errors.ErrDatabaseQuery,
		},
		{
			name:        "database error on update password",
			id:          1,
			oldPassword: "oldpass",
			newPassword: "newpass",
			mockSetup: func() {
				mockRepo.On("GetByID", mock.Anything, int64(1)).Return(
					&model.User{
						ID:       1,
						Password: "oldpass",
					}, nil).Once()
				mockRepo.On("UpdatePassword", mock.Anything, int64(1), "newpass").Return(assert.AnError).Once()
			},
			expectedError: custom_errors.ErrDatabaseQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.UpdatePassword(context.Background(), tt.id, tt.oldPassword, tt.newPassword)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_UpdateAvatar(t *testing.T) {
	service, mockRepo, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name      string
		id        int64
		avatarURL string
		mockSetup func()
		wantErr   error
	}{
		{
			name:      "successful update",
			id:        1,
			avatarURL: "https://example.com/avatar.jpg",
			mockSetup: func() {
				mockRepo.On("UpdateAvatar", mock.Anything, int64(1), "https://example.com/avatar.jpg").Return(nil).Once()
			},
			wantErr: nil,
		},
		{
			name:      "user not found",
			id:        999,
			avatarURL: "https://example.com/avatar.jpg",
			mockSetup: func() {
				mockRepo.On("UpdateAvatar", mock.Anything, int64(999), "https://example.com/avatar.jpg").Return(custom_errors.ErrUserNotFound).Once()
			},
			wantErr: custom_errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.UpdateAvatar(context.Background(), tt.id, tt.avatarURL)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
