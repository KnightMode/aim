# AI Task Manager - Setup Guide

Complete step-by-step guide to set up and run the AI Task Manager application.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Node.js** 18.0.0 or higher
- **npm** 10.0.0 or higher
- **Go** 1.21 or higher
- **MySQL** 8.0 or higher

### Required API Keys

You'll need to obtain the following API keys:

1. **Anthropic API Key**: Sign up at [Anthropic](https://console.anthropic.com/) and create an API key
2. **GitHub Personal Access Token**: Create at [GitHub Settings → Developer settings → Personal access tokens](https://github.com/settings/tokens)
   - Required scope: `repo` (Full control of private repositories)
3. **Confluence API Token** (Optional): Create at [Atlassian Account Settings](https://id.atlassian.com/manage-profile/security/api-tokens)

## Installation Steps

### 1. Clone the Repository

```bash
git clone <repository-url>
cd ai-task-manager
```

### 2. Install Root Dependencies

```bash
npm install
```

This installs Turborepo for managing the monorepo.

### 3. Set Up MySQL Database

Start MySQL and create the database:

```bash
mysql -u root -p
```

In the MySQL shell:

```sql
CREATE DATABASE ai_task_manager;
EXIT;
```

### 4. Configure Backend

Navigate to the backend directory and set up configuration:

```bash
cd apps/backend
cp .env.example .env
```

Edit `.env` and fill in your configuration:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
FRONTEND_URL=http://localhost:5173

# Database Configuration (update with your MySQL credentials)
DB_HOST=localhost
DB_PORT=3306
DB_NAME=ai_task_manager
DB_USER=root
DB_PASSWORD=your_mysql_password_here
DB_CHARSET=utf8mb4

# Task Queue Configuration
WORKER_COUNT=3
QUEUE_POLL_INTERVAL=2s
QUEUE_BUFFER_SIZE=100

# AI Provider (Anthropic Claude)
ANTHROPIC_API_KEY=sk-ant-api03-your-key-here
ANTHROPIC_MODEL=claude-4-5-sonnet-20241022
ANTHROPIC_MAX_TOKENS=4096

# Git & GitHub Configuration
GIT_USERNAME=your-github-username
GIT_EMAIL=your-email@example.com
GITHUB_TOKEN=ghp_your-github-token-here

# Confluence Configuration (Optional - for documentation tasks)
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

Install Go dependencies:

```bash
go mod download
```

### 5. Configure Frontend

Navigate to the frontend directory:

```bash
cd ../frontend
```

The `.env.local` file has already been created with default values. If you need to change the API URLs, edit:

```bash
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080/ws
```

Install frontend dependencies:

```bash
npm install
```

## Running the Application

### Option 1: Run Everything from Root (Recommended)

From the root directory of the project:

```bash
# Start both backend and frontend
npm run dev
```

This uses Turborepo to run both applications in parallel.

### Option 2: Run Backend and Frontend Separately

**Terminal 1 - Backend:**
```bash
cd apps/backend
go run cmd/server/main.go
```

**Terminal 2 - Frontend:**
```bash
cd apps/frontend
npm run dev
```

### Option 3: Run Individual Services

**Backend only:**
```bash
npm run backend:dev
```

**Frontend only:**
```bash
npm run frontend:dev
```

## Accessing the Application

Once both services are running:

- **Frontend UI**: http://localhost:5173
- **Backend API**: http://localhost:8080/api
- **WebSocket**: ws://localhost:8080/ws
- **Health Check**: http://localhost:8080/api/health

## Verification Steps

### 1. Check Backend is Running

Visit: http://localhost:8080/api/health

You should see:
```json
{
  "status": "ok",
  "queue_size": 0,
  "in_progress": 0
}
```

### 2. Check Frontend is Running

Visit: http://localhost:5173

You should see the AI Task Manager interface with a Kanban board.

### 3. Check Database Connection

In the backend terminal, you should see:
```
Database connection established successfully
Database migrations completed successfully
```

### 4. Check WebSocket Connection

In the browser console (F12), you should see:
```
WebSocket connected
```

## Creating Your First Task

### Coding Task Example

1. Click **"+ New Task"** button
2. Fill in the form:
   - **Title**: Add dark mode toggle
   - **Description**: Implement dark mode with localStorage persistence
   - **Task Type**: Coding
   - **Priority**: Medium
   - **Repository URL**: https://github.com/yourusername/your-repo
   - **Branch**: main
   - **Files to Modify**: src/components/Settings.tsx
   - **Instruction**: Add a dark mode toggle button with persistence

3. Click **"Create Task"**

The task will:
- Appear in the "To Do" column
- Automatically move to "Queued"
- Move to "In Progress" when picked up by an agent
- Move to "Completed" with a PR URL (or "Failed" if there's an error)

### Documentation Task Example

1. Click **"+ New Task"** button
2. Fill in the form:
   - **Title**: Update API documentation
   - **Description**: Generate docs from code
   - **Task Type**: Documentation
   - **Priority**: Medium
   - **Confluence Space**: DEV
   - **Page Title**: API Documentation
   - **Content Source**: File or Text
   - **Source Value**: (file path or content)
   - **Instruction**: Generate comprehensive API documentation

3. Click **"Create Task"**

## Troubleshooting

### Backend Won't Start

**Error: Failed to connect to database**
- Check MySQL is running: `mysql -u root -p`
- Verify database exists: `SHOW DATABASES;`
- Check credentials in `.env` file

**Error: ANTHROPIC_API_KEY is required**
- Ensure you've set your Anthropic API key in `.env`

**Error: GITHUB_TOKEN is required**
- Ensure you've set your GitHub token in `.env`

### Frontend Won't Build

**Error: Cannot find module**
- Delete `node_modules` and reinstall: `rm -rf node_modules && npm install`

**Error: Failed to fetch tasks**
- Ensure backend is running on port 8080
- Check CORS settings in backend

### WebSocket Not Connecting

**Issue: Real-time updates not working**
- Check backend logs for WebSocket connection messages
- Verify `FRONTEND_URL` in backend `.env` matches frontend URL
- Check browser console for connection errors

### Tasks Stuck in "Queued"

**Issue: Tasks not executing**
- Check backend logs for worker errors
- Verify `WORKER_COUNT` > 0 in `.env`
- Check task queue logs in backend terminal

### Git/GitHub Errors

**Error: Failed to clone repository**
- Verify repository URL is correct and accessible
- Check GitHub token has `repo` scope
- Ensure you have access to the repository

**Error: Failed to push**
- Verify Git username and email in `.env`
- Check GitHub token permissions

## Building for Production

### Backend

```bash
cd apps/backend
go build -o bin/server cmd/server/main.go
./bin/server
```

### Frontend

```bash
cd apps/frontend
npm run build
```

The built files will be in `dist/` directory. Serve with any static file server.

## Next Steps

- Read the [README.md](./README.md) for architecture overview
- Check the [plan](./opencode/plan/plan.md) for implementation details
- Create your first task and watch the AI agent execute it!
- View execution logs in real-time by clicking on any task

## Support

For issues, please check:
1. Backend logs in the terminal
2. Frontend console (F12 in browser)
3. MySQL connection and database status
4. API key validity and permissions

## Security Notes

- Never commit `.env` files to version control
- Keep your API keys secure
- Use environment-specific configurations for production
- Review all AI-generated code in PRs before merging
