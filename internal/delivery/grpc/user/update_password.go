package user_grpc

import (
	"context"

	"pinstack-user-service/internal/custom_errors"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UpdatePasswordRequest struct {
	Id          int64  `validate:"required,gt=0"`
	OldPassword string `validate:"required"`
	NewPassword string `validate:"required"`
}

func (s *UserGRPCService) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	input := UpdatePasswordRequest{
		Id:          req.Id,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	if err := validate.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.userService.UpdatePassword(ctx, req.Id, req.OldPassword, req.NewPassword); err != nil {
		switch err {
		case custom_errors.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case custom_errors.ErrInvalidPassword:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}
