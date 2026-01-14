package main

import (
	"ai-task-manager/backend/internal/agents"
	"ai-task-manager/backend/internal/api"
	"ai-task-manager/backend/internal/config"
	"ai-task-manager/backend/internal/database"
	"ai-task-manager/backend/internal/queue"
	"ai-task-manager/backend/internal/websocket"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting AI Task Manager Backend...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	// Connect to database
	err = database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	err = database.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()
	log.Println("WebSocket hub started")

	// Initialize agents
	db := database.GetDB()
	codingAgent := agents.NewCodingAgent(db, wsHub, cfg)
	docsAgent := agents.NewDocsAgent(db, wsHub, cfg)

	// Create agent router
	agentRouter := agents.NewRouter(codingAgent, docsAgent)
	log.Println("Agents registered")

	// Initialize task queue
	taskQueue := queue.NewTaskQueue(
		cfg.Queue.WorkerCount,
		cfg.Queue.PollInterval,
		cfg.Queue.BufferSize,
		db,
		wsHub,
		agentRouter,
	)
	taskQueue.Start()
	defer taskQueue.Stop()

	// Setup HTTP router
	router := api.SetupRouter(cfg, db, wsHub)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("Frontend URL: %s", cfg.Server.FrontendURL)
	log.Printf("API: http://localhost:%s/api", cfg.Server.Port)
	log.Printf("WebSocket: ws://localhost:%s/ws", cfg.Server.Port)

	// Handle graceful shutdown
	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	taskQueue.Stop()
	log.Println("Server stopped gracefully")
}
