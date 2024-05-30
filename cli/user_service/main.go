package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"src/pkg/data"
	"src/pkg/faker"
	"src/pkg/service"
	"src/user_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("service is stated")
	useDb := os.Getenv("USE_DATABASE")

	s := grpc.NewServer()

	if useDb == "postgres" {
		db := data.GetConnection()
		// defer db.Close()
		pgInstance := data.NewPostgresDB(db)
		user_proto.RegisterUserServiceServer(s, &service.UserService{
			Db: pgInstance,
		})
	} else if useDb == "inmemory" {
		seedUsers := faker.GetSeed()

		db := data.NewInMemoryUserRepository(seedUsers)
		user_proto.RegisterUserServiceServer(s, &service.UserService{
			Db: db,
		})
	} else {
		log.Fatal("USE_DATABASE not set")
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
