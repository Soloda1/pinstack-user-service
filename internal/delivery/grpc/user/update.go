package user_grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

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
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FullName != nil {
		user.FullName = req.FullName
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}

	updatedUser, err := s.userService.Update(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
