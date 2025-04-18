package server

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/c4erries/raptor-gateway/internal/config"
)

type Server struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Server {
	return &Server{log: log}
}

// Запуск сервера. Регистрация сервисов и подъём http сервера
func (s *Server) Start(config *config.Config) {
	// Создаем контекст
	//ctx := context.Background()

	// Создаем мультиплексор
	mux := http.NewServeMux()

	// Запускаем HTTP-сервер
	s.log.Debug("API Gateway is running on" + config.HTTPServer.Addr)
	if err := http.ListenAndServe(config.HTTPServer.Addr,
		s.errorHandlingMiddleware(s.logmiddleware(mux))); err != nil {
		s.log.Error("failed to serve:", slog.Any("error", err))
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

func (s *Server) logmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Debug("Запрос на %v с методом %v", r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) errorHandlingMiddleware(next http.Handler) http.Handler {
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

func (s *Server) Stop() {

}
