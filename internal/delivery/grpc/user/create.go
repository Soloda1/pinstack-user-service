package user_grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pinstack-user-service/internal/model"
	"pinstack-user-service/internal/utils"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Request struct {
	Username  string `validate:"required,min=3"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=8"`
	FullName  string `validate:"omitempty"`
	Bio       string `validate:"omitempty"`
	AvatarURL string `validate:"omitempty,url"`
}

func (s *UserGRPCService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	input := Request{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FullName:  utils.StrPtrToStr(req.FullName),
		Bio:       utils.StrPtrToStr(req.Bio),
		AvatarURL: utils.StrPtrToStr(req.AvatarUrl),
	}

	if err := validate.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := &model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FullName:  req.FullName,
		Bio:       req.Bio,
		AvatarURL: req.AvatarUrl,
	}

	createdUser, err := s.userService.Create(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.User{
		Id:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		Bio:       createdUser.Bio,
		AvatarUrl: createdUser.AvatarURL,
		CreatedAt: timestamppb.New(createdUser.CreatedAt),
		UpdatedAt: timestamppb.New(createdUser.UpdatedAt),
	}

	return resp, nil
}
