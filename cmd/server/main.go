package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/shivamp1998/vpn_backend/internal/database"
	server "github.com/shivamp1998/vpn_backend/internal/server"
	pb "github.com/shivamp1998/vpn_backend/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error in loading environment file")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	err = database.Connect(MONGODB_URI)

	if err != nil {
		log.Fatal("Error in connection to database", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":50051"
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Error in listening to port", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &server.Server{})
	reflection.Register(grpcServer)

	fmt.Print("Server connected on port", port)
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatal("Error in serving the grpc server")
	}
}
