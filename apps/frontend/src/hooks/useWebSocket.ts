import { useEffect, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { wsService } from '../services/websocket';
import type { WSMessage } from '../types';

export const useWebSocket = () => {
  const queryClient = useQueryClient();
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    // Connect to WebSocket
    wsService.connect();

    // Subscribe to connection status changes
    const unsubscribeStatus = wsService.subscribeToConnectionStatus((connected) => {
      setIsConnected(connected);
    });

    // Subscribe to messages
    const unsubscribe = wsService.subscribe((message: WSMessage) => {
      console.log('WebSocket message:', message);

      switch (message.type) {
        case 'task_status_changed':
        case 'task_completed':
        case 'task_failed':
        case 'agent_started':
          // Invalidate tasks query to refetch updated data
          queryClient.invalidateQueries({ queryKey: ['tasks'] });
          if (message.task_id) {
            queryClient.invalidateQueries({ queryKey: ['task', message.task_id] });
            queryClient.invalidateQueries({ queryKey: ['logs', message.task_id] });
          }
          break;

        case 'execution_log':
          // Invalidate logs query for the specific task
          if (message.task_id) {
            queryClient.invalidateQueries({ queryKey: ['logs', message.task_id] });
          }
          break;

        default:
          break;
      }
    });

    // Cleanup on unmount
    return () => {
      unsubscribe();
      unsubscribeStatus();
      wsService.disconnect();
    };
  }, [queryClient]);

  return {
    isConnected,
  };
};
