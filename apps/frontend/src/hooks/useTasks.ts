import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { taskAPI } from '../services/api';
import type { CreateTaskRequest, UpdateTaskRequest, UpdateStatusRequest } from '../types';

export const useTasks = (params?: { status?: string; tag?: string; limit?: number }) => {
  return useQuery({
    queryKey: ['tasks', params],
    queryFn: () => taskAPI.getAll(params),
    refetchInterval: 5000, // Refetch every 5 seconds as backup to WebSocket
  });
};

export const useTask = (id: number) => {
  return useQuery({
    queryKey: ['task', id],
    queryFn: () => taskAPI.getById(id),
    enabled: !!id,
  });
};

export const useCreateTask = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateTaskRequest) => taskAPI.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
  });
};

export const useUpdateTask = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateTaskRequest }) =>
      taskAPI.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      queryClient.invalidateQueries({ queryKey: ['task', variables.id] });
    },
  });
};

export const useDeleteTask = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => taskAPI.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
    },
  });
};

export const useUpdateTaskStatus = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateStatusRequest }) =>
      taskAPI.updateStatus(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] });
      queryClient.invalidateQueries({ queryKey: ['task', variables.id] });
    },
  });
};
