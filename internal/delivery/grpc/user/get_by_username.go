package user_grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetByUsernameRequest struct {
	Username string `validate:"required,min=3"`
}

var getByUsernameValidator = validator.New()

func (s *UserGRPCService) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.User, error) {
	input := GetByUsernameRequest{Username: req.Username}
	if err := getByUsernameValidator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.userService.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Bio:       user.Bio,
		AvatarUrl: user.AvatarURL,
	}, nil
}
