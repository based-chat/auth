// Package grpcserver provides grpc server.
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	srv "github.com/based-chat/auth/pkg/user/v1"
	"github.com/based-chat/auth/utilites/mathx"
)

const (
	grpcPort              = "50052"
	grpcHost              = "localhost"
	errorIDInvalid        = "invalid ID"
	errorNameRequired     = "name is required"
	errorEmailRequired    = "email is required"
	errorPasswordRequired = "password is required"
)

var (
	errFailedListen = errors.New("failed to listen")
	errFailedServe  = errors.New("failed to serve")
	errFailedSeed   = errors.New("failed to seed fakeit")
)

type server struct {
	srv.UnimplementedUserV1Server
}

// Create creates a user with a randomly generated name and email.
// It returns the user's id in the response.
func (s *server) Create(_ context.Context, req *srv.CreateRequest) (*srv.CreateResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, errorNameRequired)
	}

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, errorEmailRequired)
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, errorPasswordRequired)
	}

	return &srv.CreateResponse{
		Id: mathx.Abs(gofakeit.Int64()),
	}, nil
}

// Get retrieves a user by their id.
// It returns the user's id, name, email, role, created_at, and updated_at in the response.
// The user's details are randomly generated.
func (s *server) Get(_ context.Context, req *srv.GetRequest) (*srv.GetResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, errorIDInvalid)
	}
	// Return a random user
	date := gofakeit.Date()

	return &srv.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      srv.UserRole_USER,
		CreatedAt: timestamppb.New(date),
		UpdatedAt: timestamppb.New(date.Add(time.Minute)),
	}, nil
}

// Update updates a user with a randomly generated name and email.
// It returns the user's id, name, email, role, created_at, and updated_at in the response.
// The user's details are randomly generated.
// If the request ID is invalid (less than or equal to 0), it returns an error.
func (s *server) Update(_ context.Context, req *srv.UpdateRequest) (*srv.GetResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, errorIDInvalid)
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
		CreatedAt: timestamppb.Now(),
		UpdatedAt: timestamppb.Now(),
	}, nil
}

// Delete deletes a user by their ID.
// It returns an error if the user ID is invalid (less than or equal to 0).
// It returns a DeleteResponse with the "deleted" field set to true if the user was deleted.
func (s *server) Delete(_ context.Context, req *srv.DeleteRequest) (*srv.DeleteResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, errorIDInvalid)
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
	var lc net.ListenConfig

	listen, err := lc.Listen(context.Background(), "tcp", net.JoinHostPort(grpcHost, grpcPort))
	if err != nil {
		log.Fatalf(errFailedListen.Error(), err)
	}

	// Seed the random number generator
	if err := gofakeit.Seed(time.Now().UnixNano()); err != nil {
		log.Fatalf(errFailedSeed.Error(), err)
	}

	// Start the grpc server
	s := grpc.NewServer()
	reflection.Register(s)
	srv.RegisterUserV1Server(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Fatalf(errFailedServe.Error(), err)
	}
}
