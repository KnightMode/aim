# AI Task Manager - Complete Implementation Plan

## 📋 Project Overview

An AI-powered task management system with a Jira-like interface where tasks are automatically executed by intelligent agents based on their tags (coding, documentation, etc.). Tasks are queued and processed automatically with real-time status updates.

### Key Features
- **Jira-like Kanban Board**: Drag-and-drop interface with multiple status columns
- **Intelligent Agents**: Auto-execute tasks based on tags (coding, documentation)
- **Real-time Updates**: Live status changes and execution logs via WebSocket
- **Auto-execution**: Tasks execute immediately upon creation
- **GitHub Integration**: Coding agent creates PRs for code changes
- **Confluence Integration**: Documentation agent publishes to Confluence

### User Requirements Summary
- ✅ Frontend: React with TypeScript
- ✅ Backend: Golang
- ✅ Database: MySQL (local installation)
- ✅ Task types: Coding & Documentation
- ✅ Agents: Auto-execute with external repo access
- ✅ Always create PRs for review
- ✅ User provides API keys
- ✅ Credentials in .env file
- ✅ Personal use (single user for now)

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     React Frontend                          │
│                    (localhost:5173)                         │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Kanban Board│  │ Task Details │  │ Logs Viewer  │      │
│  │  (Drag/Drop)│  │    Modal     │  │  (Real-time) │      │
│  └─────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────────┬────────────────────────────────────┘
                         │ REST API + WebSocket
┌────────────────────────┴────────────────────────────────────┐
│                    Golang Backend                           │
│                    (localhost:8080)                         │
│  ┌──────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │ API Layer│  │ Task Queue   │  │ Agent Router │         │
│  │  (Gin)   │  │  (Channels)  │  │  (Tag-based) │         │
│  └──────────┘  └──────────────┘  └──────────────┘         │
│                                                             │
│  ┌──────────────────────────────────────────────┐         │
│  │         Agent Execution Engine                │         │
│  │  ┌──────────────┐    ┌──────────────────┐   │         │
│  │  │ Coding Agent │    │ Documentation    │   │         │
│  │  │ - Git Clone  │    │ Agent            │   │         │
│  │  │ - AI Code Gen│    │ - Confluence API │   │         │
│  │  │ - GitHub PR  │    │ - AI Doc Gen     │   │         │
│  │  └──────────────┘    └──────────────────┘   │         │
│  └──────────────────────────────────────────────┘         │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────┴────────────────────────────────────┐
│                MySQL Database (Local)                       │
│     Tasks | Execution Logs | Agent Configurations          │
└─────────────────────────────────────────────────────────────┘
```

---

## 🛠️ Technology Stack

### Frontend Stack
```
React 18+ with TypeScript
├── Vite (build tool - fast HMR)
├── @dnd-kit/core + @dnd-kit/sortable (drag & drop)
├── @tanstack/react-query v5 (server state management)
├── Axios (HTTP client)
├── Tailwind CSS (styling)
├── React Router v6 (routing)
└── WebSocket client (real-time updates)
```

### Backend Stack
```
Go 1.21+
├── Gin (HTTP framework)
├── GORM (ORM for MySQL)
├── Gorilla WebSocket (WebSocket server)
├── go-git/go-git (Git operations)
├── github.com/google/go-github (GitHub API)
├── joho/godotenv (environment variables)
└── Anthropic SDK for Go (Claude API)
```

### Database
```
MySQL 8.0 (Local Installation)
└── Reasons:
    - ACID transactions for task queue integrity
    - Row-level locking prevents duplicate execution
    - JSON column support for flexible metadata
    - Excellent indexing for status queries
    - Mature, reliable for queue-based systems
```

### External APIs
- **Anthropic Claude API**: AI code/doc generation
- **GitHub API**: Pull request creation
- **Confluence REST API**: Documentation publishing
- **Git**: Repository operations

---

## 📂 Complete Project Structure

```
ai-task-manager/
│
├── frontend/                               # React Application
│   ├── public/
│   │   └── vite.svg
│   ├── src/
│   │   ├── components/
│   │   │   ├── Board/
│   │   │   │   ├── KanbanBoard.tsx         # Main board container
│   │   │   │   ├── Column.tsx              # Droppable column (Todo, In Progress, etc.)
│   │   │   │   ├── TaskCard.tsx            # Draggable task card
│   │   │   │   └── StatusBadge.tsx         # Colored status indicator
│   │   │   ├── Task/
│   │   │   │   ├── CreateTaskModal.tsx     # Task creation form
│   │   │   │   ├── TaskDetailModal.tsx     # Task details + logs viewer
│   │   │   │   ├── TaskForm.tsx            # Reusable task form
│   │   │   │   ├── TagSelector.tsx         # coding/documentation tags
│   │   │   │   ├── PrioritySelect.tsx      # Priority dropdown
│   │   │   │   └── MetadataForm.tsx        # Dynamic metadata inputs
│   │   │   ├── Execution/
│   │   │   │   ├── ExecutionLogs.tsx       # Log stream viewer
│   │   │   │   ├── LogEntry.tsx            # Single log line
│   │   │   │   └── AgentStatus.tsx         # Agent activity indicator
│   │   │   ├── Layout/
│   │   │   │   ├── MainLayout.tsx          # App shell
│   │   │   │   ├── Header.tsx              # Top navbar
│   │   │   │   └── Sidebar.tsx             # Side navigation (future)
│   │   │   └── Common/
│   │   │       ├── Button.tsx
│   │   │       ├── Modal.tsx
│   │   │       ├── Input.tsx
│   │   │       └── Toast.tsx               # Notification toasts
│   │   ├── hooks/
│   │   │   ├── useTasks.ts                 # React Query for tasks
│   │   │   ├── useWebSocket.ts             # WebSocket connection
│   │   │   ├── useTaskMutation.ts          # Create/update/delete
│   │   │   └── useExecutionLogs.ts         # Fetch logs
│   │   ├── services/
│   │   │   ├── api.ts                      # Axios instance + interceptors
│   │   │   └── websocket.ts                # WebSocket service
│   │   ├── types/
│   │   │   └── index.ts                    # TypeScript interfaces
│   │   ├── utils/
│   │   │   ├── statusColors.ts             # Status → color mapping
│   │   │   ├── dateUtils.ts                # Date formatting
│   │   │   └── constants.ts                # App constants
│   │   ├── App.tsx                         # Root component
│   │   ├── main.tsx                        # Entry point
│   │   └── index.css                       # Global styles
│   ├── .env.example
│   ├── .env.local                          # Local config (gitignored)
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   ├── tsconfig.json
│   └── index.html
│
├── backend/                                # Go Application
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                     # Entry point - starts server
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   │   ├── task_handler.go         # Task CRUD endpoints
│   │   │   │   ├── log_handler.go          # Execution log endpoints
│   │   │   │   ├── agent_handler.go        # Agent status endpoints
│   │   │   │   └── ws_handler.go           # WebSocket handler
│   │   │   ├── middleware/
│   │   │   │   ├── cors.go                 # CORS configuration
│   │   │   │   ├── logger.go               # Request logging
│   │   │   │   └── error.go                # Error handling
│   │   │   └── router.go                   # Route definitions
│   │   ├── models/
│   │   │   ├── task.go                     # Task model + GORM tags
│   │   │   ├── execution_log.go            # ExecutionLog model
│   │   │   └── agent_config.go             # AgentConfig model (optional)
│   │   ├── database/
│   │   │   ├── db.go                       # MySQL connection + config
│   │   │   └── migrations.go               # Auto-migrate schemas
│   │   ├── queue/
│   │   │   ├── queue.go                    # Task queue implementation
│   │   │   ├── worker.go                   # Worker pool
│   │   │   └── dispatcher.go               # Task dispatcher
│   │   ├── agents/
│   │   │   ├── agent.go                    # Agent interface
│   │   │   ├── registry.go                 # Agent registry
│   │   │   ├── router.go                   # Tag-based routing
│   │   │   ├── base_agent.go               # Base agent with logging
│   │   │   ├── coding_agent.go             # Coding task implementation
│   │   │   └── docs_agent.go               # Documentation implementation
│   │   ├── services/
│   │   │   ├── task_service.go             # Task business logic
│   │   │   ├── git_service.go              # Git operations wrapper
│   │   │   ├── github_service.go           # GitHub API client
│   │   │   ├── confluence_service.go       # Confluence API client
│   │   │   └── llm_service.go              # Claude API wrapper
│   │   ├── websocket/
│   │   │   ├── hub.go                      # WebSocket hub (broadcast)
│   │   │   └── client.go                   # WebSocket client handler
│   │   └── config/
│   │       └── config.go                   # Config struct + loader
│   ├── .env.example
│   ├── .env                                # Local config (gitignored)
│   ├── go.mod
│   ├── go.sum
│   └── README.md
│
├── .gitignore
├── docker-compose.yml                      # For future Docker setup
└── README.md                               # Project documentation
```

---

## 🗄️ Database Schema (MySQL)

### tasks table
```sql
CREATE TABLE tasks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    
    -- Basic info
    title VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Classification
    tags JSON NOT NULL COMMENT '["coding"] or ["documentation"]',
    priority INT DEFAULT 0 COMMENT 'Higher = more important',
    
    -- Status tracking
    status VARCHAR(50) NOT NULL DEFAULT 'todo' COMMENT 'todo, queued, in_progress, completed, failed',
    assigned_agent VARCHAR(100) COMMENT 'coding_agent, docs_agent',
    
    -- Results
    result TEXT COMMENT 'Success output (e.g., PR URL)',
    error_msg TEXT COMMENT 'Error message if failed',
    
    -- Task-specific data
    metadata JSON COMMENT 'Repo URL, Confluence space, etc.',
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    started_at TIMESTAMP NULL COMMENT 'When agent picked up',
    completed_at TIMESTAMP NULL COMMENT 'When finished',
    
    -- Indexes for performance
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_tags ((CAST(tags AS CHAR(255) ARRAY)))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### execution_logs table
```sql
CREATE TABLE execution_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    task_id BIGINT NOT NULL,
    
    -- Log details
    agent_name VARCHAR(100),
    log_level VARCHAR(20) NOT NULL COMMENT 'info, warning, error, success',
    message TEXT NOT NULL,
    
    -- Timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    
    -- Indexes
    INDEX idx_task_id (task_id),
    INDEX idx_created_at (created_at),
    INDEX idx_log_level (log_level)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### agent_configs table (Optional - for future use)
```sql
CREATE TABLE agent_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    agent_name VARCHAR(100) UNIQUE NOT NULL,
    
    -- Configuration
    config JSON COMMENT 'Agent-specific settings',
    enabled BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

---

## 🔄 Task Status Flow

```
┌─────────────┐
│ User Creates│
│    Task     │
└──────┬──────┘
       │
       ↓
   [todo] ────────────────→ Initial status on creation
       │                    Status: "todo"
       │                    Color: Gray
       │
       ↓ (Auto-queued immediately)
       │
   [queued] ──────────────→ Waiting for available worker
       │                    Status: "queued"
       │                    Color: Yellow
       │
       ↓ (Worker picks up task)
       │
[in_progress] ─────────────→ Agent actively executing
       │                    Status: "in_progress"
       │                    Color: Blue
       │                    Logs streaming in real-time
       │
       ↓
    ┌──┴──┐
    ↓     ↓
[completed] [failed]
Success!    Error occurred
Green       Red
```

### Status Definitions

1. **todo** (Gray)
   - Just created, not yet queued
   - Initial state before queue pickup
   - Duration: < 1 second

2. **queued** (Yellow)
   - In task queue channel
   - Waiting for available worker
   - Duration: Depends on queue length

3. **in_progress** (Blue)
   - Worker has picked up task
   - Agent is actively executing
   - Logs streaming to UI
   - Duration: Varies (30s - 5min typically)

4. **completed** (Green)
   - Successfully finished
   - Result contains output (e.g., PR URL)
   - Final state

5. **failed** (Red)
   - Error occurred during execution
   - ErrorMsg contains details
   - Can be manually retried
   - Final state (unless retried)

---

## 📡 API Design

### REST Endpoints

#### Task Management
```
GET    /api/tasks
       Query params: ?status=completed&tag=coding&limit=50
       Returns: List of tasks with pagination
       
POST   /api/tasks
       Body: { title, description, tags[], priority, metadata }
       Returns: Created task (auto-queued)
       
GET    /api/tasks/:id
       Returns: Single task with full details
       
PUT    /api/tasks/:id
       Body: { title, description, metadata }
       Returns: Updated task
       
DELETE /api/tasks/:id
       Returns: 204 No Content
       
PATCH  /api/tasks/:id/status
       Body: { status: "queued" }  # For manual retry
       Returns: Updated task
```

#### Execution Logs
```
GET    /api/tasks/:id/logs
       Returns: All logs for a task (chronological)
       
GET    /api/logs/recent
       Query params: ?limit=100
       Returns: Recent logs across all tasks
```

#### Agent Status
```
GET    /api/agents
       Returns: List of all agents + enabled status
       
GET    /api/agents/stats
       Returns: { queued: 3, in_progress: 2, completed: 45, failed: 1 }
```

#### Health Check
```
GET    /api/health
       Returns: { status: "ok", queue_size: 3, workers: 3 }
```

### WebSocket Protocol

**Connection**: `ws://localhost:8080/ws`

**Event Types** (Server → Client):

```typescript
// Task status changed
{
  type: 'task_status_changed',
  task_id: 123,
  status: 'in_progress',
  timestamp: '2026-01-14T10:30:00Z'
}

// New execution log
{
  type: 'execution_log',
  task_id: 123,
  log_level: 'info',
  message: 'Cloning repository...',
  timestamp: '2026-01-14T10:30:05Z'
}

// Task completed successfully
{
  type: 'task_completed',
  task_id: 123,
  result: 'Pull request created: https://github.com/user/repo/pull/456',
  timestamp: '2026-01-14T10:32:00Z'
}

// Task failed
{
  type: 'task_failed',
  task_id: 123,
  error: 'Git clone failed: repository not found',
  timestamp: '2026-01-14T10:31:00Z'
}

// Agent started working
{
  type: 'agent_started',
  task_id: 123,
  agent_name: 'coding_agent',
  timestamp: '2026-01-14T10:30:00Z'
}
```

---

## 🤖 Agent System Design

### Agent Interface (Go)

```go
package agents

import "ai-task-manager/internal/models"

// Agent is the interface that all agents must implement
type Agent interface {
    // Name returns the agent's unique identifier
    Name() string
    
    // CanHandle returns true if this agent can process tasks with the given tags
    CanHandle(tags []string) bool
    
    // Execute processes the task and returns an error if it fails
    // The agent should:
    // 1. Update task status
    // 2. Create execution logs
    // 3. Perform the actual work
    // 4. Update task result or error_msg
    Execute(task *models.Task) error
}
```

### Task Model (Go)

```go
package models

import (
    "database/sql/driver"
    "encoding/json"
    "time"
)

type Task struct {
    ID            uint           `json:"id" gorm:"primaryKey"`
    Title         string         `json:"title" gorm:"not null"`
    Description   string         `json:"description" gorm:"type:text"`
    Tags          StringArray    `json:"tags" gorm:"type:json;not null"`
    Status        string         `json:"status" gorm:"not null;default:todo"`
    Priority      int            `json:"priority" gorm:"default:0"`
    AssignedAgent string         `json:"assigned_agent"`
    Result        string         `json:"result" gorm:"type:text"`
    ErrorMsg      string         `json:"error_msg" gorm:"type:text"`
    Metadata      JSON           `json:"metadata" gorm:"type:json"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    StartedAt     *time.Time     `json:"started_at"`
    CompletedAt   *time.Time     `json:"completed_at"`
}

// StringArray for storing []string in JSON
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
    return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
    return json.Unmarshal(value.([]byte), a)
}

// JSON for storing map[string]interface{} in JSON
type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
    return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
    return json.Unmarshal(value.([]byte), j)
}
```

---

## 🔧 Coding Agent Specification

### Capabilities
1. ✅ Clone Git repository from metadata
2. ✅ Create feature branch: `ai-task-{task_id}`
3. ✅ Use Claude API to generate code changes
4. ✅ Apply changes to specified files
5. ✅ Commit with descriptive message
6. ✅ Push to remote
7. ✅ Create GitHub Pull Request
8. ✅ Add PR URL to task result

### Task Metadata Schema (JSON)
```json
{
  "repo_url": "https://github.com/username/repo",
  "branch": "main",
  "files_to_modify": ["src/components/Button.tsx"],
  "instruction": "Add loading state to Button component with spinner icon"
}
```

### Execution Flow
```
1. [INFO]    Starting coding task for: {title}
2. [INFO]    Cloning repository from {repo_url}...
3. [INFO]    Repository cloned successfully
4. [INFO]    Creating branch: ai-task-{id}
5. [INFO]    Branch created successfully
6. [INFO]    Reading existing code from {files}...
7. [INFO]    Analyzing code structure...
8. [INFO]    Generating code changes with Claude AI...
9. [INFO]    AI generation complete
10. [INFO]   Applying changes to {files}...
11. [INFO]   Changes applied successfully
12. [INFO]   Committing changes...
13. [INFO]   Commit created: {commit_hash}
14. [INFO]   Pushing to remote repository...
15. [INFO]   Push successful
16. [INFO]   Creating pull request on GitHub...
17. [SUCCESS] ✓ Pull request created: {pr_url}

Status: completed
Result: "Pull request created: https://github.com/user/repo/pull/123"
```

### Error Handling
```
Git clone fails
├→ Log: [ERROR] Failed to clone repository: {error}
├→ Set status: failed
└→ Set error_msg: "Git clone failed: {details}"

AI generation fails
├→ Log: [WARNING] First attempt failed, retrying...
├→ Retry once
└→ If still fails:
   ├→ Log: [ERROR] AI generation failed after retry
   ├→ Set status: failed
   └→ Set error_msg: "AI generation failed: {details}"

Push fails (e.g., auth error)
├→ Log: [ERROR] Failed to push: {error}
├→ Set status: failed
└→ Set error_msg: "Git push failed: {details}"

PR creation fails
├→ Log: [WARNING] PR creation failed but code is pushed
├→ Set status: completed
├→ Set result: "Code pushed successfully. Manual PR needed."
└→ Note: Still mark as completed since code changes are saved
```

### Implementation Notes
- Use `go-git` for Git operations
- Use Anthropic SDK for Claude API calls
- Use `google/go-github` for PR creation
- Clone to temp directory: `/tmp/ai-task-manager/task-{id}/`
- Clean up temp directory after completion
- Include task description in PR body
- Set PR title to task title

---

## 📝 Documentation Agent Specification

### Capabilities
1. ✅ Connect to Confluence using API token
2. ✅ Search for existing page by title
3. ✅ Use Claude API to generate documentation
4. ✅ Format content in Confluence Storage Format
5. ✅ Create new page or update existing
6. ✅ Add labels/tags to page
7. ✅ Return page URL in result

### Task Metadata Schema (JSON)
```json
{
  "confluence_space": "DEV",
  "page_title": "API Documentation",
  "parent_page_id": 123456,
  "content_source": "file",
  "source_value": "src/api/routes.go",
  "instruction": "Generate comprehensive API documentation from route definitions"
}
```

**Content Source Options**:
- `"file"` → Read from file path in source_value
- `"url"` → Fetch from URL in source_value
- `"text"` → Use source_value directly as content

### Execution Flow
```
1. [INFO]    Starting documentation task for: {title}
2. [INFO]    Connecting to Confluence at {url}...
3. [INFO]    Connection successful
4. [INFO]    Searching for existing page: {page_title}
5. [INFO]    Page found / Page not found
6. [INFO]    Reading source content from {source}...
7. [INFO]    Content loaded successfully
8. [INFO]    Generating documentation with Claude AI...
9. [INFO]    AI generation complete
10. [INFO]   Formatting content for Confluence...
11. [INFO]   Creating/Updating page in space {space}...
12. [INFO]   Page saved successfully
13. [INFO]   Adding labels: documentation, ai-generated
14. [SUCCESS] ✓ Documentation published: {page_url}

Status: completed
Result: "Documentation published: https://domain.atlassian.net/wiki/spaces/DEV/pages/123456789"
```

### Confluence Storage Format
Convert markdown to Confluence format:
```
Markdown                  →  Confluence Storage Format
-----------------------------------------------------------
# Heading 1              →  <h1>Heading 1</h1>
## Heading 2             →  <h2>Heading 2</h2>
**bold**                 →  <strong>bold</strong>
*italic*                 →  <em>italic</em>
`code`                   →  <code>code</code>
```code block```        →  <ac:structured-macro ac:name="code">...
[link](url)              →  <a href="url">link</a>
- bullet                 →  <ul><li>bullet</li></ul>
```

### Implementation Notes
- Use Confluence REST API v2
- Authentication: Basic auth (username + API token)
- Check if page exists before creating
- Update existing page version number
- Add metadata: `ai-generated`, `auto-updated`
- Handle rate limits (100 requests/minute)

---

## ⚙️ Task Queue Architecture

### Queue Implementation (Go Channels)

```go
package queue

import (
    "ai-task-manager/internal/models"
    "ai-task-manager/internal/websocket"
    "ai-task-manager/internal/agents"
    "gorm.io/gorm"
    "time"
)

type TaskQueue struct {
    taskChannel chan *models.Task
    maxWorkers  int
    db          *gorm.DB
    wsHub       *websocket.Hub
    agentRouter *agents.Router
    stopChan    chan struct{}
}

func NewTaskQueue(
    maxWorkers int,
    db *gorm.DB,
    wsHub *websocket.Hub,
    router *agents.Router,
) *TaskQueue {
    return &TaskQueue{
        taskChannel: make(chan *models.Task, 100),
        maxWorkers:  maxWorkers,
        db:          db,
        wsHub:       wsHub,
        agentRouter: router,
        stopChan:    make(chan struct{}),
    }
}

func (q *TaskQueue) Start() {
    // Start polling for queued tasks
    go q.pollDatabaseForTasks()
    
    // Start worker pool
    for i := 0; i < q.maxWorkers; i++ {
        go q.worker(i)
    }
}

func (q *TaskQueue) pollDatabaseForTasks() {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Fetch queued tasks with row-level locking
            var tasks []models.Task
            q.db.Where("status = ?", "queued").
                Order("priority DESC, created_at ASC").
                Limit(10).
                Find(&tasks)
            
            // Add to channel
            for _, task := range tasks {
                select {
                case q.taskChannel <- &task:
                    // Task added to queue
                default:
                    // Channel full, will retry next poll
                }
            }
        case <-q.stopChan:
            return
        }
    }
}

func (q *TaskQueue) worker(workerID int) {
    for task := range q.taskChannel {
        q.executeTask(task)
    }
}

func (q *TaskQueue) executeTask(task *models.Task) {
    // Update status to in_progress
    now := time.Now()
    task.Status = "in_progress"
    task.StartedAt = &now
    q.db.Save(task)
    
    // Broadcast status change
    q.wsHub.BroadcastTaskStatus(task.ID, "in_progress")
    
    // Get appropriate agent
    agent := q.agentRouter.GetAgentForTask(task)
    
    // Broadcast agent started
    q.wsHub.BroadcastAgentStarted(task.ID, agent.Name())
    
    // Execute task
    err := agent.Execute(task)
    
    // Update final status
    completedTime := time.Now()
    task.CompletedAt = &completedTime
    
    if err != nil {
        task.Status = "failed"
        task.ErrorMsg = err.Error()
        q.wsHub.BroadcastTaskFailed(task.ID, err.Error())
    } else {
        task.Status = "completed"
        q.wsHub.BroadcastTaskCompleted(task.ID, task.Result)
    }
    
    q.db.Save(task)
}

func (q *TaskQueue) Stop() {
    close(q.stopChan)
    close(q.taskChannel)
}

func (q *TaskQueue) EnqueueTask(taskID uint) error {
    var task models.Task
    if err := q.db.First(&task, taskID).Error; err != nil {
        return err
    }
    
    task.Status = "queued"
    q.db.Save(&task)
    
    return nil
}
```

### Key Features
- **3 concurrent workers** (configurable via WORKER_COUNT)
- **Polling interval**: 2 seconds
- **Priority-based**: Higher priority tasks executed first
- **Row-level locking**: Prevent duplicate execution
- **Graceful shutdown**: Close channels, wait for workers
- **Channel buffer**: 100 tasks
- **Auto-queue on create**: POST /api/tasks automatically queues

---

## 🌐 Environment Configuration

### .env File (Backend)
```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
FRONTEND_URL=http://localhost:5173

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ai_task_manager
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_CHARSET=utf8mb4

# Task Queue Configuration
WORKER_COUNT=3
QUEUE_POLL_INTERVAL=2s
QUEUE_BUFFER_SIZE=100

# AI Provider (Anthropic Claude)
ANTHROPIC_API_KEY=sk-ant-api03-xxxxxxxxxxxxxxxxxxxxx
ANTHROPIC_MODEL=claude-3-5-sonnet-20241022
ANTHROPIC_MAX_TOKENS=4096

# Git & GitHub Configuration
GIT_USERNAME=your-github-username
GIT_EMAIL=your-email@example.com
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx

# Confluence Configuration (Optional - user provides when needed)
CONFLUENCE_URL=https://your-domain.atlassian.net/wiki
CONFLUENCE_USERNAME=your-email@example.com
CONFLUENCE_API_TOKEN=your-confluence-api-token

# Agent Configuration
CODING_AGENT_ENABLED=true
DOCS_AGENT_ENABLED=true
AUTO_EXECUTE_TASKS=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Temp Directory for Git Operations
TEMP_DIR=/tmp/ai-task-manager
```

### .env.local File (Frontend)
```bash
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080/ws
```

---

## 📝 Phase-by-Phase Implementation Plan

### Phase 1: Project Setup & Database (Est: 2-3 hours)

**Backend Setup:**
- [ ] Initialize Go module: `go mod init ai-task-manager`
- [ ] Create directory structure (cmd/, internal/)
- [ ] Install dependencies: gin, gorm, mysql driver, godotenv
- [ ] Create .env.example with all config variables
- [ ] Set up MySQL database connection
- [ ] Implement auto-migration for tables
- [ ] Test database connection

**Frontend Setup:**
- [ ] Create Vite React app: `npm create vite@latest frontend -- --template react-ts`
- [ ] Install dependencies: @dnd-kit, @tanstack/react-query, axios, tailwindcss
- [ ] Configure Tailwind CSS
- [ ] Create .env.local with API URLs
- [ ] Set up basic folder structure
- [ ] Create index.css with Tailwind directives

**Database:**
- [ ] Create MySQL database: `CREATE DATABASE ai_task_manager;`
- [ ] Verify connection from Go app
- [ ] Run migrations (auto-migrate with GORM)
- [ ] Insert test data (optional)

**Verification:**
- [ ] Backend starts without errors
- [ ] Frontend dev server runs
- [ ] Database tables created successfully

---

### Phase 2: Backend Foundation (Est: 4-5 hours)

- [ ] **Config Management**
  - [ ] Implement config/config.go with struct for all env vars
  - [ ] Load config using godotenv
  - [ ] Validate required config on startup

- [ ] **Models**
  - [ ] Implement models/task.go with GORM tags
  - [ ] Implement models/execution_log.go
  - [ ] Add JSON field helpers (StringArray, JSON types)
  - [ ] Add model validation

- [ ] **API Router Setup**
  - [ ] Create api/router.go with Gin router
  - [ ] Set up CORS middleware for http://localhost:5173
  - [ ] Add logging middleware
  - [ ] Add error handling middleware
  - [ ] Create health check endpoint

- [ ] **Task Handlers**
  - [ ] POST /api/tasks - Create task (auto-queue)
  - [ ] GET /api/tasks - List tasks with filters
  - [ ] GET /api/tasks/:id - Get single task
  - [ ] PUT /api/tasks/:id - Update task
  - [ ] DELETE /api/tasks/:id - Delete task
  - [ ] PATCH /api/tasks/:id/status - Update status

- [ ] **Log Handlers**
  - [ ] GET /api/tasks/:id/logs - Get logs for task
  - [ ] GET /api/logs/recent - Recent logs

- [ ] **WebSocket Hub**
  - [ ] Implement websocket/hub.go (connection manager)
  - [ ] Implement websocket/client.go (client handler)
  - [ ] Add broadcast methods
  - [ ] Create WS endpoint: /ws

- [ ] **Testing**
  - [ ] Test all endpoints with curl/Postman
  - [ ] Verify CORS works from frontend
  - [ ] Test WebSocket connection

---

### Phase 3: Task Queue System (Est: 3-4 hours)

- [ ] **Queue Implementation**
  - [ ] Implement queue/queue.go with channel-based queue
  - [ ] Add polling mechanism (ticker every 2s)
  - [ ] Implement priority-based fetching
  - [ ] Add worker pool (3 workers)

- [ ] **Worker Implementation**
  - [ ] Implement worker function
  - [ ] Add task status transitions
  - [ ] Integrate with WebSocket hub for broadcasts
  - [ ] Add error handling

- [ ] **Integration**
  - [ ] Start queue on server startup
  - [ ] Auto-queue tasks on POST /api/tasks
  - [ ] Test with mock task execution

- [ ] **Graceful Shutdown**
  - [ ] Handle SIGINT/SIGTERM signals
  - [ ] Wait for active tasks to complete
  - [ ] Close channels properly

- [ ] **Testing**
  - [ ] Create 5 tasks, verify sequential execution
  - [ ] Verify status transitions via WebSocket
  - [ ] Test concurrent execution (3 workers)

---

### Phase 4: Agent System Foundation (Est: 3-4 hours)

- [ ] **Agent Interface**
  - [ ] Define Agent interface in agents/agent.go
  - [ ] Create base agent struct with common methods

- [ ] **Agent Registry**
  - [ ] Implement agents/registry.go
  - [ ] Register coding agent
  - [ ] Register docs agent

- [ ] **Agent Router**
  - [ ] Implement agents/router.go
  - [ ] Route by tags: ["coding"] → CodingAgent
  - [ ] Route by tags: ["documentation"] → DocsAgent

- [ ] **Base Agent Utilities**
  - [ ] Helper: CreateExecutionLog(task, level, message)
  - [ ] Helper: UpdateTaskStatus(task, status)
  - [ ] Helper: BroadcastLog(task, message)

- [ ] **Mock Agent**
  - [ ] Create mock agent for testing
  - [ ] Simulate 5-second task execution
  - [ ] Add sample logs

- [ ] **Testing**
  - [ ] Create task with ["coding"] tag → routes to coding agent
  - [ ] Create task with ["documentation"] tag → routes to docs agent
  - [ ] Verify logs created in database

---

### Phase 5: Coding Agent (Est: 6-8 hours)

- [ ] **Git Service**
  - [ ] Implement services/git_service.go
  - [ ] Clone repository to temp directory
  - [ ] Create branch
  - [ ] Commit changes
  - [ ] Push to remote (with auth)

- [ ] **GitHub Service**
  - [ ] Implement services/github_service.go
  - [ ] Use google/go-github library
  - [ ] Create pull request
  - [ ] Add PR body with task details

- [ ] **LLM Service**
  - [ ] Implement services/llm_service.go
  - [ ] Integrate Anthropic Claude SDK
  - [ ] Create prompt template for code generation
  - [ ] Parse AI response

- [ ] **Coding Agent**
  - [ ] Implement agents/coding_agent.go
  - [ ] Implement Execute() method
  - [ ] Full execution flow (10 steps from spec)
  - [ ] Add comprehensive error handling
  - [ ] Add retry logic for AI calls

- [ ] **Testing**
  - [ ] Create test repository on GitHub
  - [ ] Create coding task via API
  - [ ] Verify PR created successfully
  - [ ] Test error scenarios (invalid repo, auth failure)

---

### Phase 6: Documentation Agent (Est: 5-6 hours)

- [ ] **Confluence Service**
  - [ ] Implement services/confluence_service.go
  - [ ] Authenticate with API token
  - [ ] Search for page by title
  - [ ] Create new page
  - [ ] Update existing page
  - [ ] Add labels to page

- [ ] **Markdown to Confluence Converter**
  - [ ] Convert markdown to Confluence Storage Format
  - [ ] Handle headings, bold, italic, code blocks
  - [ ] Handle lists and links

- [ ] **Documentation Agent**
  - [ ] Implement agents/docs_agent.go
  - [ ] Implement Execute() method
  - [ ] Full execution flow (8 steps from spec)
  - [ ] Add error handling

- [ ] **Testing**
  - [ ] Create Confluence space for testing
  - [ ] Create documentation task via API
  - [ ] Verify page created/updated
  - [ ] Test with different content sources (file, text)

---

### Phase 7: Frontend - Core UI (Est: 5-6 hours)

- [ ] **Layout Components**
  - [ ] MainLayout.tsx with header
  - [ ] Header.tsx with title and "New Task" button

- [ ] **Kanban Board**
  - [ ] KanbanBoard.tsx with DndContext
  - [ ] Column.tsx for each status (todo, queued, in_progress, completed, failed)
  - [ ] Make columns droppable

- [ ] **Task Card**
  - [ ] TaskCard.tsx with draggable functionality
  - [ ] Display title, tags, priority
  - [ ] StatusBadge.tsx with color coding
  - [ ] Click to open detail modal

- [ ] **Styling**
  - [ ] Tailwind classes for Kanban board
  - [ ] Responsive layout
  - [ ] Card shadows and hover effects
  - [ ] Status badge colors

- [ ] **Testing**
  - [ ] Render board with mock data
  - [ ] Test drag and drop between columns
  - [ ] Verify responsive design

---

### Phase 8: Frontend - Task Management (Est: 4-5 hours)

- [ ] **Create Task Modal**
  - [ ] CreateTaskModal.tsx with form
  - [ ] Title and description inputs
  - [ ] Tag selector (coding/documentation)
  - [ ] Priority selector (low/medium/high)
  - [ ] Dynamic metadata form based on tag

- [ ] **Metadata Forms**
  - [ ] Coding task metadata (repo_url, branch, files, instruction)
  - [ ] Documentation task metadata (confluence_space, page_title, etc.)
  - [ ] Form validation

- [ ] **Task Detail Modal**
  - [ ] TaskDetailModal.tsx with tabs
  - [ ] Tab 1: Task info (editable)
  - [ ] Tab 2: Execution logs (real-time)
  - [ ] Edit and Delete buttons

- [ ] **Common Components**
  - [ ] Button.tsx (primary, secondary, danger)
  - [ ] Input.tsx with label
  - [ ] Modal.tsx wrapper
  - [ ] Toast.tsx for notifications

- [ ] **Testing**
  - [ ] Open create modal, fill form, submit
  - [ ] Open task detail, verify info displayed
  - [ ] Edit task title and save
  - [ ] Delete task with confirmation

---

### Phase 9: Frontend - API Integration (Est: 3-4 hours)

- [ ] **API Service**
  - [ ] services/api.ts with Axios instance
  - [ ] Base URL from env: VITE_API_BASE_URL
  - [ ] Request/response interceptors
  - [ ] Error handling

- [ ] **TypeScript Types**
  - [ ] types/index.ts with Task, ExecutionLog, Agent interfaces
  - [ ] Metadata types for coding/docs tasks

- [ ] **React Query Hooks**
  - [ ] hooks/useTasks.ts for fetching tasks
  - [ ] hooks/useTaskMutation.ts for create/update/delete
  - [ ] hooks/useExecutionLogs.ts for logs
  - [ ] Set up QueryClient in main.tsx

- [ ] **Integration**
  - [ ] Fetch tasks on board mount
  - [ ] Create task on form submit
  - [ ] Update task on edit
  - [ ] Delete task with mutation
  - [ ] Optimistic updates for better UX

- [ ] **Toast Notifications**
  - [ ] Success: "Task created successfully"
  - [ ] Error: "Failed to create task: {error}"
  - [ ] Show on mutation success/error

- [ ] **Testing**
  - [ ] Create task → appears on board
  - [ ] Edit task → updates immediately
  - [ ] Delete task → removed from board
  - [ ] Test error handling (backend down)

---

### Phase 10: Real-time Updates (Est: 4-5 hours)

- [ ] **WebSocket Service**
  - [ ] services/websocket.ts with connection management
  - [ ] Auto-reconnect on disconnect
  - [ ] Event listener registration
  - [ ] Parse incoming messages

- [ ] **useWebSocket Hook**
  - [ ] Connect on mount
  - [ ] Subscribe to events
  - [ ] Return connection status
  - [ ] Clean up on unmount

- [ ] **Event Handlers**
  - [ ] task_status_changed → update task in React Query cache
  - [ ] execution_log → append to logs list
  - [ ] task_completed → show success toast + update task
  - [ ] task_failed → show error toast + update task
  - [ ] agent_started → show "Agent working..." indicator

- [ ] **UI Updates**
  - [ ] Move task card to new column on status change
  - [ ] Stream logs in task detail modal
  - [ ] Auto-scroll logs to bottom
  - [ ] Show agent status badge on card

- [ ] **Testing**
  - [ ] Create task, watch it move through columns
  - [ ] Open task detail, watch logs stream
  - [ ] Verify real-time updates without refresh
  - [ ] Test with multiple browser tabs

---

### Phase 11: Execution Logs Viewer (Est: 2-3 hours)

- [ ] **ExecutionLogs Component**
  - [ ] Display logs in chronological order
  - [ ] Color-code by log level (info=gray, success=green, warning=yellow, error=red)
  - [ ] Add timestamps (format: HH:mm:ss)
  - [ ] Auto-scroll toggle button

- [ ] **LogEntry Component**
  - [ ] Single log line with icon
  - [ ] Truncate long messages with "Show more"
  - [ ] Monospace font for technical logs

- [ ] **Features**
  - [ ] Auto-scroll by default
  - [ ] Pause auto-scroll on user scroll up
  - [ ] Filter by log level (show all / errors only)
  - [ ] Copy logs to clipboard button

- [ ] **Testing**
  - [ ] Create task, watch logs stream in real-time
  - [ ] Verify color coding
  - [ ] Test auto-scroll behavior
  - [ ] Test filter functionality

---

### Phase 12: Polish & Error Handling (Est: 3-4 hours)

- [ ] **Loading States**
  - [ ] Skeleton loaders for task cards
  - [ ] Spinner on task creation
  - [ ] Loading indicator in modals
  - [ ] Disabled buttons while loading

- [ ] **Empty States**
  - [ ] "No tasks yet" message on empty board
  - [ ] Empty column illustrations
  - [ ] "No logs" message

- [ ] **Error Boundaries**
  - [ ] Wrap app in error boundary
  - [ ] Show friendly error message
  - [ ] "Reload" button

- [ ] **Error Messages**
  - [ ] Improve API error messages
  - [ ] User-friendly wording
  - [ ] Include action items

- [ ] **Confirmation Dialogs**
  - [ ] "Are you sure?" on task delete
  - [ ] Warning if deleting in-progress task

- [ ] **Retry Failed Tasks**
  - [ ] "Retry" button on failed task cards
  - [ ] Call PATCH /api/tasks/:id/status with status=queued

- [ ] **Testing**
  - [ ] Test all loading states
  - [ ] Trigger errors (invalid input, backend down)
  - [ ] Delete task with confirmation
  - [ ] Retry failed task successfully

---

### Phase 13: Testing & Documentation (Est: 3-4 hours)

- [ ] **End-to-End Testing**
  - [ ] Test full coding task flow (create → execute → PR created)
  - [ ] Test full documentation task flow (create → execute → page published)
  - [ ] Test concurrent execution (3 tasks at once)
  - [ ] Test task priority (high priority first)
  - [ ] Test error scenarios (invalid repo, wrong credentials)

- [ ] **README.md**
  - [ ] Project description
  - [ ] Features list
  - [ ] Tech stack
  - [ ] Prerequisites
  - [ ] Installation instructions
  - [ ] Configuration (.env setup)
  - [ ] Usage examples
  - [ ] Screenshots (optional)
  - [ ] Troubleshooting section

- [ ] **.env.example Files**
  - [ ] Backend: Document all variables
  - [ ] Frontend: Document VITE_ variables

- [ ] **Code Comments**
  - [ ] Add comments to complex functions
  - [ ] Document agent execution flow
  - [ ] Add package-level comments

- [ ] **API Documentation**
  - [ ] Document all endpoints (optional: use Swagger)
  - [ ] Include request/response examples

---

### Phase 14: Deployment Preparation (Est: 2-3 hours)

- [ ] **Docker Compose (Optional)**
  - [ ] Create docker-compose.yml
  - [ ] Service: MySQL database
  - [ ] Service: Go backend
  - [ ] Service: React frontend (nginx)
  - [ ] Volume for MySQL data persistence

- [ ] **Environment Variables**
  - [ ] Production-ready .env.example
  - [ ] Security notes (don't commit .env)

- [ ] **Health Checks**
  - [ ] GET /api/health endpoint
  - [ ] Return queue stats
  - [ ] Return database status

- [ ] **Graceful Shutdown**
  - [ ] Handle SIGTERM in Go app
  - [ ] Wait for active tasks
  - [ ] Close DB connections

- [ ] **Build Scripts**
  - [ ] Backend: `go build -o bin/server cmd/server/main.go`
  - [ ] Frontend: `npm run build`
  - [ ] Combined: `make build`

- [ ] **Documentation**
  - [ ] Deployment guide (local)
  - [ ] Docker instructions
  - [ ] Production considerations

---

## 🎯 Success Criteria

The project is **complete and successful** when:

### Core Functionality
1. ✅ User can create a coding task with repo URL and instruction
2. ✅ Task auto-executes immediately after creation
3. ✅ Task status updates in real-time: todo → queued → in_progress → completed
4. ✅ Real-time logs stream to UI as agent works
5. ✅ Coding agent successfully creates GitHub PR
6. ✅ PR link appears in task result field
7. ✅ User can create documentation task with Confluence details
8. ✅ Documentation agent generates and publishes to Confluence
9. ✅ Confluence page URL appears in task result

### UI/UX
10. ✅ Kanban board displays tasks in correct columns
11. ✅ Drag and drop works smoothly between columns
12. ✅ Status badges show correct colors (gray/yellow/blue/green/red)
13. ✅ Task detail modal shows all task info
14. ✅ Execution logs stream in real-time with color coding
15. ✅ Toast notifications on success/failure

### Concurrency
16. ✅ Multiple tasks execute concurrently (3 workers)
17. ✅ Queue handles priority correctly (high priority first)
18. ✅ No duplicate task execution (MySQL locking works)

### Error Handling
19. ✅ Failed tasks appear in Failed column with error details
20. ✅ User can retry failed tasks
21. ✅ Graceful error messages displayed

### Integration
22. ✅ GitHub PR creation works with real repository
23. ✅ Confluence page creation works with real space
24. ✅ Claude API generates quality code/documentation

---

## ⏱️ Estimated Timeline

- **Total Estimated Time**: 50-60 hours
- **Part-time (4 hours/day)**: 12-15 days
- **Full-time (8 hours/day)**: 6-8 days

### Critical Path
```
Phase 1-2: Backend Foundation (6-8 hrs)
    ↓
Phase 3-4: Queue + Agent System (6-8 hrs)
    ↓
Phase 5: Coding Agent (6-8 hrs)
    ↓
Phase 6: Docs Agent (5-6 hrs)
    ↓
Phase 7-9: Frontend (12-15 hrs)
    ↓
Phase 10-12: Real-time + Polish (9-12 hrs)
    ↓
Phase 13-14: Testing + Docs (5-7 hrs)
```

---

## 🔒 Security Considerations

### Credentials
- ✅ Store all credentials in .env (never commit)
- ✅ Use .gitignore for .env files
- ✅ Provide .env.example with dummy values

### API Tokens
- ✅ GitHub token: Minimum scope (`repo` only)
- ✅ Confluence token: Restrict to specific space
- ✅ Anthropic key: Set spending limits

### Code Execution
- ✅ Always create PRs (never push to main directly)
- ✅ Clone to temp directory (isolated)
- ✅ Clean up temp files after execution
- ✅ Human review required before merging PRs

### Input Validation
- ✅ Validate all API inputs (title, description, metadata)
- ✅ Sanitize user input before passing to AI
- ✅ Validate URLs before cloning

---

## 📚 Key Dependencies

### Backend (Go)
```go
require (
    github.com/gin-gonic/gin v1.9.1              // HTTP framework
    gorm.io/gorm v1.25.5                         // ORM
    gorm.io/driver/mysql v1.5.2                  // MySQL driver
    github.com/gorilla/websocket v1.5.1          // WebSocket
    github.com/go-git/go-git/v5 v5.11.0          // Git operations
    github.com/google/go-github/v57 v57.0.0      // GitHub API
    github.com/joho/godotenv v1.5.1              // .env loader
    github.com/liushuangls/go-anthropic v0.5.0   // Anthropic Claude
)
```

### Frontend (React)
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.0",
    "@dnd-kit/core": "^6.1.0",
    "@dnd-kit/sortable": "^8.0.0",
    "@tanstack/react-query": "^5.14.0",
    "axios": "^1.6.2",
    "zustand": "^4.4.7"
  },
  "devDependencies": {
    "@types/react": "^18.2.43",
    "typescript": "^5.3.3",
    "vite": "^5.0.8",
    "tailwindcss": "^3.4.0",
    "autoprefixer": "^10.4.16",
    "postcss": "^8.4.32"
  }
}
```

---

## 🚀 Quick Start (Post-Implementation)

### Prerequisites
```bash
# Required
✅ Go 1.21+
✅ Node.js 18+
✅ MySQL 8.0

# API Keys (user provides)
✅ Anthropic API key
✅ GitHub personal access token (scope: repo)
✅ (Optional) Confluence API token
```

### Installation
```bash
# 1. Clone repository
git clone <repo-url>
cd ai-task-manager

# 2. Setup backend
cd backend
cp .env.example .env
# Edit .env with your API keys and database credentials
go mod download
go run cmd/server/main.go

# 3. Setup frontend (new terminal)
cd frontend
npm install
cp .env.example .env.local
npm run dev

# 4. Create MySQL database
mysql -u root -p
CREATE DATABASE ai_task_manager;
exit

# 5. Access application
# Frontend: http://localhost:5173
# Backend: http://localhost:8080
```

---

## 📖 Usage Examples

### Example 1: Coding Task
```json
POST /api/tasks
{
  "title": "Add dark mode toggle",
  "description": "Implement dark mode toggle in settings with persistence",
  "tags": ["coding"],
  "priority": 2,
  "metadata": {
    "repo_url": "https://github.com/user/my-app",
    "branch": "main",
    "files_to_modify": ["src/components/Settings.tsx"],
    "instruction": "Add dark mode toggle with localStorage persistence"
  }
}
```

**Result**: PR created at https://github.com/user/my-app/pull/123

### Example 2: Documentation Task
```json
POST /api/tasks
{
  "title": "Update API docs",
  "description": "Document new user endpoints",
  "tags": ["documentation"],
  "priority": 1,
  "metadata": {
    "confluence_space": "DEV",
    "page_title": "User API Documentation",
    "content_source": "file",
    "source_value": "backend/internal/api/handlers/user_handler.go",
    "instruction": "Generate comprehensive API documentation"
  }
}
```

**Result**: Page published at https://domain.atlassian.net/wiki/spaces/DEV/pages/123456

---

## 🐛 Troubleshooting

### Issue: Tasks stuck in "queued"
**Cause**: Workers not running or crashed  
**Solution**: Check backend logs, restart server

### Issue: GitHub PR creation fails
**Cause**: Invalid token or insufficient permissions  
**Solution**: Regenerate token with `repo` scope

### Issue: WebSocket not connecting
**Cause**: CORS or port issue  
**Solution**: Verify FRONTEND_URL in .env matches frontend origin

### Issue: MySQL connection refused
**Cause**: MySQL not running or wrong credentials  
**Solution**: Start MySQL, verify DB_HOST/DB_PORT/DB_PASSWORD

---

## 🔄 Future Enhancements

### Phase 2 (Post-MVP)
- [ ] User authentication (JWT)
- [ ] Multi-user support
- [ ] Task templates
- [ ] Task dependencies (task B after task A)
- [ ] Cron-like scheduling
- [ ] Email/Slack notifications
- [ ] Analytics dashboard
- [ ] Custom agent plugins
- [ ] Task comments
- [ ] File attachments

---

## 📝 Final Notes

- All credentials in .env file (not database)
- Tasks auto-execute on creation
- Coding agent ALWAYS creates PRs
- Status updates in real-time via WebSocket
- Failed tasks can be retried from UI
- MySQL prevents duplicate execution
- All timestamps in UTC

---

**END OF IMPLEMENTATION PLAN**

Ready to proceed with implementation!
