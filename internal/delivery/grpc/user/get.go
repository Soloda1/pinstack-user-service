package user_grpc

import (
	"context"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
)

func (s *UserGRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	// TODO: реализовать получение пользователя по ID
	return nil, nil
}
