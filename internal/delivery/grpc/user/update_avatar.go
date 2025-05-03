package user_grpc

import (
	"context"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserGRPCService) UpdateAvatar(ctx context.Context, req *pb.UpdateAvatarRequest) (*emptypb.Empty, error) {
	// TODO: реализовать обновление аватара пользователя
	return nil, nil
}
