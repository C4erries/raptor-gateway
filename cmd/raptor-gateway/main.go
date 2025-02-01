package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/c4erries/raptor-proto/userpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	// Создаем контекст
	ctx := context.Background()

	// Создаем мультиплексор
	mux := runtime.NewServeMux()

	// Опции для подключения к gRPC-сервису
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Регистрируем UserService
	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// Запускаем HTTP-сервер
	log.Println("API Gateway is running on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
