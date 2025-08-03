package main

import (
	"context"
	"github/com/ammar-nousher-ali/students-api/internal/config"
	"github/com/ammar-nousher-ali/students-api/internal/http/handlers/auth"
	"github/com/ammar-nousher-ali/students-api/internal/http/handlers/course"
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

	//students
	router.HandleFunc("POST /api/students", middleware.JWTMiddleware(student.New(storage)))
	router.HandleFunc("POST /api/students/batch", middleware.JWTMiddleware(student.NewBatch(storage)))
	router.HandleFunc("GET /api/students/{id}", middleware.JWTMiddleware(student.GetById(storage)))
	router.HandleFunc("GET /api/students", middleware.JWTMiddleware(student.GetList(storage)))
	router.HandleFunc("DELETE /api/students/{id}", middleware.JWTMiddleware(student.DeleteStudent(storage)))
	router.HandleFunc("PUT /api/students/{id}", middleware.JWTMiddleware(student.UpdateStudent(storage)))
	router.HandleFunc("GET /api/students/search", middleware.JWTMiddleware(student.SearchStudent(storage)))

	//courses
	router.HandleFunc("POST /api/courses", middleware.JWTMiddleware(course.New(storage)))
	router.HandleFunc("POST /api/courses/batch", middleware.JWTMiddleware(course.NewBatch(storage)))
	router.HandleFunc("GET /api/courses/{id}", middleware.JWTMiddleware(course.GetById(storage)))
	router.HandleFunc("GET /api/courses", middleware.JWTMiddleware(course.GetAll(storage)))
	router.HandleFunc("PUT /api/courses/{id}", middleware.JWTMiddleware(course.Update(storage)))
	router.HandleFunc("DELETE /api/courses/{id}", middleware.JWTMiddleware(course.Delete(storage)))

	corsHandler := enableCORS(router)

	//setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: corsHandler,
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

// CORS middleware function
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
