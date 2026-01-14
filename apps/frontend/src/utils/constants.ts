export const TASK_STATUSES = {
  TODO: 'todo',
  QUEUED: 'queued',
  IN_PROGRESS: 'in_progress',
  COMPLETED: 'completed',
  FAILED: 'failed',
} as const;

export const TASK_TAGS = {
  CODING: 'coding',
  DOCUMENTATION: 'documentation',
} as const;

export const LOG_LEVELS = {
  INFO: 'info',
  WARNING: 'warning',
  ERROR: 'error',
  SUCCESS: 'success',
} as const;

export const STATUS_LABELS = {
  [TASK_STATUSES.TODO]: 'To Do',
  [TASK_STATUSES.QUEUED]: 'Queued',
  [TASK_STATUSES.IN_PROGRESS]: 'In Progress',
  [TASK_STATUSES.COMPLETED]: 'Completed',
  [TASK_STATUSES.FAILED]: 'Failed',
};

export const TAG_LABELS = {
  [TASK_TAGS.CODING]: 'Coding',
  [TASK_TAGS.DOCUMENTATION]: 'Documentation',
};

export const PRIORITY_OPTIONS = [
  { value: 0, label: 'Low' },
  { value: 1, label: 'Medium' },
  { value: 2, label: 'High' },
];
