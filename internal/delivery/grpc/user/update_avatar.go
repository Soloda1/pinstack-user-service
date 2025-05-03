package user_grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UpdateAvatarRequest struct {
	Id        int64  `validate:"required,gt=0"`
	AvatarUrl string `validate:"required,url"`
}

var updateAvatarValidator = validator.New()

func (s *UserGRPCService) UpdateAvatar(ctx context.Context, req *pb.UpdateAvatarRequest) (*emptypb.Empty, error) {
	input := UpdateAvatarRequest{
		Id:        req.Id,
		AvatarUrl: req.AvatarUrl,
	}
	if err := updateAvatarValidator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.userService.UpdateAvatar(ctx, req.Id, req.AvatarUrl); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
