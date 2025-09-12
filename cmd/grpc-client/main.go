// Package grpcclient provides grpc client.
package grpcclient

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
	gprcPort            = ":50052"
	grmcHost            = "localhost"
	maxTimeout          = 1 * time.Second
	cntSymbolsPasssword = 8
)

// main provides an example of how to use the grpc client to interact with the
// grpc server. It demonstrates how to create, get, update, and delete a user.
//
// It first creates a connection to the grpc server, then creates a user with a
// randomly generated name, email, and password. It then gets the user by their
// id, updates the user's name and email, and finally deletes the user.
func main() {
	conn, err := grpc.NewClient(net.JoinHostPort(grmcHost, gprcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		connErr := conn.Close()
		if connErr != nil {
			log.Default().Println("failed to close connection: " + connErr.Error())
		}
	}()

	c := srv.NewUserV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	createResponse, err := c.Create(ctx, &srv.CreateRequest{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, cntSymbolsPasssword),
		Role:     srv.UserRole_USER,
	})
	if err != nil {
		log.Default().Println("failed to create: ", err)
	}

	_, err = c.Get(ctx, &srv.GetRequest{
		Id: createResponse.GetId(),
	})
	if err != nil {
		log.Default().Println("failed to create: ", err)
	}

	_, err = c.Update(ctx, &srv.UpdateRequest{
		Id:    createResponse.GetId(),
		Name:  &wrapperspb.StringValue{Value: gofakeit.Name()},
		Email: &wrapperspb.StringValue{Value: gofakeit.Email()},
	})
	if err != nil {
		log.Default().Print("failed to update: ", err)
		log.Default().Println()
	}

	_, err = c.Delete(ctx, &srv.DeleteRequest{
		Id: createResponse.GetId(),
	})
	if err != nil {
		log.Default().Println("failed to delete: ", err)
	}
}
