// Package grpcclient provides grpc client.
package main

import (
	"context"
	"errors"
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
	grpcPort           = "50052"
	grpcHost           = "localhost"
	maxTimeout         = 1 * time.Second
	cntSymbolsPassword = 8
)

var (
	errFailedConnect         = errors.New("failed to connect: %v")
	errFailedCloseConnection = "failed to close connection: %v"
	errFailedGetUser         = "failed to get user: %v"
	errFailedCreateUser      = "failed to create user: %v"
	errFailedUpdateUser      = "failed to update user: %v"
	errFailedDeleteUser      = "failed to delete user: %v"
)

// main provides an example of how to use the grpc client to interact with the
// grpc server. It demonstrates how to create, get, update, and delete a user.
//
// It first creates a connection to the grpc server, then creates a user with a
// randomly generated name, email, and password. It then gets the user by their
// id, updates the user's name and email, and finally deletes the user.
func main() {
	addr := net.JoinHostPort(grpcHost, grpcPort)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Default().Printf(errFailedConnect.Error(), err)
		return
	}

	defer func() {
		connErr := conn.Close()
		if connErr != nil {
			log.Default().Printf(errFailedCloseConnection, connErr)
		}
	}()

	c := srv.NewUserV1Client(conn)

	ctxCreateUser, cancelCreateUser := context.WithTimeout(context.Background(), maxTimeout)

	createResponse, err := c.Create(ctxCreateUser, &srv.CreateRequest{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, false, cntSymbolsPassword),
		Role:     srv.UserRole_USER,
	})
	if err != nil {
		log.Default().Printf(errFailedCreateUser, err)
		return
	}

	cancelCreateUser()

	ctxGetUser, cancelGetUser := context.WithTimeout(context.Background(), maxTimeout)

	_, err = c.Get(ctxGetUser, &srv.GetRequest{
		Id: createResponse.GetId(),
	})
	if err != nil {
		log.Default().Printf(errFailedGetUser, err)
	}

	cancelGetUser()

	ctxUpdateUser, cancelUpdateUser := context.WithTimeout(context.Background(), maxTimeout)

	_, err = c.Update(ctxUpdateUser, &srv.UpdateRequest{
		Id:    createResponse.GetId(),
		Name:  &wrapperspb.StringValue{Value: gofakeit.Name()},
		Email: &wrapperspb.StringValue{Value: gofakeit.Email()},
	})
	if err != nil {
		log.Default().Printf(errFailedUpdateUser, err)
	}

	cancelUpdateUser()

	ctxDeleteUser, cancelDeleteUser := context.WithTimeout(context.Background(), maxTimeout)

	_, err = c.Delete(ctxDeleteUser, &srv.DeleteRequest{
		Id: createResponse.GetId(),
	})
	if err != nil {
		log.Default().Printf(errFailedDeleteUser, err)
	}

	cancelDeleteUser()
}
