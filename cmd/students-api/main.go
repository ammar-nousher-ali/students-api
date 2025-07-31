package main

import (
	"context"
	"github/com/ammar-nousher-ali/students-api/internal/config"
	"github/com/ammar-nousher-ali/students-api/internal/http/handlers/auth"
	"github/com/ammar-nousher-ali/students-api/internal/http/handlers/student"
	"github/com/ammar-nousher-ali/students-api/internal/middleware"
	"github/com/ammar-nousher-ali/students-api/internal/storage/sqlite"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//load config

	cfg := config.MustLoad()

	//database setup

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)

	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router

	router := http.NewServeMux()

	//Public routes
	router.HandleFunc("POST /api/signup", auth.Signup(storage))
	router.HandleFunc("POST /api/signin", auth.SignIn(storage))

	//Protected routes
	router.HandleFunc("POST /api/students", middleware.JWTMiddleware(student.New(storage)))
	router.HandleFunc("POST /api/students/batch", middleware.JWTMiddleware(student.NewBatch(storage)))
	router.HandleFunc("GET /api/students/{id}", middleware.JWTMiddleware(student.GetById(storage)))
	router.HandleFunc("GET /api/students", middleware.JWTMiddleware(student.GetList(storage)))
	router.HandleFunc("DELETE /api/students/{id}", middleware.JWTMiddleware(student.DeleteStudent(storage)))
	router.HandleFunc("PUT /api/students/{id}", middleware.JWTMiddleware(student.UpdateStudent(storage)))
	router.HandleFunc("GET /api/students/search", middleware.JWTMiddleware(student.SearchStudent(storage)))

	//setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server %s", err)

		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully.")

}
