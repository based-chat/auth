// Package grpcclient provides grpc client.
package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"

	srv "github.com/based-chat/auth/pkg/user/v1"
)

const (
	grpcPort           = ":50052"
	grpcHost           = "localhost"
	maxTimeout         = 1 * time.Second
	cntSymbolsPassword = 8
)

// main provides an example of how to use the grpc client to interact with the
// grpc server. It demonstrates how to create, get, update, and delete a user.
//
// It first creates a connection to the grpc server, then creates a user with a
// randomly generated name, email, and password. It then gets the user by their
// id, updates the user's name and email, and finally deletes the user.
func main() {
	// Connect to the grpc server
	addr := net.JoinHostPort(grpcHost, grpcPort)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		connErr := conn.Close()
		if connErr != nil {
			log.Default().Println("failed to close connection: " + connErr.Error())
		}
	}()

	// make a grpc client
	c := srv.NewUserV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	// Create a user
	createResponse, err := c.Create(ctx, &srv.CreateRequest{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, cntSymbolsPassword),
		Role:     srv.UserRole_USER,
	})
	if err != nil {
		log.Default().Println("failed to get: ", err)
		return
	}

	// Get a user
	_, err = c.Get(ctx, &srv.GetRequest{
		Id: createResponse.GetId(),
	})
	if err != nil {
		log.Default().Println("failed to create: ", err)
	}

	// Update a user
	_, err = c.Update(ctx, &srv.UpdateRequest{
		Id:    createResponse.GetId(),
		Name:  &wrapperspb.StringValue{Value: gofakeit.Name()},
		Email: &wrapperspb.StringValue{Value: gofakeit.Email()},
	})
	if err != nil {
		log.Default().Println("failed to update: ", err)
	}

	// Delete a user
	_, err = c.Delete(ctx, &srv.DeleteRequest{
		Id: createResponse.GetId(),
	})
	if err != nil {
		log.Default().Println("failed to delete: ", err)
	}
}
