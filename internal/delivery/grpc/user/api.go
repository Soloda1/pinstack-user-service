package user_grpc

import (
	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"pinstack-user-service/internal/logger"
	user_service "pinstack-user-service/internal/service/user"
)

type GRPCServer struct {
	pb.UnimplementedUserServiceServer
	userService *user_service.Service
	log         *logger.Logger
}

func NewGRPCServer(userService *user_service.Service, log *logger.Logger) *GRPCServer {
	return &GRPCServer{
		userService: userService,
		log:         log,
	}
}
