package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.statusCode = status
	rw.ResponseWriter.WriteHeader(status)
}

func main() {
	// load .env file — only in development
	// in production env vars are set by the platform
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found — using environment variables")
	}

	// load and validate config
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("config error:", err)
	}

	// open database
	db, err := openDB(cfg.DBPath)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	// setup stores
	taskStore, err := NewSQLiteStore(db)
	if err != nil {
		log.Fatal("failed to setup task store:", err)
	}

	userStore := NewUserStore(db)
	if err := userStore.migrate(); err != nil {
		log.Fatal("failed to migrate users:", err)
	}

	// setup handlers — pass config
	taskHandler := NewHandler(taskStore)
	authHandler := NewAuthHandler(userStore, cfg)

	// routes
	mux := http.NewServeMux()

	// public
	mux.HandleFunc("POST /auth/register", authHandler.register)
	mux.HandleFunc("POST /auth/login", authHandler.login)

	// protected
	protected := http.NewServeMux()
	protected.HandleFunc("GET /tasks", taskHandler.listTasks)
	protected.HandleFunc("POST /tasks", taskHandler.createTask)
	protected.HandleFunc("GET /tasks/{id}", taskHandler.getTask)
	protected.HandleFunc("PUT /tasks/{id}", taskHandler.updateTask)
	protected.HandleFunc("DELETE /tasks/{id}", taskHandler.deleteTask)

	mux.Handle("/tasks", authMiddleware(protected, cfg))
	mux.Handle("/tasks/", authMiddleware(protected, cfg))

	// server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      loggerMiddleware(mux),
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	go func() {
		fmt.Printf("Server running on http://localhost:%s (env: %s)\n", cfg.Port, cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Forced shutdown:", err)
	}

	log.Println("Server exited cleanly")
}
