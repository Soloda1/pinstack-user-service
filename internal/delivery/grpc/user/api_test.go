package user_grpc_test

import (
	"context"
	"github.com/soloda1/pinstack-proto-definitions/custom_errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	user_grpc "pinstack-user-service/internal/delivery/grpc/user"
	"pinstack-user-service/internal/logger"
	"pinstack-user-service/internal/model"
	"pinstack-user-service/mocks"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
)

func setupTest(t *testing.T) (*user_grpc.UserGRPCService, *mocks.UserService, func()) {
	log := logger.New("test")
	mockService := mocks.NewUserService(t)
	handler := user_grpc.NewUserGRPCService(mockService, log)
	return handler, mockService, func() {}
}

func TestUserGRPCService_CreateUser(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *pb.CreateUserRequest
		mock    func()
		want    *pb.User
		wantErr error
	}{
		{
			name: "successful creation",
			req: &pb.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
				FullName: stringPtr("Test User"),
			},
			mock: func() {
				mockService.EXPECT().Create(
					context.Background(),
					&model.User{
						Username: "testuser",
						Email:    "test@example.com",
						Password: "password",
						FullName: stringPtr("Test User"),
					},
				).Return(&model.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password",
					FullName: stringPtr("Test User"),
				}, nil)
			},
			want: &pb.User{
				Id:       1,
				Username: "testuser",
				Email:    "test@example.com",
				FullName: stringPtr("Test User"),
			},
			wantErr: nil,
		},
		{
			name: "duplicate username",
			req: &pb.CreateUserRequest{
				Username: "testuser",
				Email:    "another@example.com",
				Password: "password",
				FullName: stringPtr("Test User"),
			},
			mock: func() {
				mockService.EXPECT().Create(
					context.Background(),
					&model.User{
						Username: "testuser",
						Email:    "another@example.com",
						Password: "password",
						FullName: stringPtr("Test User"),
					},
				).Return(nil, custom_errors.ErrUsernameExists)
			},
			want:    nil,
			wantErr: status.Error(codes.AlreadyExists, custom_errors.ErrUsernameExists.Error()),
		},
		{
			name: "duplicate email",
			req: &pb.CreateUserRequest{
				Username: "anotheruser",
				Email:    "test@example.com",
				Password: "password",
				FullName: stringPtr("Test User"),
			},
			mock: func() {
				mockService.EXPECT().Create(
					context.Background(),
					&model.User{
						Username: "anotheruser",
						Email:    "test@example.com",
						Password: "password",
						FullName: stringPtr("Test User"),
					},
				).Return(nil, custom_errors.ErrEmailExists)
			},
			want:    nil,
			wantErr: status.Error(codes.AlreadyExists, custom_errors.ErrEmailExists.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := handler.CreateUser(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.FullName, got.FullName)
				assert.NotZero(t, got.Id)
			}
		})
	}
}

func TestUserGRPCService_GetUser(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *pb.GetUserRequest
		mock    func()
		want    *pb.User
		wantErr error
	}{
		{
			name: "successful get",
			req: &pb.GetUserRequest{
				Id: 1,
			},
			mock: func() {
				mockService.EXPECT().Get(
					context.Background(),
					int64(1),
				).Return(&model.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password",
					FullName: stringPtr("Test User"),
				}, nil)
			},
			want: &pb.User{
				Id:       1,
				Username: "testuser",
				Email:    "test@example.com",
				FullName: stringPtr("Test User"),
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: &pb.GetUserRequest{
				Id: 999,
			},
			mock: func() {
				mockService.EXPECT().Get(
					context.Background(),
					int64(999),
				).Return(nil, custom_errors.ErrUserNotFound)
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, custom_errors.ErrUserNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := handler.GetUser(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.FullName, got.FullName)
			}
		})
	}
}

func TestUserGRPCService_GetUserByUsername(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *pb.GetUserByUsernameRequest
		mock    func()
		want    *pb.User
		wantErr error
	}{
		{
			name: "successful get",
			req: &pb.GetUserByUsernameRequest{
				Username: "testuser",
			},
			mock: func() {
				mockService.EXPECT().GetByUsername(
					context.Background(),
					"testuser",
				).Return(&model.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password",
					FullName: stringPtr("Test User"),
				}, nil)
			},
			want: &pb.User{
				Id:       1,
				Username: "testuser",
				Email:    "test@example.com",
				FullName: stringPtr("Test User"),
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: &pb.GetUserByUsernameRequest{
				Username: "nonexistent",
			},
			mock: func() {
				mockService.EXPECT().GetByUsername(
					context.Background(),
					"nonexistent",
				).Return(nil, custom_errors.ErrUserNotFound)
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, custom_errors.ErrUserNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := handler.GetUserByUsername(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.FullName, got.FullName)
			}
		})
	}
}

func TestUserGRPCService_GetUserByEmail(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *pb.GetUserByEmailRequest
		mock    func()
		want    *pb.User
		wantErr error
	}{
		{
			name: "successful get",
			req: &pb.GetUserByEmailRequest{
				Email: "test@example.com",
			},
			mock: func() {
				mockService.EXPECT().GetByEmail(
					context.Background(),
					"test@example.com",
				).Return(&model.User{
					ID:       1,
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password",
					FullName: stringPtr("Test User"),
				}, nil)
			},
			want: &pb.User{
				Id:       1,
				Username: "testuser",
				Email:    "test@example.com",
				FullName: stringPtr("Test User"),
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: &pb.GetUserByEmailRequest{
				Email: "nonexistent@example.com",
			},
			mock: func() {
				mockService.EXPECT().GetByEmail(
					context.Background(),
					"nonexistent@example.com",
				).Return(nil, custom_errors.ErrUserNotFound)
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, custom_errors.ErrUserNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := handler.GetUserByEmail(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.FullName, got.FullName)
			}
		})
	}
}

func TestUserGRPCService_UpdateUser(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *pb.UpdateUserRequest
		mock    func()
		want    *pb.User
		wantErr error
	}{
		{
			name: "successful update",
			req: &pb.UpdateUserRequest{
				Id:       1,
				Username: strPtr("updateduser"),
				Email:    strPtr("updated@example.com"),
				FullName: strPtr("Updated User"),
			},
			mock: func() {
				mockService.EXPECT().Update(
					context.Background(),
					&model.User{
						ID:       1,
						Username: "updateduser",
						Email:    "updated@example.com",
						FullName: stringPtr("Updated User"),
					},
				).Return(&model.User{
					ID:       1,
					Username: "updateduser",
					Email:    "updated@example.com",
					FullName: stringPtr("Updated User"),
				}, nil)
			},
			want: &pb.User{
				Id:       1,
				Username: "updateduser",
				Email:    "updated@example.com",
				FullName: stringPtr("Updated User"),
			},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: &pb.UpdateUserRequest{
				Id:       999,
				Username: strPtr("nonexistent"),
				Email:    strPtr("nonexistent@example.com"),
				FullName: strPtr("Nonexistent User"),
			},
			mock: func() {
				mockService.EXPECT().Update(
					context.Background(),
					&model.User{
						ID:       999,
						Username: "nonexistent",
						Email:    "nonexistent@example.com",
						FullName: stringPtr("Nonexistent User"),
					},
				).Return(nil, custom_errors.ErrUserNotFound)
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, custom_errors.ErrUserNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := handler.UpdateUser(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Id, got.Id)
				assert.Equal(t, tt.want.Username, got.Username)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.FullName, got.FullName)
			}
		})
	}
}

func TestUserGRPCService_DeleteUser(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *pb.DeleteUserRequest
		mock    func()
		want    *emptypb.Empty
		wantErr error
	}{
		{
			name: "successful delete",
			req: &pb.DeleteUserRequest{
				Id: 1,
			},
			mock: func() {
				mockService.EXPECT().Delete(
					context.Background(),
					int64(1),
				).Return(nil)
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: &pb.DeleteUserRequest{
				Id: 999,
			},
			mock: func() {
				mockService.EXPECT().Delete(
					context.Background(),
					int64(999),
				).Return(custom_errors.ErrUserNotFound)
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, custom_errors.ErrUserNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := handler.DeleteUser(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestUserGRPCService_SearchUsers(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *pb.SearchUsersRequest
		mock    func()
		want    *pb.SearchUsersResponse
		wantErr error
	}{
		{
			name: "search by username prefix",
			req: &pb.SearchUsersRequest{
				Query:  "test",
				Offset: 0,
				Limit:  10,
			},
			mock: func() {
				mockService.EXPECT().Search(
					context.Background(),
					"test",
					0,
					10,
				).Return([]*model.User{
					{
						ID:        1,
						Username:  "testuser1",
						Email:     "test1@example.com",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						FullName:  stringPtr("Test User 1"),
					},
					{
						ID:        2,
						Username:  "testuser2",
						Email:     "test2@example.com",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						FullName:  stringPtr("Test User 2"),
					},
				}, 2, nil)
			},
			want: &pb.SearchUsersResponse{
				Users: []*pb.User{
					{
						Id:       1,
						Username: "testuser1",
						Email:    "test1@example.com",
						FullName: stringPtr("Test User 1"),
					},
					{
						Id:       2,
						Username: "testuser2",
						Email:    "test2@example.com",
						FullName: stringPtr("Test User 2"),
					},
				},
				Total: 2,
			},
			wantErr: nil,
		},
		{
			name: "no results",
			req: &pb.SearchUsersRequest{
				Query:  "nonexistent",
				Offset: 0,
				Limit:  10,
			},
			mock: func() {
				mockService.EXPECT().Search(
					context.Background(),
					"nonexistent",
					0,
					10,
				).Return([]*model.User{}, 0, nil)
			},
			want: &pb.SearchUsersResponse{
				Users: []*pb.User{},
				Total: 0,
			},
			wantErr: nil,
		},
		{
			name: "invalid page size",
			req: &pb.SearchUsersRequest{
				Query:  "test",
				Offset: 0,
				Limit:  101, // exceeds max limit of 100
			},
			mock:    func() {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "validation error"),
		},
		{
			name: "negative offset",
			req: &pb.SearchUsersRequest{
				Query:  "test",
				Offset: -1,
				Limit:  10,
			},
			mock:    func() {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "validation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := handler.SearchUsers(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.InvalidArgument, st.Code())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.Total, got.Total)
				assert.Equal(t, len(tt.want.Users), len(got.Users))
				for i, wantUser := range tt.want.Users {
					assert.Equal(t, wantUser.Id, got.Users[i].Id)
					assert.Equal(t, wantUser.Username, got.Users[i].Username)
					assert.Equal(t, wantUser.Email, got.Users[i].Email)
					assert.Equal(t, wantUser.FullName, got.Users[i].FullName)
				}
			}
		})
	}
}

func TestUserGRPCService_UpdatePassword(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name          string
		req           *pb.UpdatePasswordRequest
		mockSetup     func(mockService *mocks.UserService)
		expectedError error
	}{
		{
			name: "successful password update",
			req: &pb.UpdatePasswordRequest{
				Id:          1,
				OldPassword: "oldpass",
				NewPassword: "newpass",
			},
			mockSetup: func(mockService *mocks.UserService) {
				mockService.EXPECT().
					UpdatePassword(mock.Anything, int64(1), "oldpass", "newpass").
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			req: &pb.UpdatePasswordRequest{
				Id:          999,
				OldPassword: "oldpass",
				NewPassword: "newpass",
			},
			mockSetup: func(mockService *mocks.UserService) {
				mockService.EXPECT().
					UpdatePassword(mock.Anything, int64(999), "oldpass", "newpass").
					Return(custom_errors.ErrUserNotFound)
			},
			expectedError: status.Error(codes.NotFound, custom_errors.ErrUserNotFound.Error()),
		},
		{
			name: "invalid old password",
			req: &pb.UpdatePasswordRequest{
				Id:          1,
				OldPassword: "wrongpass",
				NewPassword: "newpass",
			},
			mockSetup: func(mockService *mocks.UserService) {
				mockService.EXPECT().
					UpdatePassword(mock.Anything, int64(1), "wrongpass", "newpass").
					Return(custom_errors.ErrInvalidPassword)
			},
			expectedError: status.Error(codes.InvalidArgument, custom_errors.ErrInvalidPassword.Error()),
		},
		{
			name: "invalid request - missing old password",
			req: &pb.UpdatePasswordRequest{
				Id:          1,
				NewPassword: "newpass",
			},
			mockSetup:     func(mockService *mocks.UserService) {},
			expectedError: status.Error(codes.InvalidArgument, "Key: 'UpdatePasswordRequest.OldPassword' Error:Field validation for 'OldPassword' failed on the 'required' tag"),
		},
		{
			name: "invalid request - missing new password",
			req: &pb.UpdatePasswordRequest{
				Id:          1,
				OldPassword: "oldpass",
			},
			mockSetup:     func(mockService *mocks.UserService) {},
			expectedError: status.Error(codes.InvalidArgument, "Key: 'UpdatePasswordRequest.NewPassword' Error:Field validation for 'NewPassword' failed on the 'required' tag"),
		},
		{
			name: "invalid request - invalid id",
			req: &pb.UpdatePasswordRequest{
				Id:          0,
				OldPassword: "oldpass",
				NewPassword: "newpass",
			},
			mockSetup:     func(mockService *mocks.UserService) {},
			expectedError: status.Error(codes.InvalidArgument, "Key: 'UpdatePasswordRequest.Id' Error:Field validation for 'Id' failed on the 'required' tag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mockService)

			resp, err := handler.UpdatePassword(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, &emptypb.Empty{}, resp)
			}
		})
	}
}

func TestUserGRPCService_UpdateAvatar(t *testing.T) {
	handler, mockService, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		req     *pb.UpdateAvatarRequest
		mock    func()
		want    *emptypb.Empty
		wantErr error
	}{
		{
			name: "successful update",
			req: &pb.UpdateAvatarRequest{
				Id:        1,
				AvatarUrl: "https://example.com/avatar.jpg",
			},
			mock: func() {
				mockService.EXPECT().UpdateAvatar(
					context.Background(),
					int64(1),
					"https://example.com/avatar.jpg",
				).Return(nil)
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "user not found",
			req: &pb.UpdateAvatarRequest{
				Id:        999,
				AvatarUrl: "https://example.com/avatar.jpg",
			},
			mock: func() {
				mockService.EXPECT().UpdateAvatar(
					context.Background(),
					int64(999),
					"https://example.com/avatar.jpg",
				).Return(custom_errors.ErrUserNotFound)
			},
			want:    nil,
			wantErr: status.Error(codes.NotFound, custom_errors.ErrUserNotFound.Error()),
		},
		{
			name: "invalid avatar URL",
			req: &pb.UpdateAvatarRequest{
				Id:        1,
				AvatarUrl: "not-a-url",
			},
			mock:    func() {},
			want:    nil,
			wantErr: status.Error(codes.InvalidArgument, "validation error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := handler.UpdateAvatar(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				if st.Code() == codes.InvalidArgument {
					assert.Equal(t, codes.InvalidArgument, st.Code())
				} else {
					assert.Equal(t, tt.wantErr.Error(), err.Error())
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}

func stringPtr(s string) *string {
	return &s
}
