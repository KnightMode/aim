# Quick Start Guide

Get the AI Task Manager up and running in 5 minutes!

## Prerequisites Check

Ensure you have:
- ✅ Node.js 18+ installed (`node --version`)
- ✅ Go 1.21+ installed (`go version`)
- ✅ MySQL 8.0+ installed and running

## Step 1: Database Setup (2 minutes)

```bash
# Start MySQL (if not running)
# macOS: brew services start mysql
# Linux: sudo systemctl start mysql

# Create database
mysql -u root -p
```

In MySQL shell:
```sql
CREATE DATABASE ai_task_manager;
EXIT;
```

## Step 2: Configure Backend (1 minute)

Edit `apps/backend/.env`:

```bash
# Required - Update these values:
DB_PASSWORD=your_mysql_password

# Get your Anthropic API key from: https://console.anthropic.com/
ANTHROPIC_API_KEY=sk-ant-api03-your-key-here

# Get GitHub token from: https://github.com/settings/tokens (select 'repo' scope)
GITHUB_TOKEN=ghp_your-github-token-here

GIT_USERNAME=your-github-username
GIT_EMAIL=your-email@example.com

# Optional (only for documentation tasks):
CONFLUENCE_URL=https://your-domain.atlassian.net/wiki
CONFLUENCE_USERNAME=your-email@example.com
CONFLUENCE_API_TOKEN=your-token
```

## Step 3: Install & Run (1 minute)

```bash
# Install all dependencies
npm install

# Run both backend and frontend
npm run dev
```

## Step 4: Access the Application

Open your browser:
- **Application**: http://localhost:5173
- **API Health**: http://localhost:8080/api/health

You should see the Kanban board interface!

## Step 5: Create Your First Task

### Coding Task Example:

1. Click **"+ New Task"**
2. Fill in:
   - **Title**: Add README file
   - **Description**: Create a basic README with project info
   - **Task Type**: Coding
   - **Priority**: Medium
   - **Repository URL**: https://github.com/YOUR_USERNAME/YOUR_REPO
   - **Branch**: main
   - **Files to Modify**: README.md
   - **Instruction**: Create a README.md file with project title, description, and installation instructions

3. Click **"Create Task"**

Watch it automatically:
- Move to "Queued"
- Move to "In Progress"
- Execute with AI
- Create a GitHub PR
- Move to "Completed" with the PR link!

## Troubleshooting

### Backend won't start?

```bash
# Check if MySQL is running
mysql -u root -p

# Verify environment variables
cat apps/backend/.env | grep -v "^#"

# Check backend logs
cd apps/backend
go run cmd/server/main.go
```

### Frontend won't load?

```bash
# Clear cache and reinstall
cd apps/frontend
rm -rf node_modules
npm install
npm run dev
```

### WebSocket not connecting?

Check the browser console (F12) for errors. Ensure:
- Backend is running on port 8080
- `FRONTEND_URL` in `.env` is `http://localhost:5173`

### Tasks not executing?

Verify your API keys:
```bash
# Check backend logs for errors
cd apps/backend
go run cmd/server/main.go

# You should see:
# - "Database connection established successfully"
# - "Starting task queue with 3 workers..."
# - "Agents registered"
```

## What's Next?

- **View Logs**: Click on any task to see real-time execution logs
- **Drag & Drop**: Move tasks between columns
- **Retry Failed**: Click "Retry" button on failed tasks
- **Create Documentation**: Create a task with "Documentation" tag

## Need More Help?

- **Detailed Setup**: See [SETUP.md](./SETUP.md)
- **Architecture**: See [README.md](./README.md)
- **Implementation Details**: See `.opencode/plan/plan.md`

## API Keys Reference

### Anthropic (Required for AI features)
1. Visit: https://console.anthropic.com/
2. Sign up/login
3. Create API key
4. Copy to `ANTHROPIC_API_KEY` in `.env`

### GitHub (Required for coding tasks)
1. Visit: https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scope: `repo` (Full control of private repositories)
4. Copy token to `GITHUB_TOKEN` in `.env`

### Confluence (Optional - for documentation tasks)
1. Visit: https://id.atlassian.com/manage-profile/security/api-tokens
2. Create API token
3. Copy to `CONFLUENCE_API_TOKEN` in `.env`
4. Add your Confluence URL and email

## Success Indicators

✅ Backend console shows:
```
Database connection established successfully
Database migrations completed successfully
Starting task queue with 3 workers...
Agents registered
Server starting on 0.0.0.0:8080
```

✅ Frontend shows:
- Kanban board with 5 columns
- "+ New Task" button in header
- No console errors

✅ Test task:
- Creates successfully
- Moves through statuses automatically
- Shows logs in real-time
- Completes with result (PR URL)

Enjoy building with AI! 🚀
