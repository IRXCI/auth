package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/IRXCI/auth/pkg/auth"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserAPIServer
}

func (s *server) CreateUser(_ context.Context, _ *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	log.Printf("CreateUser is working...")

	return &desc.CreateUserResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) GetUser(_ context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	log.Printf("Note id: %d", req.GetId())

	return &desc.GetUserResponse{
		Id: req.GetId(),
		UserAuth: &desc.User{
			Name:  gofakeit.Name(),
			Email: gofakeit.Email(),
			Role:  desc.Role_USER},

		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *server) UpdateUser(_ context.Context, _ *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	log.Printf("User updated")
	return nil, nil
}

func (s *server) DeleteUser(_ context.Context, _ *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Printf("User delete")
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserAPIServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
