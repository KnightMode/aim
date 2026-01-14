# Backend - AI Task Manager

Go backend service for the AI Task Manager application.

## Tech Stack

- Go 1.21+
- Gin (HTTP framework)
- GORM (ORM for MySQL)
- Gorilla WebSocket
- go-git (Git operations)
- go-github (GitHub API)
- Anthropic SDK (Claude API)

## Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                     # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/                   # HTTP handlers
│   │   ├── middleware/                 # Middleware (CORS, logging)
│   │   └── router.go                   # Route definitions
│   ├── models/                         # Database models
│   ├── database/                       # DB connection & migrations
│   ├── queue/                          # Task queue implementation
│   ├── agents/                         # Agent implementations
│   ├── services/                       # External services (Git, GitHub, etc.)
│   ├── websocket/                      # WebSocket hub
│   └── config/                         # Configuration management
├── .env.example                        # Environment variables template
├── go.mod                              # Go dependencies
└── README.md                           # This file
```

## Setup

1. Copy `.env.example` to `.env` and configure:
```bash
cp .env.example .env
```

2. Install dependencies:
```bash
go mod download
```

3. Run database migrations:
The application will auto-migrate schemas on startup using GORM.

4. Start the server:
```bash
go run cmd/server/main.go
```

Or use npm script from root:
```bash
npm run backend:dev
```

## API Endpoints

### Tasks
- `GET /api/tasks` - List all tasks
- `POST /api/tasks` - Create new task (auto-queued)
- `GET /api/tasks/:id` - Get single task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task
- `PATCH /api/tasks/:id/status` - Update task status

### Logs
- `GET /api/tasks/:id/logs` - Get execution logs for a task
- `GET /api/logs/recent` - Get recent logs

### Agents
- `GET /api/agents` - List all agents
- `GET /api/agents/stats` - Get agent statistics

### WebSocket
- `WS /ws` - WebSocket connection for real-time updates

### Health
- `GET /api/health` - Health check endpoint

## Environment Variables

See `.env.example` for all available configuration options.

Required variables:
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`
- `ANTHROPIC_API_KEY`
- `GITHUB_TOKEN`
- `GIT_USERNAME`, `GIT_EMAIL`

Optional:
- `CONFLUENCE_URL`, `CONFLUENCE_USERNAME`, `CONFLUENCE_API_TOKEN`

## Development

Build the application:
```bash
go build -o bin/server cmd/server/main.go
```

Run tests:
```bash
go test ./...
```

Clean build artifacts:
```bash
rm -rf bin/
```

## Agent System

The backend implements two types of agents:

### Coding Agent
- Clones Git repositories
- Generates code changes using Claude AI
- Creates GitHub pull requests

### Documentation Agent
- Generates documentation using Claude AI
- Publishes to Confluence

Agents are automatically routed based on task tags.
