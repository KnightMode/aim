export interface Task {
  id: number;
  title: string;
  description: string;
  tags: string[];
  status: TaskStatus;
  priority: number;
  assigned_agent?: string;
  result?: string;
  error_msg?: string;
  metadata: TaskMetadata;
  created_at: string;
  updated_at: string;
  started_at?: string;
  completed_at?: string;
}

export type TaskStatus = 'todo' | 'queued' | 'in_progress' | 'completed' | 'failed';

export interface TaskMetadata {
  [key: string]: any;
  // Coding task metadata
  repo_url?: string;
  branch?: string;
  files_to_modify?: string[];
  instruction?: string;
  // Documentation task metadata
  confluence_space?: string;
  page_title?: string;
  content_source?: 'file' | 'url' | 'text';
  source_value?: string;
}

export interface ExecutionLog {
  id: number;
  task_id: number;
  agent_name: string;
  log_level: LogLevel;
  message: string;
  created_at: string;
}

export type LogLevel = 'info' | 'warning' | 'error' | 'success';

export interface CreateTaskRequest {
  title: string;
  description: string;
  tags: string[];
  priority: number;
  metadata: TaskMetadata;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  metadata?: TaskMetadata;
}

export interface UpdateStatusRequest {
  status: TaskStatus;
}

export interface AgentInfo {
  name: string;
  enabled: boolean;
  tags: string[];
}

export interface AgentStats {
  queued: number;
  in_progress: number;
  completed: number;
  failed: number;
  total: number;
}

export interface WSMessage {
  type: 'task_status_changed' | 'execution_log' | 'task_completed' | 'task_failed' | 'agent_started';
  task_id?: number;
  status?: TaskStatus;
  log_level?: LogLevel;
  message?: string;
  result?: string;
  error?: string;
  agent_name?: string;
  timestamp: string;
}
