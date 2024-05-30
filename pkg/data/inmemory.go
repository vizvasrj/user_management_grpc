package data

import (
	"context"
	"fmt"
	"src/user_proto"
)

type InMemoryUser struct {
	Users []*user_proto.User
}

func NewInMemoryUserRepository(users []*user_proto.User) *InMemoryUser {
	return &InMemoryUser{Users: users}
}

func (r *InMemoryUser) DbGetUserById(ctx context.Context, id int32) (*user_proto.User, error) {
	for _, user := range r.Users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with ID %d not found", id)
}

func (r *InMemoryUser) DbGetUsersByIds(ctx context.Context, ids []int32) ([]*user_proto.User, error) {
	foundUsers := make([]*user_proto.User, 0)
	for _, id := range ids {
		for _, user := range r.Users {
			if user.Id == id {
				foundUsers = append(foundUsers, user)
				break
			}
		}
	}
	return foundUsers, nil
}

func (s *InMemoryUser) DbSearchUsers(ctx context.Context, req *user_proto.SearchUsersRequest) ([]*user_proto.User, error) {
	var filteredUsers []*user_proto.User
	for _, u := range s.Users {
		// Check all fields based on the request
		if req.Id != 0 && u.Id != req.Id {
			continue
		}
		if req.Fname != "" && u.Fname != req.Fname {
			continue
		}
		if req.City != "" && u.City != req.City {
			continue
		}
		if req.Phone != 0 && u.Phone != req.Phone {
			continue
		}
		// Handle married (boolean or string)
		if req.Married != nil {
			if req.Married.IsMarried != u.Married {
				continue
			}
		}
		if req.Height != nil {
			if req.Height.StartValue != 0 || req.Height.EndValue != 0 {
				if u.Height < req.Height.StartValue || u.Height > req.Height.EndValue {
					continue
				}
			}
		}
		filteredUsers = append(filteredUsers, u)
	}
	return filteredUsers, nil
}
