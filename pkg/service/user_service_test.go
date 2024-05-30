package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"src/pkg/data"
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
	useDb := os.Getenv("USE_DATABASE")
	fmt.Println("USE_DATABASE: ", useDb)
	s := grpc.NewServer()
	if useDb == "postgres" {
		db := data.GetConnection()
		// defer db.Close()

		pgInstance := data.NewPostgresDB(db)

		user_proto.RegisterUserServiceServer(s, &UserService{Db: pgInstance})

	} else if useDb == "inmemory" {
		users := []*user_proto.User{
			{Id: 1, Fname: "Jane", City: "Tokyo", Phone: 1764788495, Height: 4.9, Married: true},
			{Id: 2, Fname: "William", City: "Chicago", Phone: 3147005163, Height: 6.0, Married: true},
			{Id: 3, Fname: "Jane", City: "Sydney", Phone: 5660612778, Height: 5.2, Married: false},
			{Id: 4, Fname: "Olivia", City: "Paris", Phone: 7499978875, Height: 5.0, Married: false},
			{Id: 5, Fname: "Olivia", City: "New York", Phone: 4934669609, Height: 5.8, Married: true},
			{Id: 6, Fname: "Olivia", City: "Berlin", Phone: 3935422070, Height: 6.2, Married: false},
			{Id: 7, Fname: "Peter", City: "Chicago", Phone: 1389433165, Height: 5.7, Married: false},
			{Id: 8, Fname: "Susan", City: "Paris", Phone: 5307397290, Height: 6.1, Married: false},
			{Id: 9, Fname: "Olivia", City: "Madrid", Phone: 5199527895, Height: 4.6, Married: false},
			{Id: 10, Fname: "Olivia", City: "Madrid", Phone: 8354340417, Height: 5.5, Married: true},
		}
		db := data.NewInMemoryUserRepository(users)
		user_proto.RegisterUserServiceServer(s, &UserService{Db: db})
	} else {
		log.Fatal("USE_DATABASE not set")
	}
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
	user, err := c.GetUserById(ctx, &user_proto.GetUserByIdRequest{Id: 2})
	if err != nil {
		t.Fatalf("GetUserById failed: %v", err)
	}
	if user.Id != 2 || user.Fname != "William" {
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
		City: "Madrid",
		// Height: &user_proto.HeightRange{
		// 	StartValue: 5.0,
		// 	EndValue:   7.0,
		// },
	})

	if err != nil {
		t.Fatalf("SearchUsers failed: %v", err)
	}

	if results.Users[0].City != "Madrid" {
		t.Errorf("Unexpected user data: %v ", results.Users[0])
	}

	if len(results.Users) == 0 {
		t.Fatalf("No users found")
	}
}
