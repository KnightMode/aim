import { useQuery } from '@tanstack/react-query';
import { logAPI } from '../services/api';

export const useTaskLogs = (taskId: number) => {
  return useQuery({
    queryKey: ['logs', taskId],
    queryFn: () => logAPI.getTaskLogs(taskId),
    enabled: !!taskId,
    refetchInterval: 2000, // Refetch every 2 seconds for real-time updates
  });
};

export const useRecentLogs = (limit?: number) => {
  return useQuery({
    queryKey: ['logs', 'recent', limit],
    queryFn: () => logAPI.getRecentLogs(limit),
  });
};
