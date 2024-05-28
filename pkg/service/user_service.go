package service

import (
	"context"
	"fmt"
	"log"
	"src/pkg/data"
	"src/user_proto"
)

type UserService struct {
	user_proto.UnimplementedUserServiceServer
	Db data.DataStore
}

func (s *UserService) GetUserById(ctx context.Context, req *user_proto.GetUserByIdRequest) (*user_proto.User, error) {
	// Fetch user by ID

	user, err := s.Db.GetUserById(ctx, req.Id)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUsersByIds(ctx context.Context, req *user_proto.GetUsersByIdsRequest) (*user_proto.Users, error) {
	// Fetch users by IDs
	foundUsers, err := s.Db.GetUsersByIds(ctx, req.Ids)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	return &user_proto.Users{Users: foundUsers}, nil
}

func (s *UserService) SearchUsers(ctx context.Context, req *user_proto.SearchUsersRequest) (*user_proto.Users, error) {
	results, err := s.Db.SearchUsers(ctx, req.Criteria)
	if err != nil {
		log.Printf("Error searching users: %v", err)
		return nil, err
	}

	return &user_proto.Users{Users: results}, nil
}
