package agents

import (
	"ai-task-manager/backend/internal/config"
	"ai-task-manager/backend/internal/models"
	"ai-task-manager/backend/internal/services"
	"ai-task-manager/backend/internal/websocket"
	"fmt"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
)

// DocsAgent handles documentation tasks
type DocsAgent struct {
	*BaseAgent
	confluenceService *services.ConfluenceService
	llmService        *services.LLMService
}

// NewDocsAgent creates a new DocsAgent
func NewDocsAgent(
	db *gorm.DB,
	wsHub *websocket.Hub,
	cfg *config.Config,
) *DocsAgent {
	confluenceService := services.NewConfluenceService(&cfg.Confluence)
	llmService := services.NewLLMService(&cfg.Anthropic)

	return &DocsAgent{
		BaseAgent:         NewBaseAgent(db, wsHub),
		confluenceService: confluenceService,
		llmService:        llmService,
	}
}

// Name returns the agent's name
func (a *DocsAgent) Name() string {
	return "docs_agent"
}

// CanHandle checks if this agent can handle the given tags
func (a *DocsAgent) CanHandle(tags []string) bool {
	for _, tag := range tags {
		if tag == models.TagDocumentation {
			return true
		}
	}
	return false
}

// Execute processes a documentation task
func (a *DocsAgent) Execute(task *models.Task) error {
	agentName := a.Name()

	a.LogInfo(task, agentName, fmt.Sprintf("Starting documentation task: %s", task.Title))

	// Extract metadata
	confluenceSpace := task.Metadata.GetString("confluence_space")
	pageTitle := task.Metadata.GetString("page_title")
	contentSource := task.Metadata.GetString("content_source")
	sourceValue := task.Metadata.GetString("source_value")

	// Use task description as the instruction
	instruction := task.Description
	if instruction == "" {
		instruction = task.Title // Fallback to title if no description
	}

	// Validate metadata
	if confluenceSpace == "" {
		a.LogError(task, agentName, "Missing confluence_space in metadata")
		return fmt.Errorf("missing confluence_space in metadata")
	}
	if pageTitle == "" {
		a.LogError(task, agentName, "Missing page_title in metadata")
		return fmt.Errorf("missing page_title in metadata")
	}
	if contentSource == "" {
		contentSource = "text"
	}

	// Get source content based on content_source type
	var sourceContent string
	var err error

	switch contentSource {
	case "file":
		a.LogInfo(task, agentName, fmt.Sprintf("Reading source content from file: %s", sourceValue))
		sourceContent, err = a.readFromFile(sourceValue)
		if err != nil {
			a.LogError(task, agentName, fmt.Sprintf("Failed to read file: %v", err))
			return fmt.Errorf("failed to read file: %w", err)
		}

	case "url":
		a.LogInfo(task, agentName, fmt.Sprintf("Fetching source content from URL: %s", sourceValue))
		sourceContent, err = a.fetchFromURL(sourceValue)
		if err != nil {
			a.LogError(task, agentName, fmt.Sprintf("Failed to fetch URL: %v", err))
			return fmt.Errorf("failed to fetch URL: %w", err)
		}

	case "text":
		a.LogInfo(task, agentName, "Using provided text as source content")
		sourceContent = sourceValue

	default:
		a.LogError(task, agentName, fmt.Sprintf("Invalid content_source: %s", contentSource))
		return fmt.Errorf("invalid content_source: %s", contentSource)
	}

	a.LogInfo(task, agentName, "Content loaded successfully")

	// Generate documentation with AI
	a.LogInfo(task, agentName, "Generating documentation with Claude AI...")
	documentation, err := a.llmService.GenerateDocumentation(instruction, sourceContent)
	if err != nil {
		a.LogError(task, agentName, fmt.Sprintf("AI generation failed: %v", err))
		return fmt.Errorf("AI generation failed: %w", err)
	}
	a.LogInfo(task, agentName, "AI generation complete")

	// Search for existing page
	a.LogInfo(task, agentName, fmt.Sprintf("Searching for existing page: %s", pageTitle))
	existingPage, err := a.confluenceService.SearchPageByTitle(confluenceSpace, pageTitle)
	if err != nil {
		a.LogWarning(task, agentName, fmt.Sprintf("Search failed: %v", err))
	}

	var pageURL string

	if existingPage != nil {
		// Update existing page
		a.LogInfo(task, agentName, fmt.Sprintf("Updating existing page (ID: %s)...", existingPage.ID))
		pageURL, err = a.confluenceService.UpdatePage(
			existingPage.ID,
			pageTitle,
			documentation,
			existingPage.Version.Number,
		)
		if err != nil {
			a.LogError(task, agentName, fmt.Sprintf("Failed to update page: %v", err))
			return fmt.Errorf("failed to update page: %w", err)
		}
		a.LogInfo(task, agentName, "Page updated successfully")
	} else {
		// Create new page
		a.LogInfo(task, agentName, fmt.Sprintf("Creating new page in space %s...", confluenceSpace))
		pageURL, err = a.confluenceService.CreatePage(confluenceSpace, pageTitle, documentation)
		if err != nil {
			a.LogError(task, agentName, fmt.Sprintf("Failed to create page: %v", err))
			return fmt.Errorf("failed to create page: %w", err)
		}
		a.LogInfo(task, agentName, "Page created successfully")
	}

	a.LogSuccess(task, agentName, fmt.Sprintf("✓ Documentation published: %s", pageURL))
	task.Result = fmt.Sprintf("Documentation published: %s", pageURL)

	return nil
}

// readFromFile reads content from a file
func (a *DocsAgent) readFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// fetchFromURL fetches content from a URL
func (a *DocsAgent) fetchFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
