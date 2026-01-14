package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Task represents a task in the system
type Task struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Title         string         `json:"title" gorm:"not null;size:255"`
	Description   string         `json:"description" gorm:"type:text"`
	Tags          StringArray    `json:"tags" gorm:"type:json;not null"`
	Status        string         `json:"status" gorm:"not null;default:todo;size:50;index"`
	Priority      int            `json:"priority" gorm:"default:0"`
	AssignedAgent string         `json:"assigned_agent" gorm:"size:100"`
	Result        string         `json:"result" gorm:"type:text"`
	ErrorMsg      string         `json:"error_msg" gorm:"type:text;column:error_msg"`
	Metadata      JSON           `json:"metadata" gorm:"type:json"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	StartedAt     *time.Time     `json:"started_at"`
	CompletedAt   *time.Time     `json:"completed_at"`
}

// TableName specifies the table name for the Task model
func (Task) TableName() string {
	return "tasks"
}

// StringArray is a custom type for storing []string in JSON
type StringArray []string

// Value implements the driver.Valuer interface
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// JSON is a custom type for storing map[string]interface{} in JSON
type JSON map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// GetString safely retrieves a string value from JSON metadata
func (j JSON) GetString(key string) string {
	if val, ok := j[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetStringSlice safely retrieves a string slice from JSON metadata
func (j JSON) GetStringSlice(key string) []string {
	if val, ok := j[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return []string{}
}

// Task status constants
const (
	StatusTodo       = "todo"
	StatusQueued     = "queued"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// Task tag constants
const (
	TagCoding         = "coding"
	TagDocumentation  = "documentation"
)
