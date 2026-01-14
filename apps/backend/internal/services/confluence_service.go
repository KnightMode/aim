package services

import (
	"ai-task-manager/backend/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ConfluenceService handles Confluence API operations
type ConfluenceService struct {
	baseURL  string
	username string
	apiToken string
	client   *http.Client
}

// NewConfluenceService creates a new ConfluenceService
func NewConfluenceService(cfg *config.ConfluenceConfig) *ConfluenceService {
	return &ConfluenceService{
		baseURL:  strings.TrimSuffix(cfg.URL, "/"),
		username: cfg.Username,
		apiToken: cfg.APIToken,
		client:   &http.Client{},
	}
}

// ConfluencePage represents a Confluence page
type ConfluencePage struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Title   string `json:"title"`
	Body    ConfluenceBody `json:"body"`
	Space   ConfluenceSpace `json:"space"`
	Version ConfluenceVersion `json:"version,omitempty"`
}

type ConfluenceBody struct {
	Storage ConfluenceStorage `json:"storage"`
}

type ConfluenceStorage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}

type ConfluenceSpace struct {
	Key string `json:"key"`
}

type ConfluenceVersion struct {
	Number int `json:"number"`
}

// SearchPageByTitle searches for a page by title in a space
func (s *ConfluenceService) SearchPageByTitle(spaceKey string, title string) (*ConfluencePage, error) {
	url := fmt.Sprintf("%s/rest/api/content?spaceKey=%s&title=%s&expand=body.storage,version",
		s.baseURL, spaceKey, title)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Results []ConfluencePage `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, nil // Page not found
	}

	return &result.Results[0], nil
}

// CreatePage creates a new Confluence page
func (s *ConfluenceService) CreatePage(spaceKey string, title string, content string) (string, error) {
	// Convert markdown to Confluence storage format
	storageContent := s.markdownToConfluence(content)

	page := ConfluencePage{
		Type:   "page",
		Title:  title,
		Space:  ConfluenceSpace{Key: spaceKey},
		Body: ConfluenceBody{
			Storage: ConfluenceStorage{
				Value:          storageContent,
				Representation: "storage",
			},
		},
	}

	jsonData, err := json.Marshal(page)
	if err != nil {
		return "", fmt.Errorf("failed to marshal page: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/content", s.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.apiToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdPage ConfluencePage
	if err := json.NewDecoder(resp.Body).Decode(&createdPage); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	pageURL := fmt.Sprintf("%s/pages/viewpage.action?pageId=%s", s.baseURL, createdPage.ID)
	return pageURL, nil
}

// UpdatePage updates an existing Confluence page
func (s *ConfluenceService) UpdatePage(pageID string, title string, content string, version int) (string, error) {
	// Convert markdown to Confluence storage format
	storageContent := s.markdownToConfluence(content)

	page := ConfluencePage{
		ID:      pageID,
		Type:    "page",
		Title:   title,
		Version: ConfluenceVersion{Number: version + 1},
		Body: ConfluenceBody{
			Storage: ConfluenceStorage{
				Value:          storageContent,
				Representation: "storage",
			},
		},
	}

	jsonData, err := json.Marshal(page)
	if err != nil {
		return "", fmt.Errorf("failed to marshal page: %w", err)
	}

	url := fmt.Sprintf("%s/rest/api/content/%s", s.baseURL, pageID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(s.username, s.apiToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to update page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("update failed with status %d: %s", resp.StatusCode, string(body))
	}

	pageURL := fmt.Sprintf("%s/pages/viewpage.action?pageId=%s", s.baseURL, pageID)
	return pageURL, nil
}

// markdownToConfluence converts markdown to Confluence storage format
func (s *ConfluenceService) markdownToConfluence(markdown string) string {
	// Basic markdown to HTML conversion for Confluence
	content := markdown

	// Headers
	content = strings.ReplaceAll(content, "### ", "<h3>")
	content = strings.ReplaceAll(content, "## ", "<h2>")
	content = strings.ReplaceAll(content, "# ", "<h1>")

	// Bold and italic
	content = strings.ReplaceAll(content, "**", "<strong>")
	content = strings.ReplaceAll(content, "*", "<em>")

	// Code blocks (simplified)
	content = strings.ReplaceAll(content, "`", "<code>")

	// Convert newlines to <br/> or <p> tags
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = "<p>" + line + "</p>"
		}
	}
	content = strings.Join(lines, "")

	return content
}
