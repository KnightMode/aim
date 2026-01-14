package handlers

import (
	"ai-task-manager/backend/internal/api/middleware"
	"ai-task-manager/backend/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TaskHandler handles task-related requests
type TaskHandler struct {
	db *gorm.DB
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Title       string                 `json:"title" binding:"required"`
	Description string                 `json:"description"`
	Tags        []string               `json:"tags" binding:"required"`
	Priority    int                    `json:"priority"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UpdateStatusRequest represents the request body for updating task status
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// ListTasks retrieves all tasks with optional filters
func (h *TaskHandler) ListTasks(c *gin.Context) {
	var tasks []models.Task

	query := h.db.Model(&models.Task{})

	// Apply filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if tag := c.Query("tag"); tag != "" {
		query = query.Where("JSON_CONTAINS(tags, ?)", "\""+tag+"\"")
	}

	// Apply limit
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Order by priority (desc) and created_at (asc)
	query = query.Order("priority DESC, created_at ASC").Limit(limit)

	if err := query.Find(&tasks).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to retrieve tasks")
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// CreateTask creates a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	// Validate tags
	if len(req.Tags) == 0 {
		middleware.JSONError(c, http.StatusBadRequest, "invalid_tags", "At least one tag is required")
		return
	}

	// Create task
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		Priority:    req.Priority,
		Status:      models.StatusTodo,
		Metadata:    req.Metadata,
	}

	if err := h.db.Create(&task).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to create task")
		return
	}

	// Auto-queue the task
	task.Status = models.StatusQueued
	if err := h.db.Save(&task).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to queue task")
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTask retrieves a single task by ID
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.JSONError(c, http.StatusNotFound, "not_found", "Task not found")
			return
		}
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to retrieve task")
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates a task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.JSONError(c, http.StatusNotFound, "not_found", "Task not found")
			return
		}
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to retrieve task")
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	// Update fields
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Metadata != nil {
		task.Metadata = req.Metadata
	}

	if err := h.db.Save(&task).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to update task")
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.JSONError(c, http.StatusNotFound, "not_found", "Task not found")
			return
		}
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to retrieve task")
		return
	}

	// Delete associated logs first (cascade)
	if err := h.db.Where("task_id = ?", id).Delete(&models.ExecutionLog{}).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to delete task logs")
		return
	}

	// Delete task
	if err := h.db.Delete(&task).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to delete task")
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateTaskStatus updates the status of a task
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	if err := h.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			middleware.JSONError(c, http.StatusNotFound, "not_found", "Task not found")
			return
		}
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to retrieve task")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		models.StatusTodo:       true,
		models.StatusQueued:     true,
		models.StatusInProgress: true,
		models.StatusCompleted:  true,
		models.StatusFailed:     true,
	}

	if !validStatuses[req.Status] {
		middleware.JSONError(c, http.StatusBadRequest, "invalid_status", "Invalid status value")
		return
	}

	task.Status = req.Status
	if err := h.db.Save(&task).Error; err != nil {
		middleware.JSONError(c, http.StatusInternalServerError, "database_error", "Failed to update task status")
		return
	}

	c.JSON(http.StatusOK, task)
}
