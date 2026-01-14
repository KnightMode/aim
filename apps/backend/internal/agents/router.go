package agents

import (
	"ai-task-manager/backend/internal/models"
	"log"
)

// Router routes tasks to appropriate agents based on tags
type Router struct {
	agents []Agent
}

// NewRouter creates a new Router with registered agents
func NewRouter(agents ...Agent) *Router {
	return &Router{
		agents: agents,
	}
}

// RegisterAgent registers a new agent
func (r *Router) RegisterAgent(agent Agent) {
	r.agents = append(r.agents, agent)
	log.Printf("Agent registered: %s", agent.Name())
}

// GetAgentForTask returns the appropriate agent for a task based on its tags
func (r *Router) GetAgentForTask(task *models.Task) Agent {
	for _, agent := range r.agents {
		if agent.CanHandle(task.Tags) {
			return agent
		}
	}
	return nil
}
