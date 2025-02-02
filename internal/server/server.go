package server

import (
	"context"
	"log"
	"net/http"

	userpb "github.com/c4erries/raptor-proto/userpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server interface {
	Start(c *Config)
}

// Запуск сервера. Регистрация сервисов и подъём http сервера
func Start(config *Config) {
	// Создаем контекст
	ctx := context.Background()

	// Создаем мультиплексор
	mux := runtime.NewServeMux()

	// Опции для подключения к gRPC-сервису
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Регистрируем UserService
	if err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.UserServiceUrl, opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// Запускаем HTTP-сервер
	log.Println("API Gateway is running on port 8080")
	if err := http.ListenAndServe(config.Bindaddr,
		errorHandlingMiddleware(
			asyncmiddleware(
				logmiddleware(mux)))); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func asyncmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})

		// Копируем контекст и запрос для асинхронной обработки
		ctx := r.Context()
		req := r.Clone(ctx)

		// Запускаем асинхронную обработку
		go func() {
			defer close(done)
			next.ServeHTTP(w, req)
		}()

		// Ждем завершения обработки
		<-done
	})
}

func logmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func errorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from error: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
