import type { WSMessage } from '../types';

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';

export type WSEventHandler = (message: WSMessage) => void;
export type ConnectionStatusHandler = (isConnected: boolean) => void;

class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;
  private handlers: Set<WSEventHandler> = new Set();
  private connectionStatusHandlers: Set<ConnectionStatusHandler> = new Set();
  private isConnecting = false;

  connect() {
    if (this.ws?.readyState === WebSocket.OPEN || this.isConnecting) {
      return;
    }

    this.isConnecting = true;

    try {
      this.ws = new WebSocket(WS_URL);

      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        this.isConnecting = false;
        this.notifyConnectionStatus(true);
      };

      this.ws.onmessage = (event) => {
        try {
          const messages = event.data.split('\n').filter(Boolean);
          messages.forEach((msgStr: string) => {
            const message: WSMessage = JSON.parse(msgStr);
            this.handlers.forEach((handler) => handler(message));
          });
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.isConnecting = false;
      };

      this.ws.onclose = () => {
        console.log('WebSocket closed');
        this.isConnecting = false;
        this.notifyConnectionStatus(false);
        this.reconnect();
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      this.isConnecting = false;
      this.reconnect();
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  private reconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    this.reconnectAttempts++;
    console.log(`Reconnecting... Attempt ${this.reconnectAttempts}`);

    setTimeout(() => {
      this.connect();
    }, this.reconnectDelay);
  }

  subscribe(handler: WSEventHandler) {
    this.handlers.add(handler);
    return () => {
      this.handlers.delete(handler);
    };
  }

  subscribeToConnectionStatus(handler: ConnectionStatusHandler) {
    this.connectionStatusHandlers.add(handler);
    // Immediately notify with current status
    handler(this.isConnected());
    return () => {
      this.connectionStatusHandlers.delete(handler);
    };
  }

  private notifyConnectionStatus(isConnected: boolean) {
    this.connectionStatusHandlers.forEach((handler) => handler(isConnected));
  }

  isConnected() {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}

export const wsService = new WebSocketService();
