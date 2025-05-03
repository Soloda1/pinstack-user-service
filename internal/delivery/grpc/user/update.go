package user_grpc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"pinstack-user-service/internal/model"
	"pinstack-user-service/internal/utils"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
)

type UpdateRequest struct {
	Id       int64  `validate:"required,gt=0"`
	Username string `validate:"required,min=3"`
	Email    string `validate:"required,email"`
	FullName string `validate:"omitempty"`
	Bio      string `validate:"omitempty"`
}

type UpdateResponse struct {
	User *pb.User
}

var updateValidator = validator.New()

func (s *UserGRPCService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	input := UpdateRequest{
		Id:       req.Id,
		Username: utils.StrPtrToStr(req.Username),
		Email:    utils.StrPtrToStr(req.Email),
		FullName: utils.StrPtrToStr(req.FullName),
		Bio:      utils.StrPtrToStr(req.Bio),
	}

	if err := updateValidator.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := &model.User{
		ID:       req.Id,
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Bio:      req.Bio,
	}

	updatedUser, err := s.userService.Update(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := UpdateResponse{
		User: &pb.User{
			Id:        updatedUser.ID,
			Username:  updatedUser.Username,
			Email:     updatedUser.Email,
			FullName:  updatedUser.FullName,
			Bio:       updatedUser.Bio,
			AvatarUrl: updatedUser.AvatarURL,
		},
	}

	return resp.User, nil
}
