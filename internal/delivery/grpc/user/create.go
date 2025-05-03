package user_grpc

import (
	"context"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
)

func (s *UserGRPCService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// TODO: реализовать создание пользователя
	return nil, nil
}
