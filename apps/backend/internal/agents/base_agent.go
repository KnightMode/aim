package agents

import (
	"ai-task-manager/backend/internal/models"
	"ai-task-manager/backend/internal/websocket"
	"log"
	"time"

	"gorm.io/gorm"
)

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	db    *gorm.DB
	wsHub *websocket.Hub
}

// NewBaseAgent creates a new BaseAgent
func NewBaseAgent(db *gorm.DB, wsHub *websocket.Hub) *BaseAgent {
	return &BaseAgent{
		db:    db,
		wsHub: wsHub,
	}
}

// CreateLog creates an execution log entry
func (a *BaseAgent) CreateLog(task *models.Task, agentName string, level string, message string) {
	executionLog := models.ExecutionLog{
		TaskID:    task.ID,
		AgentName: agentName,
		LogLevel:  level,
		Message:   message,
		CreatedAt: time.Now().UTC(),
	}

	if err := a.db.Create(&executionLog).Error; err != nil {
		log.Printf("Failed to create execution log: %v", err)
		return
	}

	// Broadcast log via WebSocket
	a.wsHub.BroadcastExecutionLog(
		task.ID,
		level,
		message,
		executionLog.CreatedAt.Format(time.RFC3339),
	)

	// Also log to console
	log.Printf("[Task %d] [%s] %s: %s", task.ID, agentName, level, message)
}

// LogInfo logs an info-level message
func (a *BaseAgent) LogInfo(task *models.Task, agentName string, message string) {
	a.CreateLog(task, agentName, models.LogLevelInfo, message)
}

// LogWarning logs a warning-level message
func (a *BaseAgent) LogWarning(task *models.Task, agentName string, message string) {
	a.CreateLog(task, agentName, models.LogLevelWarning, message)
}

// LogError logs an error-level message
func (a *BaseAgent) LogError(task *models.Task, agentName string, message string) {
	a.CreateLog(task, agentName, models.LogLevelError, message)
}

// LogSuccess logs a success-level message
func (a *BaseAgent) LogSuccess(task *models.Task, agentName string, message string) {
	a.CreateLog(task, agentName, models.LogLevelSuccess, message)
}

// UpdateTaskStatus updates the task status in the database
func (a *BaseAgent) UpdateTaskStatus(task *models.Task, status string) error {
	task.Status = status
	return a.db.Save(task).Error
}
