package data

import (
	"fmt"
	"src/user_proto"
)

type InMemoryUser struct {
	Users []*user_proto.User
}

func NewInMemoryUserRepository(users []*user_proto.User) *InMemoryUser {
	return &InMemoryUser{Users: users}
}

func (r *InMemoryUser) GetUserById(id int32) (*user_proto.User, error) {
	for _, user := range r.Users {
		if user.Id == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with ID %d not found", id)
}

func (r *InMemoryUser) GetUsersByIds(ids []int32) ([]*user_proto.User, error) {
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

func (r *InMemoryUser) SearchUsers(city string, married bool) ([]*user_proto.User, error) {
	foundUsers := make([]*user_proto.User, 0)
	for _, user := range r.Users {
		if city != "" && user.City == city {
			foundUsers = append(foundUsers, user)
		} else if married && user.Married {
			foundUsers = append(foundUsers, user)
		}
	}
	return foundUsers, nil
}
