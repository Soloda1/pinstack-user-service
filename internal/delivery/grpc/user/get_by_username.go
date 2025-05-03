package user_grpc

import (
	"context"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
)

func (s *UserGRPCService) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.User, error) {
	// TODO: реализовать получение пользователя по username
	return nil, nil
}
