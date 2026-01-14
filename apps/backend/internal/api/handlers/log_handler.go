package handlers

import (
	"ai-task-manager/backend/internal/api/middleware"
	"ai-task-manager/backend/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LogHandler handles execution log requests
type LogHandler struct {
	db *gorm.DB
}

// NewLogHandler creates a new LogHandler
func NewLogHandler(db *gorm.DB) *LogHandler {
	return &LogHandler{db: db}
}

// GetTaskLogs retrieves all logs for a specific task
func (h *LogHandler) GetTaskLogs(c *gin.Context) {
	taskID := c.Param("id")

	var logs []models.ExecutionLog
	if err := h.db.Where("task_id = ?", taskID).Order("created_at ASC").Find(&logs).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to retrieve logs")
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// GetRecentLogs retrieves recent logs across all tasks
func (h *LogHandler) GetRecentLogs(c *gin.Context) {
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	var logs []models.ExecutionLog
	if err := h.db.Order("created_at DESC").Limit(limit).Find(&logs).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to retrieve logs")
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
