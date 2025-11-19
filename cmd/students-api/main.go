package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zeus-mtolympus/students-api/internal/config"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Students Api"))
	})

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done // jab tak done khatam nhi hoga tab tak code aage nhi jaayega

	// graceful shutdown
	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx) // may get infinetly stuck
	if err != nil {
		slog.Error("failed to shut down server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
