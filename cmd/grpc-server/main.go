// Package grpcserver provides grpc server.
package main

import (
	"context"
	"fmt"
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
	grpcPort = ":50052"
	grpcHost = "localhost"
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
func (s *server) Get(_ context.Context, req *srv.GetRequest) (*srv.GetResponse, error) {
	if req.GetId() <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", req.GetId())
	}
	// Return a random user
	return &srv.GetResponse{
		Id:        req.GetId(), // Возвращаем тот же ID, что был запрошен
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      srv.UserRole_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

// Update updates a user with a randomly generated name and email.
// It returns the user's id, name, email, role, created_at, and updated_at in the response.
// The user's details are randomly generated.
// If the request ID is invalid (less than or equal to 0), it returns an error.
func (s *server) Update(_ context.Context, req *srv.UpdateRequest) (*srv.GetResponse, error) {
	if req.GetId() <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", req.GetId())
	}

	// set random name or request name
	name := gofakeit.Name()
	if req.GetName() != nil {
		name = req.GetName().GetValue()
	}

	// set random email or request email
	email := gofakeit.Email()
	if req.GetEmail() != nil {
		email = req.GetEmail().GetValue()
	}

	// Return a random user
	return &srv.GetResponse{
		Id:        req.GetId(),
		Name:      name,
		Email:     email,
		Role:      srv.UserRole_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(time.Now()), // Используем текущее время для updated_at
	}, nil
}

// Delete deletes a user by their ID.
// It returns an error if the user ID is invalid (less than or equal to 0).
// It returns a DeleteResponse with the "deleted" field set to true if the user was deleted.
func (s *server) Delete(_ context.Context, req *srv.DeleteRequest) (*srv.DeleteResponse, error) {
	if req.GetId() <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", req.GetId())
	}
	// Return fact that user was deleted
	return &srv.DeleteResponse{
		Deleted: true,
	}, nil
}

// main starts the grpc server and listens on the specified address.
// It seeds the random number generator and registers the user service.
// It then serves the grpc server and logs any errors that occur during serving.
func main() {
	// Listen on the specified address
	webAddress := net.JoinHostPort(grpcPort, grpcHost)
	addr, err := net.ResolveTCPAddr("tcp", webAddress)
	if err != nil {
		log.Fatalf("failed to resolve address: %v", err)
	}
	listen, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Seed the random number generator
	if err = gofakeit.Seed(time.Now().UnixNano()); err != nil {
		log.Default().Printf("failed to seed random number generator: %v", err)
	}

	// Start the grpc server
	s := grpc.NewServer()
	reflection.Register(s)
	srv.RegisterUserV1Server(s, &server{})
	err = s.Serve(listen)
	if err != nil {
		log.Default().Println("failed to serve:", err.Error())
	}
}
