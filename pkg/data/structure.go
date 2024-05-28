package data

import (
	"context"
	"src/user_proto"
)

// DataStore defines the interface for data access operations.
type DataStore interface {
	GetUserById(ctx context.Context, id int32) (*user_proto.User, error)
	GetUsersByIds(ctx context.Context, ids []int32) ([]*user_proto.User, error)
	SearchUsers(ctx context.Context, criteria []*user_proto.SearchCriteria) ([]*user_proto.User, error)
	// Add more methods as needed
}
