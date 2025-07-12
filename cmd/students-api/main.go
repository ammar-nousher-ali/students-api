package main

import (
	"context"
	"github/com/ammar-nousher-ali/students-api/internal/config"
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
	//setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("welcome to students api"))

	})

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

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))

	}

	slog.Info("server shutdown successfully.")

}
