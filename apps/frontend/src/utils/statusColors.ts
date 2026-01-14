import type { TaskStatus, LogLevel } from '../types';

export const getStatusColor = (status: TaskStatus): string => {
  const colors: Record<TaskStatus, string> = {
    todo: 'bg-[#2d3441] text-[#8a93a6] border-[#2d3441]',
    queued: 'bg-[#ffb454]/10 text-[#ffb454] border-[#ffb454]',
    in_progress: 'bg-[#00d9ff]/10 text-[#00d9ff] border-[#00d9ff]',
    completed: 'bg-[#00ff88]/10 text-[#00ff88] border-[#00ff88]',
    failed: 'bg-[#ff006e]/10 text-[#ff006e] border-[#ff006e]',
  };
  return colors[status] || colors.todo;
};

// LED dot colors for status badges
export const getStatusLEDColor = (status: TaskStatus): string => {
  const colors: Record<TaskStatus, string> = {
    todo: '#5a6478',
    queued: '#ffb454',
    in_progress: '#00d9ff',
    completed: '#00ff88',
    failed: '#ff006e',
  };
  return colors[status] || colors.todo;
};

export const getLogLevelColor = (level: LogLevel): string => {
  const colors: Record<LogLevel, string> = {
    info: 'text-[#8a93a6]',
    warning: 'text-[#ffb454]',
    error: 'text-[#ff006e]',
    success: 'text-[#00ff88]',
  };
  return colors[level] || colors.info;
};

export const getLogLevelIcon = (level: LogLevel): string => {
  const icons: Record<LogLevel, string> = {
    info: '◆',
    warning: '⚠',
    error: '✕',
    success: '✓',
  };
  return icons[level] || icons.info;
};
