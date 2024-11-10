package main

import (
	"context"
	"github.com/vlks-dev/sibavto/auth-service/internal/handlers"
	"github.com/vlks-dev/sibavto/auth-service/internal/services"
	"github.com/vlks-dev/sibavto/auth-service/internal/storage"
	"github.com/vlks-dev/sibavto/auth-service/proto/authpb"
	"github.com/vlks-dev/sibavto/shared/logger"
	"github.com/vlks-dev/sibavto/shared/utils/config"
	"github.com/vlks-dev/sibavto/shared/utils/token"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := config.MustParse("config.yaml")
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	slog := logger.NewSlog(config)
	postgresPool, err := storage.PostgresPool(ctx, config, slog)
	if err != nil {
		slog.Error("Postgres pool err: ", err.Error())
	}

	userStorage := storage.NewUserStorage(postgresPool, slog)
	jwtService := token.NewJWTService(config)
	authService := services.NewAuthService(userStorage, slog, jwtService)
	authServiceServer := handlers.NewAuthServiceServer(authService)
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, authServiceServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("Failed to listen on port 50051", err.Error())
	}

	slog.Info("gRPC Auth Server is running",
		"Port:", listener.Addr().(*net.TCPAddr).Port,
		"time:", time.Now().Format("2006-01-02 15:04:05"),
		"server", server.GetServiceInfo(),
	)
	go func() {
		if err := server.Serve(listener); err != nil {
			slog.Error("Failed to serve gRPC server", err.Error())
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Блокируем выполнение, ожидая сигнал завершения
	<-stopChan
	slog.Info("Shutting down gRPC server...")

	/*	uid, _ := uuid.Parse("4ac9e202-6b0b-474e-8514-7b13f69dcb7f")
		keyString := fmt.Sprintf("user_token:%s", uid)
		err = redisConn.Get(ctx, keyString).Err()
		if err != nil {
			slog.Error("Failed to get user token", err.Error())
			return
		}*/

	// Плавная остановка сервера
	server.GracefulStop()
	slog.Info("gRPC server stopped gracefully.")
}
