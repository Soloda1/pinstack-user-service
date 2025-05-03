package user_grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DeleteRequest struct {
	Id int64 `validate:"required,gt=0"`
}

var deleteValidator = validator.New()

func (s *UserGRPCService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	input := DeleteRequest{Id: req.Id}
	if err := deleteValidator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.userService.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
