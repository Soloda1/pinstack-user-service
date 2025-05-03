package user_grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetRequest struct {
	Id int64 `validate:"required,gt=0"`
}

var getValidator = validator.New()

func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	input := GetRequest{Id: req.Id}
	if err := getValidator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.userService.Get(ctx, req.Id)
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
