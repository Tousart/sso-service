package main

import (
	"fmt"
	"log"

	// "log/slog"
	"net"
	// "os"

	_ "github.com/lib/pq"
	"github.com/tousart/sso/config"
	"github.com/tousart/sso/grpc_server/auth"
	"github.com/tousart/sso/repository/kafka"
	"github.com/tousart/sso/repository/postgres"
	"github.com/tousart/sso/usecase/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// const (
// 	envLocal = "local"
// 	envDev   = "dev"
// 	envProd  = "prod"
// )

func main() {
	cfg := config.MustLoad()

	repo, err := postgres.CreateAuthRepo()
	if err != nil {
		log.Fatalf("failed to create repo: %v", err)
	}

	sender := kafka.NewKafkaSender([]string{"kafka:9093"}, "email_messages")
	defer sender.Writer.Close()

	service := service.CreateAuthService(repo, sender)

	serverAPI := auth.CreateServerAPI(service)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	log.Println("listener has been created")
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	server := grpc.NewServer()
	log.Println("grpc server has been created")

	auth.Register(server, serverAPI)
	log.Println("auth server has been registered")

	reflection.Register(server)
	log.Println("added reflection")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Логгер
	// log := setupLogger(cfg.Env)

	// Приложение

	// Запуск gRPC-сервера
}

// func setupLogger(env string) *slog.Logger {
// 	var log *slog.Logger

// 	switch env {
// 	case envLocal:
// 		log = slog.New(
// 			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
// 		)
// 	case envDev:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
// 		)
// 	case envProd:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
// 		)
// 	}

// 	return log
// }
