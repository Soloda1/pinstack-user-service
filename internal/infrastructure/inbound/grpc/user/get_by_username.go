package user_grpc

import (
	"context"
	"errors"

	"github.com/soloda1/pinstack-proto-definitions/custom_errors"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GetByUsernameRequest struct {
	Username string `validate:"required,min=3"`
}

func (s *UserGRPCService) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.User, error) {
	input := GetByUsernameRequest{Username: req.Username}
	if err := validate.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, custom_errors.ErrUserNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarURL,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Password:  user.Password,
	}, nil
}
