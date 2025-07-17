package main

import (
	"fmt"
	"log"
	"strings"

	// "log/slog"
	"net"
	// "os"

	_ "github.com/lib/pq"
	"github.com/tousart/sso/config"
	"github.com/tousart/sso/grpc_server/auth"
	"github.com/tousart/sso/repository/kafka"
	"github.com/tousart/sso/repository/postgres"
	"github.com/tousart/sso/repository/redis"
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
	cfgPath := config.ParseFlags()
	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	authRepo, err := postgres.CreateAuthRepo(cfg)
	if err != nil {
		log.Fatalf("failed to create auth repo: %v", err)
	}

	tokenRepo, err := redis.NewTokenRepo(cfg)
	if err != nil {
		log.Fatalf("failed to create token repo: %v", err)
	}

	sender := kafka.NewKafkaSender(strings.Split(cfg.Kafka.Brokers, ","), cfg.Kafka.TopicName)
	defer sender.Writer.Close()

	service := service.CreateAuthService(authRepo, tokenRepo, sender)

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

	// on for tests, off for prod
	reflection.Register(server)
	log.Println("added reflection")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Логгер
	// log := setupLogger(cfg.Env)
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
