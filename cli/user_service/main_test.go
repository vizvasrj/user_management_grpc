package main

import (
	"context"
	"fmt"
	"net"
	"src/pkg/data"
	"src/pkg/service"
	"src/user_proto"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	fmt.Println("Test service is stated")
	lis = bufconn.Listen(bufSize)
	db := data.GetConnection()
	// defer db.Close()

	pgInstance := data.NewPostgresDB(db)

	s := grpc.NewServer()
	user_proto.RegisterUserServiceServer(s, &service.UserService{Db: pgInstance})
	go s.Serve(lis)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGetUserById(t *testing.T) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	c := user_proto.NewUserServiceClient(conn)
	ctx := context.Background()

	// Test fetching an existing user
	user, err := c.GetUserById(ctx, &user_proto.GetUserByIdRequest{Id: 1})
	if err != nil {
		t.Fatalf("GetUserById failed: %v", err)
	}
	if user.Id != 1 || user.Fname != "Jane" {
		t.Errorf("Unexpected user data: %v", user)
	}

}

func TestGetUsersByIds(t *testing.T) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	c := user_proto.NewUserServiceClient(conn)
	ctx := context.Background()

	// Test fetching existing users
	users, err := c.GetUsersByIds(ctx, &user_proto.GetUsersByIdsRequest{Ids: []int32{1, 2}})
	if err != nil {
		t.Fatalf("GetUsersByIds failed: %v", err)
	}
	if len(users.Users) != 2 {
		t.Errorf("Unexpected number of users: %v", len(users.Users))
	}

}

func TestSearchUsers(t *testing.T) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	c := user_proto.NewUserServiceClient(conn)
	ctx := context.Background()

	// Test searching for users
	results, err := c.SearchUsers(ctx, &user_proto.SearchUsersRequest{
		Filters:      map[string]string{"city": "Tokyo"},
		RangeFilters: map[string]string{"height": "4.0 6.0"},
	})

	if err != nil {
		t.Fatalf("SearchUsers failed: %v", err)
	}
	if len(results.Users) != 1 {
		t.Errorf("Unexpected number of users: %v", len(results.Users))
	}
	if results.Users[0].City != "Tokyo" {
		t.Errorf("Unexpected user data: %v ", results.Users[0])
	}
}
