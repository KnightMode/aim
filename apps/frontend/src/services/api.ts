import axios from 'axios';
import type {
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
  UpdateStatusRequest,
  ExecutionLog,
  AgentInfo,
  AgentStats,
} from '../types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

// Task API
export const taskAPI = {
  getAll: async (params?: { status?: string; tag?: string; limit?: number }) => {
    const response = await api.get<{ tasks: Task[] }>('/tasks', { params });
    return response.data.tasks;
  },

  getById: async (id: number) => {
    const response = await api.get<Task>(`/tasks/${id}`);
    return response.data;
  },

  create: async (data: CreateTaskRequest) => {
    const response = await api.post<Task>('/tasks', data);
    return response.data;
  },

  update: async (id: number, data: UpdateTaskRequest) => {
    const response = await api.put<Task>(`/tasks/${id}`, data);
    return response.data;
  },

  delete: async (id: number) => {
    await api.delete(`/tasks/${id}`);
  },

  updateStatus: async (id: number, data: UpdateStatusRequest) => {
    const response = await api.patch<Task>(`/tasks/${id}/status`, data);
    return response.data;
  },
};

// Log API
export const logAPI = {
  getTaskLogs: async (taskId: number) => {
    const response = await api.get<{ logs: ExecutionLog[] }>(`/tasks/${taskId}/logs`);
    return response.data.logs;
  },

  getRecentLogs: async (limit?: number) => {
    const response = await api.get<{ logs: ExecutionLog[] }>('/logs/recent', {
      params: { limit },
    });
    return response.data.logs;
  },
};

// Agent API
export const agentAPI = {
  getAll: async () => {
    const response = await api.get<{ agents: AgentInfo[] }>('/agents');
    return response.data.agents;
  },

  getStats: async () => {
    const response = await api.get<AgentStats>('/agents/stats');
    return response.data;
  },
};

export default api;
