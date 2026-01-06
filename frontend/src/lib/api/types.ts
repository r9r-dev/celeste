// System Stats (from WebSocket 'stats' message)
export interface SystemStats {
	cpuUsage: number;
	cpuCores: number;
	memoryUsed: number;
	memoryTotal: number;
	memoryPercent: number;
	diskUsed: number;
	diskTotal: number;
	diskPercent: number;
	uptime: number;
	hostname: string;
	os: string;
	platform: string;
}

// Container types
export interface PortBinding {
	hostIp: string;
	hostPort: number;
	containerPort: number;
	protocol: string;
}

export interface ContainerInfo {
	id: string;
	name: string;
	image: string;
	status: string;
	state: string;
	created: number;
	ports: PortBinding[];
	labels: Record<string, string>;
	networkMode: string;
}

export interface ContainerStats {
	cpuPercent: number;
	memoryUsage: number;
	memoryLimit: number;
	memoryPercent: number;
	networkRx: number;
	networkTx: number;
	blockRead: number;
	blockWrite: number;
}

// Stack types
export interface StackInfo {
	name: string;
	path: string;
	status: 'running' | 'partial' | 'stopped';
	services: number;
	runningServices: number;
}

// Volume types
export interface VolumeInfo {
	name: string;
	driver: string;
	mountpoint: string;
	createdAt: string;
	labels: Record<string, string>;
	usedBy: string[];
}

// Network types
export interface NetworkInfo {
	id: string;
	name: string;
	driver: string;
	scope: string;
	internal: boolean;
	containers: string[];
}

// Image types
export interface ImageInfo {
	id: string;
	tags: string[];
	size: number;
	created: number;
}

// WebSocket message types
export interface WebSocketMessage<T = unknown> {
	type: string;
	payload: T;
}

export interface ContainerStatsPayload {
	containers: Record<string, ContainerStats>;
	timestamp: number;
}

// API response types
export interface ApiError {
	error: string;
}

export interface HealthResponse {
	status: string;
	version: string;
}

export interface StatusResponse {
	status: string;
}

export interface ComposeFileResponse {
	content: string;
}

export interface LogsResponse {
	logs: string;
}
