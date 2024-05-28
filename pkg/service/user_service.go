package service

import (
	"context"
	"fmt"
	"src/pkg/data"
	"src/user_proto"
)

type UserService struct {
	user_proto.UnimplementedUserServiceServer
	Db data.InMemoryUser
}

func (s *UserService) GetUserById(ctx context.Context, req *user_proto.GetUserByIdRequest) (*user_proto.User, error) {
	// Fetch user by ID
	for _, user := range s.Db.Users {
		if user.Id == req.Id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with ID %d not found", req.Id)
}

func (s *UserService) GetUsersByIds(ctx context.Context, req *user_proto.GetUsersByIdsRequest) (*user_proto.Users, error) {
	// Fetch users by IDs
	foundUsers := make([]*user_proto.User, 0)
	for _, id := range req.Ids {
		for _, user := range s.Db.Users {
			if user.Id == id {
				foundUsers = append(foundUsers, user)
				break
			}
		}
	}
	return &user_proto.Users{Users: foundUsers}, nil
}

func (s *UserService) SearchUsers(ctx context.Context, req *user_proto.SearchUsersRequest) (*user_proto.Users, error) {
	// Implement search logic based on criteria
	foundUsers := make([]*user_proto.User, 0)
	for _, user := range s.Db.Users {
		if req.City != "" && user.City == req.City { // Example search
			foundUsers = append(foundUsers, user)
		} else if req.Married && user.Married { // Example search
			foundUsers = append(foundUsers, user)
		}
		// Add more search criteria logic
	}
	return &user_proto.Users{Users: foundUsers}, nil
}
