package services

import (
	"ai-task-manager/backend/internal/config"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitService handles Git operations
type GitService struct {
	username string
	email    string
	token    string
}

// NewGitService creates a new GitService
func NewGitService(cfg *config.GitConfig, token string) *GitService {
	return &GitService{
		username: cfg.Username,
		email:    cfg.Email,
		token:    token,
	}
}

// CloneRepository clones a Git repository to the specified directory
func (s *GitService) CloneRepository(repoURL string, destPath string) (*git.Repository, error) {
	// Create auth if token is available
	auth := &http.BasicAuth{
		Username: s.username,
		Password: s.token,
	}

	repo, err := git.PlainClone(destPath, false, &git.CloneOptions{
		URL:      repoURL,
		Auth:     auth,
		Progress: os.Stdout,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	return repo, nil
}

// CreateBranch creates a new branch in the repository
func (s *GitService) CreateBranch(repo *git.Repository, branchName string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get HEAD reference
	headRef, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Create new branch
	refName := fmt.Sprintf("refs/heads/%s", branchName)
	ref := plumbing.NewHashReference(plumbing.ReferenceName(refName), headRef.Hash())

	err = repo.Storer.SetReference(ref)
	if err != nil {
		return fmt.Errorf("failed to create branch reference: %w", err)
	}

	// Checkout the new branch
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(refName),
	})

	if err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	return nil
}

// CommitChanges commits all changes in the repository
func (s *GitService) CommitChanges(repo *git.Repository, message string) (string, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	// Add all changes
	_, err = worktree.Add(".")
	if err != nil {
		return "", fmt.Errorf("failed to add changes: %w", err)
	}

	// Commit
	commit, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  s.username,
			Email: s.email,
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to commit: %w", err)
	}

	return commit.String(), nil
}

// PushToRemote pushes the current branch to the remote repository
func (s *GitService) PushToRemote(repo *git.Repository) error {
	auth := &http.BasicAuth{
		Username: s.username,
		Password: s.token,
	}

	err := repo.Push(&git.PushOptions{
		Auth:     auth,
		Progress: os.Stdout,
	})

	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// WriteFile writes content to a file in the repository
func (s *GitService) WriteFile(repoPath string, filePath string, content string) error {
	fullPath := filepath.Join(repoPath, filePath)

	// Create directories if they don't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Write file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ReadFile reads content from a file in the repository
func (s *GitService) ReadFile(repoPath string, filePath string) (string, error) {
	fullPath := filepath.Join(repoPath, filePath)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

// GetRepositoryStructure returns a list of all relevant files in the repository
func (s *GitService) GetRepositoryStructure(repoPath string) ([]string, error) {
	var files []string

	// Directories to skip
	skipDirs := map[string]bool{
		".git":         true,
		"node_modules": true,
		"vendor":       true,
		"dist":         true,
		"build":        true,
		".next":        true,
		"target":       true,
		"__pycache__":  true,
		".idea":        true,
		".vscode":      true,
	}

	// File extensions to skip
	skipExtensions := map[string]bool{
		".lock":     true,
		".log":      true,
		".png":      true,
		".jpg":      true,
		".jpeg":     true,
		".gif":      true,
		".ico":      true,
		".svg":      true,
		".woff":     true,
		".woff2":    true,
		".ttf":      true,
		".eot":      true,
		".min.js":   true,
		".min.css":  true,
		".map":      true,
		".DS_Store": true,
	}

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories we don't care about
		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip files with certain extensions
		ext := filepath.Ext(info.Name())
		if skipExtensions[ext] || skipExtensions[info.Name()] {
			return nil
		}

		// Get relative path from repo root
		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return err
		}

		files = append(files, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk repository: %w", err)
	}

	return files, nil
}
