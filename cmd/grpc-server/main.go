// Package grpcserver provides grpc server.
package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	srv "github.com/based-chat/auth/pkg/user/v1"
	"github.com/based-chat/auth/utilites/mathx"
)

const (
	gprcPort = ":50052"
	grmcHost = "localhost"
)

type server struct {
	srv.UnimplementedUserV1Server
}

// Create creates a user with a randomly generated name and email.
// It returns the user's id in the response.
func (s *server) Create(_ context.Context, _ *srv.CreateRequest) (*srv.CreateResponse, error) {
	return &srv.CreateResponse{
		Id: mathx.Abs(gofakeit.Int64()),
	}, nil
}

// Get retrieves a user by their id.
// It returns the user's id, name, email, role, created_at, and updated_at in the response.
// The user's details are randomly generated.
func (s *server) Get(_ context.Context, _ *srv.GetRequest) (*srv.GetResponse, error) {
	return &srv.GetResponse{
		Id:        mathx.Abs(gofakeit.Int64()),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      srv.UserRole_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

// Update updates a user by their id.
// It returns the user's id, name, email, role, created_at, and updated_at in the response.
// The user's details are randomly generated.
func (s *server) Update(_ context.Context, _ *srv.UpdateRequest) (*srv.GetResponse, error) {
	return &srv.GetResponse{
		Id:        mathx.Abs(gofakeit.Int64()),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      srv.UserRole_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *server) Delete(_ context.Context, _ *srv.DeleteRequest) (*srv.DeleteResponse, error) {
	return &srv.DeleteResponse{
		Deleted: true,
	}, nil
}

// main starts the grpc server and listens on the specified address.
// It seeds the random number generator and registers the user service.
// It then serves the grpc server and logs any errors that occur during serving.
func main() {
	webAddress := net.JoinHostPort(grmcHost, gprcPort)
	addr, err := net.ResolveTCPAddr("tcp", webAddress)
	if err != nil {
		log.Fatalf("failed to resolve address: %v", err)
	}

	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err = gofakeit.Seed(time.Now().UnixNano()); err != nil {
		log.Default().Printf("failed to seed random number generator: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	srv.RegisterUserV1Server(s, &server{})
	err = s.Serve(listen)
	if err != nil {
		log.Default().Println("failed to serve:", err.Error())
	}
}
