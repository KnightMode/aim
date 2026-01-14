package agents

import (
	"ai-task-manager/backend/internal/models"
)

// Agent is the interface that all agents must implement
type Agent interface {
	// Name returns the agent's unique identifier
	Name() string

	// CanHandle returns true if this agent can process tasks with the given tags
	CanHandle(tags []string) bool

	// Execute processes the task and returns an error if it fails
	Execute(task *models.Task) error
}
