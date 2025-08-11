package user_grpc

import (
	user_service "pinstack-user-service/internal/domain/ports/input"
	ports "pinstack-user-service/internal/domain/ports/output"

	"github.com/go-playground/validator/v10"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
)

var validate = validator.New()

type UserGRPCService struct {
	pb.UnimplementedUserServiceServer
	userService user_service.UserService
	log         ports.Logger
}

func NewUserGRPCService(userService user_service.UserService, log ports.Logger) *UserGRPCService {
	return &UserGRPCService{
		userService: userService,
		log:         log,
	}
}
