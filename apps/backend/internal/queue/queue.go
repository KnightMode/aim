package queue

import (
	"ai-task-manager/backend/internal/agents"
	"ai-task-manager/backend/internal/models"
	"ai-task-manager/backend/internal/websocket"
	"log"
	"time"

	"gorm.io/gorm"
)

// TaskQueue manages the task queue and worker pool
type TaskQueue struct {
	taskChannel chan *models.Task
	maxWorkers  int
	db          *gorm.DB
	wsHub       *websocket.Hub
	agentRouter *agents.Router
	stopChan    chan struct{}
	pollInterval time.Duration
}

// NewTaskQueue creates a new TaskQueue instance
func NewTaskQueue(
	maxWorkers int,
	pollInterval time.Duration,
	bufferSize int,
	db *gorm.DB,
	wsHub *websocket.Hub,
	router *agents.Router,
) *TaskQueue {
	return &TaskQueue{
		taskChannel: make(chan *models.Task, bufferSize),
		maxWorkers:  maxWorkers,
		db:          db,
		wsHub:       wsHub,
		agentRouter: router,
		stopChan:    make(chan struct{}),
		pollInterval: pollInterval,
	}
}

// Start initializes the queue and worker pool
func (q *TaskQueue) Start() {
	log.Printf("Starting task queue with %d workers...", q.maxWorkers)

	// Start polling for queued tasks
	go q.pollDatabaseForTasks()

	// Start worker pool
	for i := 0; i < q.maxWorkers; i++ {
		go q.worker(i)
	}

	log.Println("Task queue started successfully")
}

// pollDatabaseForTasks continuously polls the database for queued tasks
func (q *TaskQueue) pollDatabaseForTasks() {
	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Fetch queued tasks with priority ordering
			var tasks []models.Task
			q.db.Where("status = ?", models.StatusQueued).
				Order("priority DESC, created_at ASC").
				Limit(10).
				Find(&tasks)

			// Add to channel (non-blocking)
			for i := range tasks {
				select {
				case q.taskChannel <- &tasks[i]:
					log.Printf("Task %d added to queue", tasks[i].ID)
				default:
					// Channel full, will retry next poll
					log.Printf("Queue channel full, task %d will be retried", tasks[i].ID)
				}
			}

		case <-q.stopChan:
			log.Println("Stopping task polling...")
			return
		}
	}
}

// worker processes tasks from the queue
func (q *TaskQueue) worker(workerID int) {
	log.Printf("Worker %d started", workerID)

	for task := range q.taskChannel {
		log.Printf("Worker %d picked up task %d", workerID, task.ID)
		q.executeTask(task, workerID)
	}

	log.Printf("Worker %d stopped", workerID)
}

// executeTask executes a single task
func (q *TaskQueue) executeTask(task *models.Task, workerID int) {
	// Update status to in_progress
	now := time.Now().UTC()
	task.Status = models.StatusInProgress
	task.StartedAt = &now
	q.db.Save(task)

	// Broadcast status change
	q.wsHub.BroadcastTaskStatusChanged(task.ID, models.StatusInProgress, now.Format(time.RFC3339))

	// Get appropriate agent
	agent := q.agentRouter.GetAgentForTask(task)
	if agent == nil {
		log.Printf("No agent found for task %d", task.ID)
		q.markTaskFailed(task, "No agent available for this task type")
		return
	}

	task.AssignedAgent = agent.Name()
	q.db.Save(task)

	// Broadcast agent started
	q.wsHub.BroadcastAgentStarted(task.ID, agent.Name(), time.Now().UTC().Format(time.RFC3339))

	log.Printf("Worker %d: Agent %s executing task %d", workerID, agent.Name(), task.ID)

	// Execute task
	err := agent.Execute(task)

	// Update final status
	completedTime := time.Now().UTC()
	task.CompletedAt = &completedTime

	if err != nil {
		log.Printf("Worker %d: Task %d failed: %v", workerID, task.ID, err)
		task.Status = models.StatusFailed
		task.ErrorMsg = err.Error()
		q.wsHub.BroadcastTaskFailed(task.ID, err.Error(), completedTime.Format(time.RFC3339))
	} else {
		log.Printf("Worker %d: Task %d completed successfully", workerID, task.ID)
		task.Status = models.StatusCompleted
		q.wsHub.BroadcastTaskCompleted(task.ID, task.Result, completedTime.Format(time.RFC3339))
	}

	q.db.Save(task)
}

// markTaskFailed marks a task as failed
func (q *TaskQueue) markTaskFailed(task *models.Task, errorMsg string) {
	now := time.Now().UTC()
	task.Status = models.StatusFailed
	task.ErrorMsg = errorMsg
	task.CompletedAt = &now
	q.db.Save(task)
	q.wsHub.BroadcastTaskFailed(task.ID, errorMsg, now.Format(time.RFC3339))
}

// Stop gracefully shuts down the queue
func (q *TaskQueue) Stop() {
	log.Println("Stopping task queue...")
	close(q.stopChan)
	close(q.taskChannel)
	log.Println("Task queue stopped")
}

// EnqueueTask adds a task to the queue by ID
func (q *TaskQueue) EnqueueTask(taskID uint) error {
	var task models.Task
	if err := q.db.First(&task, taskID).Error; err != nil {
		return err
	}

	task.Status = models.StatusQueued
	q.db.Save(&task)

	return nil
}
