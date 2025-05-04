package user_grpc

import (
	"github.com/go-playground/validator/v10"
	"pinstack-user-service/internal/logger"
	user_service "pinstack-user-service/internal/service/user"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
)

var validate = validator.New()

type UserGRPCService struct {
	pb.UnimplementedUserServiceServer
	userService user_service.UserService
	log         *logger.Logger
}

func NewUserGRPCService(userService user_service.UserService, log *logger.Logger) *UserGRPCService {
	return &UserGRPCService{
		userService: userService,
		log:         log,
	}
}
