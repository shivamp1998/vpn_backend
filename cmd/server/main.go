package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

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
	defer database.Disconnect()

	if err != nil {
		log.Fatal("Error in connection to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = database.InitializeIndexes(ctx)

	if err != nil {
		log.Printf("Warning: Failed to initialize indexes: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":50051"
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Error in listening to port", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(server.AuthInterceptor),
	)

	mainServer := server.NewServer()
	pb.RegisterUserServiceServer(grpcServer, mainServer)
	pb.RegisterServerServiceServer(grpcServer, mainServer)
	pb.RegisterConfigServiceServer(grpcServer, mainServer)
	reflection.Register(grpcServer)

	fmt.Print("Server connected on port", port)
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatal("Error in serving the grpc server")
	}
}
