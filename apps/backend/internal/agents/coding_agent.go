package agents

import (
	"ai-task-manager/backend/internal/config"
	"ai-task-manager/backend/internal/models"
	"ai-task-manager/backend/internal/services"
	"ai-task-manager/backend/internal/websocket"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// CodingAgent handles coding tasks
type CodingAgent struct {
	*BaseAgent
	gitService    *services.GitService
	githubService *services.GitHubService
	llmService    *services.LLMService
	tempDir       string
}

// NewCodingAgent creates a new CodingAgent
func NewCodingAgent(
	db *gorm.DB,
	wsHub *websocket.Hub,
	cfg *config.Config,
) *CodingAgent {
	gitService := services.NewGitService(&cfg.Git, cfg.GitHub.Token)
	githubService := services.NewGitHubService(&cfg.GitHub)
	llmService := services.NewLLMService(&cfg.Anthropic)

	return &CodingAgent{
		BaseAgent:     NewBaseAgent(db, wsHub),
		gitService:    gitService,
		githubService: githubService,
		llmService:    llmService,
		tempDir:       "/tmp/ai-task-manager",
	}
}

// Name returns the agent's name
func (a *CodingAgent) Name() string {
	return "coding_agent"
}

// CanHandle checks if this agent can handle the given tags
func (a *CodingAgent) CanHandle(tags []string) bool {
	for _, tag := range tags {
		if tag == models.TagCoding {
			return true
		}
	}
	return false
}

// Execute processes a coding task
func (a *CodingAgent) Execute(task *models.Task) error {
	agentName := a.Name()

	a.LogInfo(task, agentName, fmt.Sprintf("Starting coding task: %s", task.Title))

	// Extract metadata
	repoURL := task.Metadata.GetString("repo_url")
	branch := task.Metadata.GetString("branch")
	if branch == "" {
		branch = "main"
	}

	// Use task description as the instruction
	instruction := task.Description
	if instruction == "" {
		instruction = task.Title // Fallback to title if no description
	}

	// Validate metadata
	if repoURL == "" {
		a.LogError(task, agentName, "Missing repo_url in metadata")
		return fmt.Errorf("missing repo_url in metadata")
	}

	// Create temp directory for this task
	taskDir := filepath.Join(a.tempDir, fmt.Sprintf("task-%d", task.ID))
	defer os.RemoveAll(taskDir) // Clean up after execution

	// Clone repository
	a.LogInfo(task, agentName, fmt.Sprintf("Cloning repository from %s...", repoURL))
	repo, err := a.gitService.CloneRepository(repoURL, taskDir)
	if err != nil {
		a.LogError(task, agentName, fmt.Sprintf("Failed to clone repository: %v", err))
		return fmt.Errorf("git clone failed: %w", err)
	}
	a.LogInfo(task, agentName, "Repository cloned successfully")

	// Create feature branch
	branchName := fmt.Sprintf("ai-task-%d", task.ID)
	a.LogInfo(task, agentName, fmt.Sprintf("Creating branch: %s", branchName))
	err = a.gitService.CreateBranch(repo, branchName)
	if err != nil {
		a.LogError(task, agentName, fmt.Sprintf("Failed to create branch: %v", err))
		return fmt.Errorf("failed to create branch: %w", err)
	}
	a.LogInfo(task, agentName, "Branch created successfully")

	// Analyze repository and determine files to modify
	a.LogInfo(task, agentName, "Analyzing repository structure...")
	repoStructure, err := a.gitService.GetRepositoryStructure(taskDir)
	if err != nil {
		a.LogError(task, agentName, fmt.Sprintf("Failed to analyze repository: %v", err))
		return fmt.Errorf("repository analysis failed: %w", err)
	}
	a.LogInfo(task, agentName, fmt.Sprintf("Found %d files in repository", len(repoStructure)))

	// Ask AI to determine which files need modification
	a.LogInfo(task, agentName, "Consulting AI to determine which files to modify...")
	filesToModify, err := a.llmService.DetermineFilesToModify(instruction, repoStructure)
	if err != nil {
		a.LogError(task, agentName, fmt.Sprintf("Failed to determine files: %v", err))
		return fmt.Errorf("file determination failed: %w", err)
	}

	if len(filesToModify) == 0 {
		a.LogWarning(task, agentName, "AI determined no files need modification")
		task.Result = "No files needed modification based on the instruction"
		return nil
	}

	a.LogInfo(task, agentName, fmt.Sprintf("AI identified %d file(s) to modify: %v", len(filesToModify), filesToModify))

	// Process each file
	for _, filePath := range filesToModify {
		a.LogInfo(task, agentName, fmt.Sprintf("Processing file: %s", filePath))

		// Read existing code
		a.LogInfo(task, agentName, fmt.Sprintf("Reading existing code from %s...", filePath))
		existingCode, err := a.gitService.ReadFile(taskDir, filePath)
		if err != nil {
			a.LogWarning(task, agentName, fmt.Sprintf("File not found, will create new file: %s", filePath))
			existingCode = ""
		}

		// Generate code changes with AI
		a.LogInfo(task, agentName, "Generating code changes with Claude AI...")
		modifiedCode, err := a.llmService.GenerateCode(instruction, existingCode, filePath)
		if err != nil {
			a.LogError(task, agentName, fmt.Sprintf("AI generation failed: %v", err))
			return fmt.Errorf("AI generation failed: %w", err)
		}
		a.LogInfo(task, agentName, "AI generation complete")

		// Write modified code
		a.LogInfo(task, agentName, fmt.Sprintf("Applying changes to %s...", filePath))
		err = a.gitService.WriteFile(taskDir, filePath, modifiedCode)
		if err != nil {
			a.LogError(task, agentName, fmt.Sprintf("Failed to write file: %v", err))
			return fmt.Errorf("failed to write file: %w", err)
		}
		a.LogInfo(task, agentName, "Changes applied successfully")
	}

	// Commit changes
	commitMessage := fmt.Sprintf("%s\n\n%s\n\nGenerated by AI Task Manager (Task #%d)",
		task.Title, task.Description, task.ID)
	a.LogInfo(task, agentName, "Committing changes...")
	commitHash, err := a.gitService.CommitChanges(repo, commitMessage)
	if err != nil {
		a.LogError(task, agentName, fmt.Sprintf("Failed to commit: %v", err))
		return fmt.Errorf("failed to commit: %w", err)
	}
	a.LogInfo(task, agentName, fmt.Sprintf("Commit created: %s", commitHash[:8]))

	// Push to remote
	a.LogInfo(task, agentName, "Pushing to remote repository...")
	err = a.gitService.PushToRemote(repo)
	if err != nil {
		a.LogError(task, agentName, fmt.Sprintf("Failed to push: %v", err))
		return fmt.Errorf("git push failed: %w", err)
	}
	a.LogInfo(task, agentName, "Push successful")

	// Create pull request
	a.LogInfo(task, agentName, "Creating pull request on GitHub...")
	prDescription := fmt.Sprintf("%s\n\n---\n**Task ID:** %d\n**Instruction:** %s\n\nThis PR was automatically generated by AI Task Manager.",
		task.Description, task.ID, instruction)

	prURL, err := a.githubService.CreatePullRequest(repoURL, branchName, branch, task.Title, prDescription)
	if err != nil {
		a.LogWarning(task, agentName, fmt.Sprintf("PR creation failed, but code is pushed: %v", err))
		task.Result = "Code pushed successfully. Manual PR creation needed."
		return nil // Still mark as success since code is pushed
	}

	a.LogSuccess(task, agentName, fmt.Sprintf("✓ Pull request created: %s", prURL))
	task.Result = fmt.Sprintf("Pull request created: %s", prURL)

	return nil
}
