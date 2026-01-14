package services

import (
	"ai-task-manager/backend/internal/config"
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
)

// GitHubService handles GitHub API operations
type GitHubService struct {
	client *github.Client
}

// NewGitHubService creates a new GitHubService
func NewGitHubService(cfg *config.GitHubConfig) *GitHubService {
	client := github.NewClient(nil).WithAuthToken(cfg.Token)

	return &GitHubService{
		client: client,
	}
}

// CreatePullRequest creates a pull request on GitHub
func (s *GitHubService) CreatePullRequest(
	repoURL string,
	branch string,
	baseBranch string,
	title string,
	description string,
) (string, error) {
	ctx := context.Background()

	// Parse repo URL to extract owner and repo name
	owner, repoName, err := s.parseRepoURL(repoURL)
	if err != nil {
		return "", err
	}

	// Create PR
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(branch),
		Base:  github.String(baseBranch),
		Body:  github.String(description),
	}

	pr, _, err := s.client.PullRequests.Create(ctx, owner, repoName, newPR)
	if err != nil {
		return "", fmt.Errorf("failed to create pull request: %w", err)
	}

	return pr.GetHTMLURL(), nil
}

// parseRepoURL extracts owner and repo name from a GitHub URL
func (s *GitHubService) parseRepoURL(repoURL string) (string, string, error) {
	// Remove .git suffix if present
	repoURL = strings.TrimSuffix(repoURL, ".git")

	// Handle https://github.com/owner/repo format
	if strings.Contains(repoURL, "github.com/") {
		parts := strings.Split(repoURL, "github.com/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid GitHub repository URL: %s", repoURL)
		}

		pathParts := strings.Split(parts[1], "/")
		if len(pathParts) < 2 {
			return "", "", fmt.Errorf("invalid GitHub repository URL: %s", repoURL)
		}

		return pathParts[0], pathParts[1], nil
	}

	// Handle git@github.com:owner/repo format
	if strings.HasPrefix(repoURL, "git@github.com:") {
		repoURL = strings.TrimPrefix(repoURL, "git@github.com:")
		parts := strings.Split(repoURL, "/")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid GitHub repository URL: %s", repoURL)
		}

		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unsupported repository URL format: %s", repoURL)
}
