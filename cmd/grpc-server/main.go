// Package grpcserver provides grpc server.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"time"

	"github.com/based-chat/auth/internal/config"
	"github.com/based-chat/auth/internal/config/env"
	"github.com/based-chat/auth/internal/mathx"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	srv "github.com/based-chat/auth/pkg/user/v1"
)

var configPath string

// init инициализирует генератор случайных данных gofakeit
// регистрирует флаг командной строки `--config-path` (по умолчанию ".env").
// При ошибке сидирования она записывается в стандартный лог.
func init() {
	err := gofakeit.Seed(time.Now().UnixNano())
	if err != nil {
		log.Default().Println(errFailedSeed, err)
	}

	flag.StringVar(&configPath, "config-path", ".env", "config path")
}

const (
	errorIDInvalid        = "invalid ID"
	errorNameRequired     = "name is required"
	errorEmailRequired    = "email is required"
	errorPasswordRequired = "password is required"
)

var (
	errFailedListen          = errors.New("failed to listen")
	errFailedServe           = errors.New("failed to serve")
	errFailedSeed            = errors.New("failed to seed fakeit")
	errFailedLoadConfig      = errors.New("failed to load config")
	errFailedConnect         = errors.New("failed to connect")
	errFailedCloseConnection = errors.New("failed to close connection")
)

type server struct {
	srv.UnimplementedUserV1Server
}

// Create создает нового пользователя. В текущей реализаци его id генерируется случайным образом.
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

// Get возвращает случайного пользователя.
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

// Update обновляет пользователя.
// В текущей реализации он возвращает случайно созданного пользователя.
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

// Delete удаляет пользователя.
// В текущей реализации он возвращает успешное удаление.
func (s *server) Delete(_ context.Context, req *srv.DeleteRequest) (*srv.DeleteResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, errorIDInvalid)
	}
	// Return fact that user was deleted
	return &srv.DeleteResponse{
		Deleted: true,
	}, nil
}

// main запускает gRPC-сервер для сервиса UserV1.
//
// Функция:
// - загружает конфигурацию из файла окружения (config.Load(".env")) и формирует gRPC и Postgres конфиги;
// - открывает TCP-листенер по адресу gRPC-конфига (gRPCConfig.Address());
// - устанавливает подключение к PostgreSQL через pgx и откладывает его закрытие;
// - создаёт gRPC-сервер, регистрирует reflection и реализацию UserV1, после чего начинает
// обслуживать входящие соединения.
// В случае ошибок загрузки конфигурации, создания листенера или установления подключения к БД функция
// завершает процесс с логированием через log.Fatalf.
// Ошибки во время работы s.Serve() логируются без явного завершения процесса.
func main() {
	flag.Parse()

	ctx := context.Background()

	// Initialize the config
	if err := config.Load(".env"); err != nil {
		log.Fatalf("%s: %v", errFailedLoadConfig.Error(), err)
	}

	gRPCConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("%s: %v", errFailedLoadConfig.Error(), err)
	}

	// Listen on the specified address
	var lc net.ListenConfig

	listen, err := lc.Listen(context.Background(), "tcp", gRPCConfig.Address())
	if err != nil {
		log.Fatalf("%s: %v", errFailedListen.Error(), err)
	}

	postgresConfig, err := env.NewPostgresConfig()
	if err != nil {
		log.Fatalf("%s: %v", errFailedLoadConfig.Error(), err)
	}

	conn, err := pgx.Connect(ctx, postgresConfig.DSN())
	if err != nil {
		log.Fatalf("%s: %v", errFailedConnect.Error(), err)
	}

	defer func() {
		err := conn.Close(ctx)
		if err != nil {
			log.Printf("%s: %v", errFailedCloseConnection.Error(), err)
		}
	}()

	// Start the grpc server
	s := grpc.NewServer()
	reflection.Register(s)
	srv.RegisterUserV1Server(s, &server{})

	if err = s.Serve(listen); err != nil {
		log.Printf("%s: %v", errFailedServe.Error(), err)
	}
}
