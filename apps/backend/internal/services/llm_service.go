package services

import (
	"ai-task-manager/backend/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/liushuangls/go-anthropic/v2"
)

// LLMService handles interactions with the Claude API
type LLMService struct {
	client    *anthropic.Client
	model     string
	maxTokens int
}

// NewLLMService creates a new LLMService
func NewLLMService(cfg *config.AnthropicConfig) *LLMService {
	client := anthropic.NewClient(cfg.APIKey)

	return &LLMService{
		client:    client,
		model:     cfg.Model,
		maxTokens: cfg.MaxTokens,
	}
}

// GenerateCode generates code based on the given prompt
func (s *LLMService) GenerateCode(instruction string, existingCode string, filePath string) (string, error) {
	prompt := fmt.Sprintf(`You are a code generation assistant. Your task is to modify the existing code according to the given instruction.

File: %s

Existing Code:
%s

Instruction: %s

Please provide the complete modified code. Return ONLY the code without any explanations or markdown formatting.`,
		filePath, existingCode, instruction)

	response, err := s.generateText(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}

	return response, nil
}

// GenerateDocumentation generates documentation based on the given content
func (s *LLMService) GenerateDocumentation(instruction string, sourceContent string) (string, error) {
	prompt := fmt.Sprintf(`You are a technical documentation writer. Generate comprehensive documentation based on the following content.

Instruction: %s

Source Content:
%s

Please provide well-structured documentation in markdown format.`,
		instruction, sourceContent)

	response, err := s.generateText(prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate documentation: %w", err)
	}

	return response, nil
}

// DetermineFilesToModify analyzes the repository and determines which files need modification
func (s *LLMService) DetermineFilesToModify(instruction string, repoFiles []string) ([]string, error) {
	// Limit the file list if it's too large (keep first 500 files)
	displayFiles := repoFiles
	if len(repoFiles) > 500 {
		displayFiles = repoFiles[:500]
	}

	prompt := fmt.Sprintf(`You are a code analysis assistant. Based on the task instruction and the repository structure, determine which files need to be modified.

Task Instruction:
%s

Repository Files:
%s

Analyze the task and the repository structure. Return ONLY a JSON array of file paths that need to be modified.
Consider:
- Which files are most relevant to the task
- Dependencies between files
- Common patterns (e.g., if modifying a component, you might need its types file)

Return format: ["path/to/file1.ext", "path/to/file2.ext"]
Return ONLY the JSON array, no explanations.`, instruction, strings.Join(displayFiles, "\n"))

	response, err := s.generateText(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to determine files: %w", err)
	}

	// Clean up the response - remove markdown code blocks if present
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Parse JSON array
	var files []string
	err = json.Unmarshal([]byte(response), &files)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file list: %w (response: %s)", err, response)
	}

	return files, nil
}

// generateText sends a prompt to Claude and returns the response
func (s *LLMService) generateText(prompt string) (string, error) {
	ctx := context.Background()

	resp, err := s.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: anthropic.Model(s.model),
		Messages: []anthropic.Message{
			{
				Role: "user",
				Content: []anthropic.MessageContent{
					anthropic.NewTextMessageContent(prompt),
				},
			},
		},
		MaxTokens: s.maxTokens,
	})

	if err != nil {
		return "", fmt.Errorf("Claude API error: %w", err)
	}

	if len(resp.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude API")
	}

	// Extract text from content blocks
	var result string
	for _, c := range resp.Content {
		// Access the Text field directly (it's a pointer)
		if c.Text != nil {
			result += *c.Text
		}
	}

	if result == "" {
		return "", fmt.Errorf("no text content in Claude API response")
	}

	return result, nil
}
