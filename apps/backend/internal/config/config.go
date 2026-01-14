package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Queue     QueueConfig
	Anthropic AnthropicConfig
	Git       GitConfig
	GitHub    GitHubConfig
	Confluence ConfluenceConfig
	Agent     AgentConfig
	Logging   LoggingConfig
}

type ServerConfig struct {
	Port        string
	Host        string
	FrontendURL string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Charset  string
}

type QueueConfig struct {
	WorkerCount    int
	PollInterval   time.Duration
	BufferSize     int
}

type AnthropicConfig struct {
	APIKey    string
	Model     string
	MaxTokens int
}

type GitConfig struct {
	Username string
	Email    string
}

type GitHubConfig struct {
	Token string
}

type ConfluenceConfig struct {
	URL      string
	Username string
	APIToken string
}

type AgentConfig struct {
	CodingAgentEnabled bool
	DocsAgentEnabled   bool
	AutoExecuteTasks   bool
}

type LoggingConfig struct {
	Level  string
	Format string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnv("SERVER_PORT", "8080"),
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Name:     getEnv("DB_NAME", "ai_task_manager"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		Queue: QueueConfig{
			WorkerCount:  getEnvAsInt("WORKER_COUNT", 3),
			PollInterval: getEnvAsDuration("QUEUE_POLL_INTERVAL", 2*time.Second),
			BufferSize:   getEnvAsInt("QUEUE_BUFFER_SIZE", 100),
		},
		Anthropic: AnthropicConfig{
			APIKey:    getEnv("ANTHROPIC_API_KEY", ""),
			Model:     getEnv("ANTHROPIC_MODEL", "claude-3-5-sonnet-20241022"),
			MaxTokens: getEnvAsInt("ANTHROPIC_MAX_TOKENS", 4096),
		},
		Git: GitConfig{
			Username: getEnv("GIT_USERNAME", ""),
			Email:    getEnv("GIT_EMAIL", ""),
		},
		GitHub: GitHubConfig{
			Token: getEnv("GITHUB_TOKEN", ""),
		},
		Confluence: ConfluenceConfig{
			URL:      getEnv("CONFLUENCE_URL", ""),
			Username: getEnv("CONFLUENCE_USERNAME", ""),
			APIToken: getEnv("CONFLUENCE_API_TOKEN", ""),
		},
		Agent: AgentConfig{
			CodingAgentEnabled: getEnvAsBool("CODING_AGENT_ENABLED", true),
			DocsAgentEnabled:   getEnvAsBool("DOCS_AGENT_ENABLED", true),
			AutoExecuteTasks:   getEnvAsBool("AUTO_EXECUTE_TASKS", true),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validate required config
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if required configuration values are present
func (c *Config) Validate() error {
	// if c.Database.Password == "" {
	// 	return fmt.Errorf("DB_PASSWORD is required")
	// }
	if c.Anthropic.APIKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY is required")
	}
	if c.Git.Username == "" {
		return fmt.Errorf("GIT_USERNAME is required")
	}
	if c.Git.Email == "" {
		return fmt.Errorf("GIT_EMAIL is required")
	}
	if c.GitHub.Token == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}
	return nil
}

// DSN returns the MySQL Data Source Name
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.Charset,
	)
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
