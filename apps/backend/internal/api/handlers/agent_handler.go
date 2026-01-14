package handlers

import (
	"ai-task-manager/backend/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AgentHandler handles agent-related requests
type AgentHandler struct {
	db *gorm.DB
}

// NewAgentHandler creates a new AgentHandler
func NewAgentHandler(db *gorm.DB) *AgentHandler {
	return &AgentHandler{db: db}
}

// AgentInfo represents agent information
type AgentInfo struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Tags    []string `json:"tags"`
}

// AgentStats represents statistics about task processing
type AgentStats struct {
	Queued     int64 `json:"queued"`
	InProgress int64 `json:"in_progress"`
	Completed  int64 `json:"completed"`
	Failed     int64 `json:"failed"`
	Total      int64 `json:"total"`
}

// ListAgents returns information about all available agents
func (h *AgentHandler) ListAgents(c *gin.Context) {
	agents := []AgentInfo{
		{
			Name:    "coding_agent",
			Enabled: true,
			Tags:    []string{models.TagCoding},
		},
		{
			Name:    "docs_agent",
			Enabled: true,
			Tags:    []string{models.TagDocumentation},
		},
	}

	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

// GetAgentStats returns statistics about task processing
func (h *AgentHandler) GetAgentStats(c *gin.Context) {
	var stats AgentStats

	// Count tasks by status
	h.db.Model(&models.Task{}).Where("status = ?", models.StatusQueued).Count(&stats.Queued)
	h.db.Model(&models.Task{}).Where("status = ?", models.StatusInProgress).Count(&stats.InProgress)
	h.db.Model(&models.Task{}).Where("status = ?", models.StatusCompleted).Count(&stats.Completed)
	h.db.Model(&models.Task{}).Where("status = ?", models.StatusFailed).Count(&stats.Failed)
	h.db.Model(&models.Task{}).Count(&stats.Total)

	c.JSON(http.StatusOK, stats)
}
