import type { WebSocketMessage } from './types';

type MessageHandler<T = unknown> = (payload: T) => void;

// Determine WebSocket URL based on environment
const getWsUrl = (): string => {
	if (typeof window === 'undefined') return 'ws://localhost:8080/ws';
	const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	return `${protocol}//${window.location.host}/ws`;
};

const WS_URL = getWsUrl();
const RECONNECT_DELAY = 1000;
const MAX_RECONNECT_DELAY = 30000;
const MAX_RECONNECT_ATTEMPTS = 10;

class WebSocketClient {
	private ws: WebSocket | null = null;
	private handlers: Map<string, Set<MessageHandler>> = new Map();
	private reconnectAttempts = 0;
	private reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
	private url: string;
	private _connected = $state(false);

	constructor(url: string = WS_URL) {
		this.url = url;
	}

	get connected(): boolean {
		return this._connected;
	}

	connect(): void {
		if (this.ws?.readyState === WebSocket.OPEN) {
			return;
		}

		try {
			this.ws = new WebSocket(this.url);

			this.ws.onopen = () => {
				console.log('[WS] Connected');
				this._connected = true;
				this.reconnectAttempts = 0;
			};

			this.ws.onclose = (event) => {
				console.log('[WS] Disconnected:', event.code, event.reason);
				this._connected = false;
				this.ws = null;
				this.scheduleReconnect();
			};

			this.ws.onerror = (error) => {
				console.error('[WS] Error:', error);
			};

			this.ws.onmessage = (event) => {
				this.handleMessage(event);
			};
		} catch (error) {
			console.error('[WS] Connection error:', error);
			this.scheduleReconnect();
		}
	}

	disconnect(): void {
		if (this.reconnectTimeout) {
			clearTimeout(this.reconnectTimeout);
			this.reconnectTimeout = null;
		}
		if (this.ws) {
			this.ws.close();
			this.ws = null;
		}
		this._connected = false;
		this.reconnectAttempts = 0;
	}

	on<T>(type: string, handler: MessageHandler<T>): () => void {
		if (!this.handlers.has(type)) {
			this.handlers.set(type, new Set());
		}
		this.handlers.get(type)!.add(handler as MessageHandler);

		// Return unsubscribe function
		return () => {
			const typeHandlers = this.handlers.get(type);
			if (typeHandlers) {
				typeHandlers.delete(handler as MessageHandler);
				if (typeHandlers.size === 0) {
					this.handlers.delete(type);
				}
			}
		};
	}

	off(type: string, handler?: MessageHandler): void {
		if (handler) {
			const typeHandlers = this.handlers.get(type);
			if (typeHandlers) {
				typeHandlers.delete(handler);
			}
		} else {
			this.handlers.delete(type);
		}
	}

	private handleMessage(event: MessageEvent): void {
		try {
			const message: WebSocketMessage = JSON.parse(event.data);
			const handlers = this.handlers.get(message.type);

			if (handlers) {
				handlers.forEach((handler) => {
					try {
						handler(message.payload);
					} catch (error) {
						console.error(`[WS] Handler error for type "${message.type}":`, error);
					}
				});
			}
		} catch (error) {
			console.error('[WS] Failed to parse message:', error);
		}
	}

	private scheduleReconnect(): void {
		if (this.reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
			console.error('[WS] Max reconnection attempts reached');
			return;
		}

		const delay = Math.min(
			RECONNECT_DELAY * Math.pow(2, this.reconnectAttempts),
			MAX_RECONNECT_DELAY
		);

		console.log(`[WS] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts + 1})`);

		this.reconnectTimeout = setTimeout(() => {
			this.reconnectAttempts++;
			this.connect();
		}, delay);
	}
}

// Export singleton instance
export const wsClient = new WebSocketClient();

// Export class for testing or custom instances
export { WebSocketClient };
