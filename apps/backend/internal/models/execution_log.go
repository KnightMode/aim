package models

import (
	"time"
)

// ExecutionLog represents a log entry for task execution
type ExecutionLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	TaskID     uint      `json:"task_id" gorm:"not null;index"`
	AgentName  string    `json:"agent_name" gorm:"size:100"`
	LogLevel   string    `json:"log_level" gorm:"not null;size:20;index"`
	Message    string    `json:"message" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
}

// TableName specifies the table name for the ExecutionLog model
func (ExecutionLog) TableName() string {
	return "execution_logs"
}

// Log level constants
const (
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelError   = "error"
	LogLevelSuccess = "success"
)
