package user_grpc

import (
	"context"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SearchRequest struct {
	Query  string `validate:"omitempty"`
	Offset int32  `validate:"gte=0"`
	Limit  int32  `validate:"gte=1,lte=100"`
}

func (s *UserGRPCService) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	input := SearchRequest{
		Query:  req.Query,
		Offset: req.Offset,
		Limit:  req.Limit,
	}
	if err := validate.Struct(input); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	users, total, err := s.userService.Search(ctx, req.Query, int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.SearchUsersResponse{
		Users: make([]*pb.User, 0, len(users)),
		Total: int64(total),
	}
	for _, u := range users {
		resp.Users = append(resp.Users, &pb.User{
			Id:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FullName:  u.FullName,
			Bio:       u.Bio,
			AvatarUrl: u.AvatarURL,
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		})
	}
	return resp, nil
}
