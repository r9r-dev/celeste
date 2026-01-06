import { wsClient } from '$lib/api/websocket';
import type { SystemStats } from '$lib/api/types';

// Default stats for initial state
const defaultStats: SystemStats = {
	cpuUsage: 0,
	cpuCores: 0,
	memoryUsed: 0,
	memoryTotal: 0,
	memoryPercent: 0,
	diskUsed: 0,
	diskTotal: 0,
	diskPercent: 0,
	uptime: 0,
	hostname: '',
	os: '',
	platform: ''
};

function createSystemStore() {
	let stats = $state<SystemStats>(defaultStats);
	let connected = $state(false);
	let error = $state<string | null>(null);
	let unsubscribe: (() => void) | null = null;

	function connect() {
		wsClient.connect();

		// Subscribe to system stats
		unsubscribe = wsClient.on<SystemStats>('stats', (payload) => {
			stats = payload;
			error = null;
		});

		// Track connection state
		$effect(() => {
			connected = wsClient.connected;
		});
	}

	function disconnect() {
		if (unsubscribe) {
			unsubscribe();
			unsubscribe = null;
		}
		wsClient.disconnect();
		connected = false;
	}

	return {
		get stats() {
			return stats;
		},
		get connected() {
			return connected;
		},
		get error() {
			return error;
		},
		connect,
		disconnect
	};
}

export const systemStore = createSystemStore();
