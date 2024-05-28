package data

// import (
// 	"context"
// 	"fmt"
// 	"src/user_proto"
// 	"strconv"
// )

// type InMemoryUser struct {
// 	Users []*user_proto.User
// }

// func NewInMemoryUserRepository(users []*user_proto.User) *InMemoryUser {
// 	return &InMemoryUser{Users: users}
// }

// func (r *InMemoryUser) GetUserById(ctx context.Context, id int32) (*user_proto.User, error) {
// 	for _, user := range r.Users {
// 		if user.Id == id {
// 			return user, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("user with ID %d not found", id)
// }

// func (r *InMemoryUser) GetUsersByIds(ctx context.Context, ids []int32) ([]*user_proto.User, error) {
// 	foundUsers := make([]*user_proto.User, 0)
// 	for _, id := range ids {
// 		for _, user := range r.Users {
// 			if user.Id == id {
// 				foundUsers = append(foundUsers, user)
// 				break
// 			}
// 		}
// 	}
// 	return foundUsers, nil
// }

// func (r *InMemoryUser) SearchUsers(ctx context.Context, req *user_proto.SearchUsersRequest) ([]*user_proto.User, error) {
// 	foundUsers := make([]*user_proto.User, 0)

// 	// Iterate over all criteria
// 	for _, criterion := range req.Criteria {
// 		// Define a function to check a single criterion
// 		checkCriterion := func(user *user_proto.User) bool {
// 			switch criterion.Field {
// 			case "city":
// 				return user.City == criterion.Value
// 			case "married":
// 				married, err := strconv.ParseBool(criterion.Value) // Handle string bool
// 				if err != nil {
// 					return false
// 				}
// 				return user.Married == married
// 			case "height":
// 				height, _ := strconv.ParseFloat(criterion.Value, 32) // Handle string float
// 				return user.Height == float32(height)
// 			case "phone":
// 				phone, _ := strconv.ParseInt(criterion.Value, 10, 64) // Handle string int64
// 				return user.Phone == phone
// 			case "fname":
// 				return user.Fname == criterion.Value
// 			default:
// 				return false // Ignore unknown fields
// 			}
// 		}

// 		// Apply the search logic
// 		for _, user := range r.Users {
// 			if checkCriterion(user) {
// 				foundUsers = append(foundUsers, user)
// 			}
// 			// Remove the break statement - We want to check all users
// 		}
// 	}
// 	return foundUsers, nil
// }
