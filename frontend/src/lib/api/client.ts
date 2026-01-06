import type {
	SystemStats,
	ContainerInfo,
	ContainerStats,
	StackInfo,
	VolumeInfo,
	NetworkInfo,
	ImageInfo,
	StatusResponse,
	HealthResponse,
	ComposeFileResponse,
	LogsResponse,
	ApiError as ApiErrorResponse
} from './types';

// Determine API base URL based on environment
const getApiBase = (): string => {
	if (typeof window === 'undefined') return 'http://localhost:8080/api';
	return `${window.location.origin}/api`;
};

const API_BASE = getApiBase();

class ApiError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
	const url = `${API_BASE}${endpoint}`;
	const response = await fetch(url, {
		...options,
		headers: {
			'Content-Type': 'application/json',
			...options.headers
		}
	});

	if (!response.ok) {
		const errorResponse = (await response.json().catch(() => ({ error: 'Unknown error' }))) as ApiErrorResponse;
		throw new ApiError(response.status, errorResponse.error || `HTTP ${response.status}`);
	}

	return response.json();
}

// System
export async function getSystemStats(): Promise<SystemStats> {
	return request<SystemStats>('/stats');
}

export async function getHealth(): Promise<HealthResponse> {
	const baseUrl = typeof window === 'undefined' ? 'http://localhost:8080' : window.location.origin;
	const response = await fetch(`${baseUrl}/health`);
	return response.json();
}

// Stacks
export async function listStacks(): Promise<StackInfo[]> {
	return request<StackInfo[]>('/stacks');
}

export async function getStack(name: string): Promise<StackInfo> {
	return request<StackInfo>(`/stacks/${encodeURIComponent(name)}`);
}

export async function startStack(name: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/stacks/${encodeURIComponent(name)}/start`, {
		method: 'POST'
	});
}

export async function stopStack(name: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/stacks/${encodeURIComponent(name)}/stop`, {
		method: 'POST'
	});
}

export async function restartStack(name: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/stacks/${encodeURIComponent(name)}/restart`, {
		method: 'POST'
	});
}

export async function pullStack(name: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/stacks/${encodeURIComponent(name)}/pull`, {
		method: 'POST'
	});
}

export async function getComposeFile(name: string): Promise<string> {
	const response = await request<ComposeFileResponse>(
		`/stacks/${encodeURIComponent(name)}/compose`
	);
	return response.content;
}

export async function updateComposeFile(name: string, content: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/stacks/${encodeURIComponent(name)}/compose`, {
		method: 'PUT',
		body: JSON.stringify({ content })
	});
}

// Containers
export async function listContainers(all = true): Promise<ContainerInfo[]> {
	return request<ContainerInfo[]>(`/containers?all=${all}`);
}

export async function getContainer(id: string): Promise<ContainerInfo> {
	return request<ContainerInfo>(`/containers/${encodeURIComponent(id)}`);
}

export async function startContainer(id: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/containers/${encodeURIComponent(id)}/start`, {
		method: 'POST'
	});
}

export async function stopContainer(id: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/containers/${encodeURIComponent(id)}/stop`, {
		method: 'POST'
	});
}

export async function restartContainer(id: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/containers/${encodeURIComponent(id)}/restart`, {
		method: 'POST'
	});
}

export async function getContainerLogs(id: string, tail = '100'): Promise<string> {
	const response = await request<LogsResponse>(
		`/containers/${encodeURIComponent(id)}/logs?tail=${tail}`
	);
	return response.logs;
}

export async function getContainerStats(id: string): Promise<ContainerStats> {
	return request<ContainerStats>(`/containers/${encodeURIComponent(id)}/stats`);
}

// Volumes
export async function listVolumes(): Promise<VolumeInfo[]> {
	return request<VolumeInfo[]>('/volumes');
}

export async function createVolume(
	name: string,
	driver = 'local',
	labels: Record<string, string> = {}
): Promise<VolumeInfo> {
	return request<VolumeInfo>('/volumes', {
		method: 'POST',
		body: JSON.stringify({ name, driver, labels })
	});
}

export async function deleteVolume(name: string, force = false): Promise<StatusResponse> {
	return request<StatusResponse>(`/volumes/${encodeURIComponent(name)}?force=${force}`, {
		method: 'DELETE'
	});
}

// Networks
export async function listNetworks(): Promise<NetworkInfo[]> {
	return request<NetworkInfo[]>('/networks');
}

export async function createNetwork(name: string, driver = 'bridge'): Promise<NetworkInfo> {
	return request<NetworkInfo>('/networks', {
		method: 'POST',
		body: JSON.stringify({ name, driver })
	});
}

export async function deleteNetwork(id: string): Promise<StatusResponse> {
	return request<StatusResponse>(`/networks/${encodeURIComponent(id)}`, {
		method: 'DELETE'
	});
}

// Images
export async function listImages(): Promise<ImageInfo[]> {
	return request<ImageInfo[]>('/images');
}

// Export all functions as api object for convenience
export const api = {
	getSystemStats,
	getHealth,
	listStacks,
	getStack,
	startStack,
	stopStack,
	restartStack,
	pullStack,
	getComposeFile,
	updateComposeFile,
	listContainers,
	getContainer,
	startContainer,
	stopContainer,
	restartContainer,
	getContainerLogs,
	getContainerStats,
	listVolumes,
	createVolume,
	deleteVolume,
	listNetworks,
	createNetwork,
	deleteNetwork,
	listImages
};
