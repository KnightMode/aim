package api

import (
	"ai-task-manager/backend/internal/api/handlers"
	"ai-task-manager/backend/internal/api/middleware"
	"ai-task-manager/backend/internal/config"
	"ai-task-manager/backend/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(cfg *config.Config, db *gorm.DB, wsHub *websocket.Hub) *gin.Engine {
	// Set Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Apply global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS(&cfg.Server))

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler(db)
	logHandler := handlers.NewLogHandler(db)
	agentHandler := handlers.NewAgentHandler(db)
	wsHandler := handlers.NewWSHandler(wsHub)

	// WebSocket endpoint
	router.GET("/ws", wsHandler.HandleWebSocket)

	// API routes
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			var queuedCount int64
			var inProgressCount int64
			db.Model(&struct{ ID uint }{}).Table("tasks").Where("status = ?", "queued").Count(&queuedCount)
			db.Model(&struct{ ID uint }{}).Table("tasks").Where("status = ?", "in_progress").Count(&inProgressCount)

			c.JSON(http.StatusOK, gin.H{
				"status":      "ok",
				"queue_size":  queuedCount,
				"in_progress": inProgressCount,
			})
		})

		// Task routes
		api.GET("/tasks", taskHandler.ListTasks)
		api.POST("/tasks", taskHandler.CreateTask)
		api.GET("/tasks/:id", taskHandler.GetTask)
		api.PUT("/tasks/:id", taskHandler.UpdateTask)
		api.DELETE("/tasks/:id", taskHandler.DeleteTask)
		api.PATCH("/tasks/:id/status", taskHandler.UpdateTaskStatus)

		// Log routes
		api.GET("/tasks/:id/logs", logHandler.GetTaskLogs)
		api.GET("/logs/recent", logHandler.GetRecentLogs)

		// Agent routes
		api.GET("/agents", agentHandler.ListAgents)
		api.GET("/agents/stats", agentHandler.GetAgentStats)
	}

	return router
}
