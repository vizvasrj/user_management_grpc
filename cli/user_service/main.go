package main

import (
	"fmt"
	"log"
	"net"
	"src/pkg/data"
	"src/pkg/faker"
	"src/pkg/service"
	"src/user_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("service is stated")
	fakeUsers := faker.GenerateFakeUsers(10)
	users := []*user_proto.User{}
	// user := &user_proto.User{Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true}
	users = append(users, fakeUsers...)

	db := data.NewInMemoryUserRepository(users)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	user_proto.RegisterUserServiceServer(s, &service.UserService{Db: *db})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
