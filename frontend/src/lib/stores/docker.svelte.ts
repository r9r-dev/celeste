import { api } from '$lib/api/client';
import { wsClient } from '$lib/api/websocket.svelte';
import type {
	ContainerInfo,
	ContainerStats,
	StackInfo,
	VolumeInfo,
	NetworkInfo,
	ImageInfo,
	ContainerStatsPayload
} from '$lib/api/types';

function createDockerStore() {
	let containers = $state<ContainerInfo[]>([]);
	let stacks = $state<StackInfo[]>([]);
	let volumes = $state<VolumeInfo[]>([]);
	let networks = $state<NetworkInfo[]>([]);
	let images = $state<ImageInfo[]>([]);
	let containerStats = $state<Map<string, ContainerStats>>(new Map());
	let loading = $state(false);
	let error = $state<string | null>(null);
	let unsubscribe: (() => void) | null = null;

	// Fetch functions
	async function fetchContainers(all = true) {
		loading = true;
		error = null;
		try {
			containers = await api.listContainers(all);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to fetch containers';
			console.error('[Docker] Error fetching containers:', e);
		} finally {
			loading = false;
		}
	}

	async function fetchStacks() {
		loading = true;
		error = null;
		try {
			stacks = await api.listStacks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to fetch stacks';
			console.error('[Docker] Error fetching stacks:', e);
		} finally {
			loading = false;
		}
	}

	async function fetchVolumes() {
		try {
			volumes = await api.listVolumes();
		} catch (e) {
			console.error('[Docker] Error fetching volumes:', e);
		}
	}

	async function fetchNetworks() {
		try {
			networks = await api.listNetworks();
		} catch (e) {
			console.error('[Docker] Error fetching networks:', e);
		}
	}

	async function fetchImages() {
		try {
			images = await api.listImages();
		} catch (e) {
			console.error('[Docker] Error fetching images:', e);
		}
	}

	async function fetchAll() {
		loading = true;
		error = null;
		try {
			await Promise.all([
				fetchContainers(),
				fetchStacks(),
				fetchVolumes(),
				fetchNetworks(),
				fetchImages()
			]);
		} finally {
			loading = false;
		}
	}

	// Subscribe to container stats from WebSocket
	function subscribeToContainerStats() {
		unsubscribe = wsClient.on<ContainerStatsPayload>('container_stats', (payload) => {
			containerStats = new Map(Object.entries(payload.containers));
		});
	}

	function unsubscribeFromContainerStats() {
		if (unsubscribe) {
			unsubscribe();
			unsubscribe = null;
		}
	}

	// Stack operations
	async function startStack(name: string) {
		try {
			await api.startStack(name);
			await fetchStacks();
			await fetchContainers();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to start stack';
			throw e;
		}
	}

	async function stopStack(name: string) {
		try {
			await api.stopStack(name);
			await fetchStacks();
			await fetchContainers();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to stop stack';
			throw e;
		}
	}

	async function restartStack(name: string) {
		try {
			await api.restartStack(name);
			await fetchStacks();
			await fetchContainers();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to restart stack';
			throw e;
		}
	}

	async function pullStack(name: string) {
		try {
			await api.pullStack(name);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to pull stack';
			throw e;
		}
	}

	// Container operations
	async function startContainer(id: string) {
		try {
			await api.startContainer(id);
			await fetchContainers();
			await fetchStacks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to start container';
			throw e;
		}
	}

	async function stopContainer(id: string) {
		try {
			await api.stopContainer(id);
			await fetchContainers();
			await fetchStacks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to stop container';
			throw e;
		}
	}

	async function restartContainer(id: string) {
		try {
			await api.restartContainer(id);
			await fetchContainers();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to restart container';
			throw e;
		}
	}

	// Volume operations
	async function createVolume(name: string, driver = 'local', labels: Record<string, string> = {}) {
		try {
			await api.createVolume(name, driver, labels);
			await fetchVolumes();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create volume';
			throw e;
		}
	}

	async function deleteVolume(name: string, force = false) {
		try {
			await api.deleteVolume(name, force);
			await fetchVolumes();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to delete volume';
			throw e;
		}
	}

	// Network operations
	async function createNetwork(name: string, driver = 'bridge') {
		try {
			await api.createNetwork(name, driver);
			await fetchNetworks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to create network';
			throw e;
		}
	}

	async function deleteNetwork(id: string) {
		try {
			await api.deleteNetwork(id);
			await fetchNetworks();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to delete network';
			throw e;
		}
	}

	// Helper: Get containers for a specific stack
	function getStackContainers(stackName: string): ContainerInfo[] {
		return containers.filter((c) => c.labels['com.docker.compose.project'] === stackName);
	}

	// Helper: Get stats for a container
	function getContainerStatsById(containerId: string): ContainerStats | undefined {
		return containerStats.get(containerId);
	}

	return {
		// State getters
		get containers() {
			return containers;
		},
		get stacks() {
			return stacks;
		},
		get volumes() {
			return volumes;
		},
		get networks() {
			return networks;
		},
		get images() {
			return images;
		},
		get containerStats() {
			return containerStats;
		},
		get loading() {
			return loading;
		},
		get error() {
			return error;
		},

		// Fetch functions
		fetchContainers,
		fetchStacks,
		fetchVolumes,
		fetchNetworks,
		fetchImages,
		fetchAll,

		// WebSocket subscriptions
		subscribeToContainerStats,
		unsubscribeFromContainerStats,

		// Stack operations
		startStack,
		stopStack,
		restartStack,
		pullStack,

		// Container operations
		startContainer,
		stopContainer,
		restartContainer,

		// Volume operations
		createVolume,
		deleteVolume,

		// Network operations
		createNetwork,
		deleteNetwork,

		// Helpers
		getStackContainers,
		getContainerStatsById
	};
}

export const dockerStore = createDockerStore();
