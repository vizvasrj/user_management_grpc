package service

import (
	"context"
	"fmt"
	"src/pkg/data"
	"src/user_proto"
	"strconv"
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
	foundUsers := make([]*user_proto.User, 0)

	// Iterate over all criteria
	for _, criterion := range req.Criteria {
		// Define a function to check a single criterion
		checkCriterion := func(user *user_proto.User) bool {
			switch criterion.Field {
			case "city":
				return user.City == criterion.Value
			case "married":
				married, _ := strconv.ParseBool(criterion.Value) // Handle string bool
				return user.Married == married
			case "height":
				height, _ := strconv.ParseFloat(criterion.Value, 32) // Handle string float
				return user.Height == float32(height)
			case "phone":
				phone, _ := strconv.ParseInt(criterion.Value, 10, 64) // Handle string int64
				return user.Phone == phone
			case "fname":
				return user.Fname == criterion.Value
			default:
				return false // Ignore unknown fields
			}
		}

		// Apply the search logic based on the operator
		switch criterion.Operator {
		case user_proto.SearchOperator_AND:
			// AND logic: All criteria must match
			for _, user := range s.Db.Users {
				if checkCriterion(user) {
					foundUsers = append(foundUsers, user)
				}
			}
		case user_proto.SearchOperator_OR:
			// OR logic: At least one criterion must match
			for _, user := range s.Db.Users {
				if checkCriterion(user) {
					foundUsers = append(foundUsers, user)
					break // Move to the next criterion
				}
			}
		default:
			return nil, fmt.Errorf("unknown search operator: %v", criterion.Operator)
		}
	}

	return &user_proto.Users{Users: foundUsers}, nil
}
