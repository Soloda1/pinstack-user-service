package user_grpc

import (
	"context"
	"errors"
	"log/slog"
	"pinstack-user-service/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"pinstack-user-service/internal/custom_errors"
	"pinstack-user-service/internal/model"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
)

type UpdateRequest struct {
	Id       int64   `validate:"required,gt=0"`
	Username *string `validate:"omitempty,min=3"`
	Email    *string `validate:"omitempty,email"`
	FullName *string `validate:"omitempty"`
	Bio      *string `validate:"omitempty"`
}

func (s *UserGRPCService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	input := UpdateRequest{
		Id:       req.Id,
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Bio:      req.Bio,
	}

	if err := validate.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := &model.User{
		ID: req.Id,
	}

	if req.Username != nil {
		user.Username = utils.StrPtrToStr(req.Username)
	}
	if req.Email != nil {
		user.Email = utils.StrPtrToStr(req.Email)
	}
	if req.FullName != nil {
		user.FullName = req.FullName
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}

	updatedUser, err := s.userService.Update(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, custom_errors.ErrUsernameExists), errors.Is(err, custom_errors.ErrEmailExists):
			s.log.Debug("Email or Username already exists received in grpc", slog.String("error", err.Error()))
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.User{
		Id:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		FullName:  updatedUser.FullName,
		Bio:       updatedUser.Bio,
		AvatarUrl: updatedUser.AvatarURL,
		CreatedAt: timestamppb.New(updatedUser.CreatedAt),
		UpdatedAt: timestamppb.New(updatedUser.UpdatedAt),
	}, nil
}
