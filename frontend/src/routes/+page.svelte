<script lang="ts">
	import { onMount, untrack } from 'svelte';
	import { Panel, ProgressBar, StatusIndicator, DataLabel, LineChart, Button } from '$lib/components/ui';
	import { systemStore } from '$lib/stores/system.svelte';
	import { dockerStore } from '$lib/stores/docker.svelte';
	import type { StackInfo, ContainerInfo, ContainerStats } from '$lib/api/types';

	// Extended stack type for the UI with containers
	interface StackWithContainers extends StackInfo {
		containers: ContainerWithStats[];
	}

	interface ContainerWithStats extends ContainerInfo {
		cpu: number;
		memory: number;
		memoryLimit: number;
	}

	// Connect to WebSocket and fetch initial data
	onMount(() => {
		systemStore.connect();
		dockerStore.subscribeToContainerStats();
		dockerStore.fetchAll();

		// Refresh data periodically
		const refreshInterval = setInterval(() => {
			dockerStore.fetchContainers();
			dockerStore.fetchStacks();
		}, 10000);

		return () => {
			systemStore.disconnect();
			dockerStore.unsubscribeFromContainerStats();
			clearInterval(refreshInterval);
		};
	});

	// Build stacks with their containers and stats
	const stacksWithContainers = $derived.by(() => {
		return dockerStore.stacks.map((stack): StackWithContainers => {
			const stackContainers = dockerStore.containers
				.filter(c => c.labels['com.docker.compose.project'] === stack.name)
				.map((container): ContainerWithStats => {
					const stats = dockerStore.containerStats.get(container.id);
					return {
						...container,
						cpu: stats?.cpuPercent ?? 0,
						memory: stats?.memoryUsage ? stats.memoryUsage / (1024 * 1024) : 0,
						memoryLimit: stats?.memoryLimit ? stats.memoryLimit / (1024 * 1024) : 0
					};
				});

			return {
				...stack,
				containers: stackContainers
			};
		});
	});

	// Generate network history for chart (will be replaced with real data later)
	let networkHistory = $state(Array.from({ length: 40 }, () => Math.random() * 100 + 50));

	// Update network history periodically
	$effect(() => {
		const interval = setInterval(() => {
			networkHistory = [...networkHistory.slice(1), Math.random() * 100 + 50];
		}, 2000);
		return () => clearInterval(interval);
	});

	let selectedStack = $state<StackWithContainers | null>(null);

	// Auto-select first stack when data loads
	// Use untrack to read selectedStack without tracking it as a dependency
	// This prevents an infinite loop (reading and writing the same state in an effect)
	$effect(() => {
		const current = untrack(() => selectedStack);
		if (!current && stacksWithContainers.length > 0) {
			selectedStack = stacksWithContainers[0];
		} else if (current) {
			// Update selected stack with new data
			const updated = stacksWithContainers.find(s => s.name === current.name);
			if (updated) {
				selectedStack = updated;
			}
		}
	});

	const formatUptime = (seconds: number) => {
		const days = Math.floor(seconds / 86400);
		const hours = Math.floor((seconds % 86400) / 3600);
		const mins = Math.floor((seconds % 3600) / 60);
		return `${days}d ${hours}h ${mins}m`;
	};

	const formatBytes = (mb: number) => {
		if (mb >= 1024) return `${(mb / 1024).toFixed(1)} GB`;
		return `${mb.toFixed(0)} MB`;
	};

	const formatBytesRaw = (bytes: number) => {
		if (bytes >= 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB`;
		if (bytes >= 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
		return `${(bytes / 1024).toFixed(0)} KB`;
	};

	// Calculate total CPU for all running containers
	const totalCpu = $derived(
		stacksWithContainers.flatMap(s => s.containers)
			.filter(c => c.state === 'running')
			.reduce((sum, c) => sum + c.cpu, 0)
	);

	// Count running/stopped containers
	const runningContainers = $derived(
		dockerStore.containers.filter(c => c.state === 'running').length
	);
	const stoppedContainers = $derived(
		dockerStore.containers.filter(c => c.state !== 'running').length
	);

	// Stack actions
	async function handleStartStack(name: string) {
		await dockerStore.startStack(name);
	}

	async function handleStopStack(name: string) {
		await dockerStore.stopStack(name);
	}

	async function handleRestartStack(name: string) {
		await dockerStore.restartStack(name);
	}
</script>

<div class="dashboard">
	<!-- Loading indicator -->
	{#if dockerStore.loading}
		<div class="loading-overlay">
			<div class="loading-spinner"></div>
			<span>LOADING...</span>
		</div>
	{/if}

	<!-- Left Panel - Stack Visualization -->
	<div class="visualization-panel">
		<Panel title="INFRASTRUCTURE MAP" expandable variant="accent">
			<div class="infra-map">
				<div class="map-container">
					<svg viewBox="0 0 500 350" class="schema-svg">
						<!-- Definitions -->
						<defs>
							<!-- Grid pattern -->
							<pattern id="grid" width="25" height="25" patternUnits="userSpaceOnUse">
								<path d="M 25 0 L 0 0 0 25" fill="none" stroke="var(--color-border)" stroke-width="0.3"/>
							</pattern>
							<!-- Glow filter -->
							<filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
								<feGaussianBlur stdDeviation="2" result="coloredBlur"/>
								<feMerge>
									<feMergeNode in="coloredBlur"/>
									<feMergeNode in="SourceGraphic"/>
								</feMerge>
							</filter>
							<!-- Arrow marker -->
							<marker id="arrow" markerWidth="6" markerHeight="6" refX="5" refY="3" orient="auto">
								<path d="M0,0 L6,3 L0,6" fill="none" stroke="var(--color-primary)" stroke-width="0.5"/>
							</marker>
						</defs>

						<!-- Background grid -->
						<rect width="100%" height="100%" fill="url(#grid)" opacity="0.5"/>

						<!-- Central Docker hub -->
						<g class="central-hub" filter="url(#glow)">
							<circle cx="250" cy="175" r="35" fill="none" stroke="var(--color-primary)" stroke-width="0.5"/>
							<circle cx="250" cy="175" r="28" fill="none" stroke="var(--color-primary)" stroke-width="0.8"/>
							<circle cx="250" cy="175" r="20" fill="rgba(0, 212, 170, 0.1)" stroke="var(--color-primary)" stroke-width="1"/>
							<!-- Docker icon placeholder -->
							<text x="250" y="172" text-anchor="middle" fill="var(--color-primary)" font-size="8" letter-spacing="0.1em">DOCKER</text>
							<text x="250" y="183" text-anchor="middle" fill="var(--color-primary-dim)" font-size="6">ENGINE</text>
						</g>

						<!-- Stack nodes with labels like satellite -->
						{#each stacksWithContainers as stack, i}
							{@const angle = (i * (360 / Math.max(stacksWithContainers.length, 1)) - 45) * (Math.PI / 180)}
							{@const radius = 120}
							{@const x = 250 + Math.cos(angle) * radius}
							{@const y = 175 + Math.sin(angle) * radius * 0.7}
							{@const isSelected = selectedStack?.name === stack.name}
							{@const healthPercent = stack.services > 0 ? Math.round((stack.runningServices / stack.services) * 100) : 0}
							{@const labelX = x + (x > 250 ? 50 : -50)}
							{@const labelY = y - 15}

							<!-- Connection line - thin like blueprint -->
							<g class="connection">
								<line
									x1="250" y1="175"
									x2={x} y2={y}
									stroke={stack.status === 'running' ? 'var(--color-primary)' : stack.status === 'partial' ? 'var(--color-warning)' : 'var(--color-text-dim)'}
									stroke-width={isSelected ? 0.8 : 0.4}
									stroke-dasharray={stack.status === 'stopped' ? '3,3' : 'none'}
									opacity={isSelected ? 1 : 0.6}
								/>
								<!-- Animated dot on line -->
								{#if stack.status === 'running'}
									<circle r="2" fill="var(--color-primary)" opacity="0.8">
										<animateMotion dur="3s" repeatCount="indefinite">
											<mpath href="#path-{i}"/>
										</animateMotion>
									</circle>
									<path id="path-{i}" d="M250,175 L{x},{y}" fill="none" stroke="none"/>
								{/if}
							</g>

							<!-- Label with percentage - like HG-ANTENNA 72% -->
							<g class="label-group" role="button" tabindex="0" style="cursor: pointer" onclick={() => selectedStack = stack} onkeydown={(e) => e.key === 'Enter' && (selectedStack = stack)}>
								<line
									x1={x} y1={y}
									x2={labelX} y2={labelY}
									stroke="var(--color-primary-dim)"
									stroke-width="0.4"
									opacity="0.6"
								/>
								<line
									x1={labelX} y1={labelY}
									x2={labelX + (x > 250 ? 40 : -40)} y2={labelY}
									stroke="var(--color-primary-dim)"
									stroke-width="0.4"
									opacity="0.6"
								/>

								<!-- Label box -->
								<rect
									x={x > 250 ? labelX : labelX - 55}
									y={labelY - 12}
									width="55"
									height="18"
									fill="var(--color-surface)"
									stroke={isSelected ? 'var(--color-primary)' : 'var(--color-border)'}
									stroke-width={isSelected ? 0.8 : 0.4}
								/>
								<text
									x={x > 250 ? labelX + 5 : labelX - 50}
									y={labelY}
									fill={stack.status === 'running' ? 'var(--color-primary)' : stack.status === 'partial' ? 'var(--color-warning)' : 'var(--color-text-dim)'}
									font-size="7"
									font-weight="500"
									letter-spacing="0.05em"
								>{stack.name.toUpperCase()}</text>

								<!-- Percentage badge -->
								<rect
									x={x > 250 ? labelX + 42 : labelX - 55}
									y={labelY + 5}
									width="28"
									height="12"
									fill={stack.status === 'running' ? 'var(--color-primary)' : stack.status === 'partial' ? 'var(--color-warning)' : 'var(--color-text-dim)'}
									opacity="0.9"
								/>
								<text
									x={x > 250 ? labelX + 56 : labelX - 41}
									y={labelY + 13}
									text-anchor="middle"
									fill="var(--color-void)"
									font-size="7"
									font-weight="600"
								>{healthPercent}%</text>
							</g>

							<!-- Stack node circle -->
							<g
								class="stack-node"
								class:selected={isSelected}
								style="cursor: pointer"
								role="button"
								tabindex="0"
								onclick={() => selectedStack = stack}
								onkeydown={(e) => e.key === 'Enter' && (selectedStack = stack)}
							>
								<circle
									cx={x} cy={y} r="18"
									fill="rgba(10, 10, 15, 0.9)"
									stroke={stack.status === 'running' ? 'var(--color-primary)' : stack.status === 'partial' ? 'var(--color-warning)' : 'var(--color-text-dim)'}
									stroke-width={isSelected ? 1 : 0.5}
								/>
								<circle
									cx={x} cy={y} r="12"
									fill="none"
									stroke={stack.status === 'running' ? 'var(--color-primary)' : stack.status === 'partial' ? 'var(--color-warning)' : 'var(--color-text-dim)'}
									stroke-width="0.3"
									opacity="0.5"
								/>
								<!-- Container count -->
								<text x={x} y={y + 3} text-anchor="middle" fill="var(--color-text)" font-size="9" font-weight="500">
									{stack.runningServices}
								</text>
							</g>
						{/each}

						<!-- No stacks message -->
						{#if stacksWithContainers.length === 0 && !dockerStore.loading}
							<text x="250" y="280" text-anchor="middle" fill="var(--color-text-dim)" font-size="10">
								No stacks found
							</text>
						{/if}

						<!-- Scanning indicator -->
						<g class="scanning">
							<text x="20" y="25" fill="var(--color-text-dim)" font-size="7" letter-spacing="0.1em">
								{systemStore.connected ? 'CONNECTED' : 'CONNECTING...'}
							</text>
							<rect x="20" y="30" width="60" height="4" fill="none" stroke="var(--color-border)" stroke-width="0.5"/>
							<rect x="20" y="30" width="25" height="4" fill="var(--color-primary)" opacity="0.8">
								<animate attributeName="width" values="0;60;0" dur="2s" repeatCount="indefinite"/>
							</rect>
						</g>
					</svg>
				</div>
			</div>
		</Panel>

		<!-- Selected Stack Details -->
		{#if selectedStack}
			<Panel title="STACK: {selectedStack.name.toUpperCase()}" class="stack-detail" variant="accent">
				<div class="stack-header">
					<StatusIndicator
						status={selectedStack.status === 'running' ? 'running' : selectedStack.status === 'partial' ? 'warning' : 'stopped'}
					/>
					<span class="stack-path">{selectedStack.path}</span>
				</div>

				<div class="containers-list">
					{#each selectedStack.containers as container}
						<div class="container-row" class:stopped={container.state !== 'running'}>
							<div class="container-info">
								<StatusIndicator status={container.state === 'running' ? 'running' : 'stopped'} size="sm" label="" />
								<span class="container-name">{container.name}</span>
								<span class="container-image">{container.image}</span>
							</div>
							<div class="container-stats">
								<span class="stat">{container.cpu.toFixed(1)}%</span>
								<span class="stat">{formatBytes(container.memory)}</span>
								<span class="stat-label">{container.status}</span>
							</div>
						</div>
					{:else}
						<div class="no-containers">No containers found for this stack</div>
					{/each}
				</div>

				<div class="stack-actions">
					<Button variant="secondary" size="sm">EDIT COMPOSE</Button>
					{#if selectedStack.status === 'stopped'}
						<Button variant="primary" size="sm" onclick={() => handleStartStack(selectedStack!.name)}>START</Button>
					{:else}
						<Button variant="ghost" size="sm" onclick={() => handleRestartStack(selectedStack!.name)}>RESTART</Button>
						<Button variant="danger" size="sm" onclick={() => handleStopStack(selectedStack!.name)}>STOP</Button>
					{/if}
				</div>
			</Panel>
		{/if}
	</div>

	<!-- Right Panel - Stats and Monitoring -->
	<div class="stats-panel">
		<!-- System Overview -->
		<Panel title="TRANSMISSION" expandable>
			<div class="transmission-stats">
				<div class="stat-row">
					<ProgressBar
						value={systemStore.stats.cpuUsage}
						label="CPU LOAD"
						variant="default"
						color="primary"
					/>
					<div class="stat-values">
						<span>{systemStore.stats.cpuUsage.toFixed(0)}%</span>
						<span class="dim">of {systemStore.stats.cpuCores} cores</span>
					</div>
				</div>
				<div class="stat-row">
					<ProgressBar
						value={systemStore.stats.memoryPercent}
						label="MEMORY QUOTA"
						variant="default"
						color="primary"
					/>
					<div class="stat-values">
						<span>{formatBytesRaw(systemStore.stats.memoryUsed)}</span>
						<span class="dim">/ {formatBytesRaw(systemStore.stats.memoryTotal)}</span>
					</div>
				</div>
				<div class="stat-row">
					<ProgressBar
						value={systemStore.stats.diskPercent}
						label="DISK ALLOCATION"
						variant="default"
						color={systemStore.stats.diskPercent > 80 ? 'warning' : 'primary'}
					/>
					<div class="stat-values">
						<span>{formatBytesRaw(systemStore.stats.diskUsed)}</span>
						<span class="dim">/ {formatBytesRaw(systemStore.stats.diskTotal)}</span>
					</div>
				</div>
			</div>
		</Panel>

		<!-- Resource Grid -->
		<div class="resource-grid">
			<Panel title="CONTAINERS" expandable>
				<div class="resource-content">
					<DataLabel label="Active" value={runningContainers} size="lg" highlight />
					<DataLabel label="Inactive" value={stoppedContainers} size="md" />
				</div>
			</Panel>
			<Panel title="SYSTEM POWER" expandable>
				<div class="resource-content">
					<DataLabel label="Uptime" value={formatUptime(systemStore.stats.uptime)} size="md" />
					<div class="power-indicator">
						<span class="power-value">{systemStore.connected ? 'ONLINE' : 'OFFLINE'}</span>
						<StatusIndicator status={systemStore.connected ? 'running' : 'stopped'} pulse size="lg" label="" />
					</div>
				</div>
			</Panel>
		</div>

		<!-- Network Chart -->
		<Panel title="NETWORK THROUGHPUT" expandable>
			<LineChart data={networkHistory} height={100} />
			<div class="network-stats">
				<DataLabel label="Stacks" value={dockerStore.stacks.length} size="md" />
				<DataLabel label="Volumes" value={dockerStore.volumes.length} size="md" />
				<DataLabel label="Networks" value={dockerStore.networks.length} size="md" />
			</div>
		</Panel>

		<!-- Logs Preview -->
		<Panel title="SYSTEM LOGS" expandable>
			<div class="logs-preview">
				<div class="log-entry">
					<span class="log-time">{new Date().toLocaleTimeString()}</span>
					<span class="log-source">system</span>
					<span class="log-message">Dashboard connected to backend</span>
				</div>
				{#if dockerStore.error}
					<div class="log-entry">
						<span class="log-time">{new Date().toLocaleTimeString()}</span>
						<span class="log-source warning">error</span>
						<span class="log-message">{dockerStore.error}</span>
					</div>
				{/if}
			</div>
		</Panel>
	</div>

	<!-- Bottom Panel - Status Bar with SOLAR style -->
	<div class="bottom-panel">
		<div class="status-bar">
			<!-- Left info -->
			<div class="status-info">
				<div class="info-block">
					<span class="info-label">{systemStore.stats.hostname || 'CELESTE'}</span>
					<StatusIndicator status={systemStore.connected ? 'running' : 'stopped'} size="sm" label={systemStore.connected ? 'ACTIVE' : 'DISCONNECTED'} />
				</div>
				<div class="info-grid">
					<div class="info-item">
						<span class="info-key">CONTAINERS</span>
						<span class="info-value">: {dockerStore.containers.length}</span>
					</div>
					<div class="info-item">
						<span class="info-key">IMAGES</span>
						<span class="info-value">: {dockerStore.images.length}</span>
					</div>
					<div class="info-item">
						<span class="info-key">VOLUMES</span>
						<span class="info-value">: {dockerStore.volumes.length}</span>
					</div>
					<div class="info-item">
						<span class="info-key">NETWORKS</span>
						<span class="info-value">: {dockerStore.networks.length}</span>
					</div>
				</div>
			</div>

			<!-- Center - Solar style bar -->
			<div class="solar-panel-section">
				<div class="solar-header">
					<span class="solar-label">CPU DISTRIBUTION</span>
					<span class="solar-value">{totalCpu.toFixed(0)}%</span>
				</div>
				<ProgressBar value={totalCpu} max={100} variant="solar" showValue={false} />
			</div>

			<!-- Right - Logs panel -->
			<div class="mini-logs">
				<div class="mini-logs-header">
					<span>STATUS</span>
					<button class="expand-btn" aria-label="Expand logs">
						<svg width="8" height="8" viewBox="0 0 8 8" fill="none">
							<path d="M1 7L7 1M7 1H2M7 1V6" stroke="currentColor" stroke-width="1"/>
						</svg>
					</button>
				</div>
				<div class="mini-log-entry">
					<span class="mini-log-time">{systemStore.stats.os}</span>
					<span class="mini-log-source">{systemStore.stats.platform}</span>
				</div>
				<div class="mini-log-text">{systemStore.connected ? 'WebSocket connected' : 'Connecting...'}</div>
			</div>
		</div>
	</div>
</div>

<style>
	.dashboard {
		display: grid;
		grid-template-columns: 1fr 380px;
		grid-template-rows: 1fr auto;
		gap: 1rem;
		height: calc(100vh - 140px);
		position: relative;
	}

	.loading-overlay {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(5, 5, 8, 0.8);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 1rem;
		z-index: 100;
		color: var(--color-primary);
		font-size: 0.8rem;
		letter-spacing: 0.1em;
	}

	.loading-spinner {
		width: 40px;
		height: 40px;
		border: 2px solid var(--color-border);
		border-top-color: var(--color-primary);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.visualization-panel {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		overflow: hidden;
	}

	.stats-panel {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
		overflow-y: auto;
	}

	.bottom-panel {
		grid-column: 1 / -1;
	}

	/* Infrastructure Map */
	.infra-map {
		height: 100%;
		min-height: 300px;
	}

	.map-container {
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.schema-svg {
		width: 100%;
		height: 100%;
		max-height: 350px;
	}

	.stack-node:hover circle {
		stroke-width: 1;
	}

	.stack-node.selected circle {
		filter: drop-shadow(0 0 6px var(--color-primary-dim));
	}

	/* Stack Detail */
	.stack-header {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding-bottom: 0.75rem;
		border-bottom: 1px solid var(--color-border);
		margin-bottom: 0.75rem;
	}

	.stack-path {
		font-size: 0.7rem;
		color: var(--color-text-dim);
		letter-spacing: 0.02em;
	}

	.containers-list {
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
	}

	.container-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.5rem 0.625rem;
		background: rgba(18, 18, 26, 0.6);
		border: 1px solid var(--color-border);
		border-left: 2px solid var(--color-primary-dim);
	}

	.container-row.stopped {
		opacity: 0.5;
		border-left-color: var(--color-text-dim);
	}

	.container-info {
		display: flex;
		align-items: center;
		gap: 0.625rem;
	}

	.container-name {
		font-size: 0.8rem;
		color: var(--color-text);
	}

	.container-image {
		font-size: 0.6rem;
		color: var(--color-text-dim);
		padding: 0.125rem 0.375rem;
		background: var(--color-border);
	}

	.container-stats {
		display: flex;
		align-items: center;
		gap: 0.875rem;
	}

	.stat {
		font-size: 0.7rem;
		font-variant-numeric: tabular-nums;
		color: var(--color-text);
		min-width: 45px;
		text-align: right;
	}

	.stat-label {
		font-size: 0.6rem;
		color: var(--color-text-dim);
	}

	.no-containers {
		font-size: 0.75rem;
		color: var(--color-text-dim);
		text-align: center;
		padding: 1rem;
	}

	.stack-actions {
		display: flex;
		gap: 0.5rem;
		margin-top: 0.75rem;
		padding-top: 0.75rem;
		border-top: 1px solid var(--color-border);
	}

	/* Transmission Stats */
	.transmission-stats {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.stat-row {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.stat-values {
		display: flex;
		justify-content: flex-end;
		gap: 0.25rem;
		font-size: 0.7rem;
		font-variant-numeric: tabular-nums;
	}

	.dim {
		color: var(--color-text-dim);
	}

	/* Resource Grid */
	.resource-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 0.75rem;
	}

	.resource-content {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.power-indicator {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.power-value {
		font-size: 0.75rem;
		color: var(--color-success);
		font-weight: 500;
		letter-spacing: 0.05em;
	}

	/* Network Stats */
	.network-stats {
		display: flex;
		gap: 1.5rem;
		margin-top: 0.75rem;
		padding-top: 0.75rem;
		border-top: 1px solid var(--color-border);
	}

	/* Logs */
	.logs-preview {
		display: flex;
		flex-direction: column;
		gap: 0.375rem;
		font-size: 0.7rem;
	}

	.log-entry {
		display: flex;
		gap: 0.625rem;
		padding: 0.375rem 0.5rem;
		background: rgba(18, 18, 26, 0.6);
		border-left: 1px solid var(--color-primary-dim);
	}

	.log-time {
		color: var(--color-text-dim);
		font-variant-numeric: tabular-nums;
		white-space: nowrap;
	}

	.log-source {
		color: var(--color-primary);
		min-width: 70px;
	}

	.log-source.warning {
		color: var(--color-warning);
	}

	.log-message {
		color: var(--color-text);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	/* Bottom Status Bar */
	.status-bar {
		display: grid;
		grid-template-columns: auto 1fr auto;
		gap: 2rem;
		padding: 0.75rem 1rem;
		background: rgba(10, 10, 15, 0.9);
		border: 1px solid var(--color-border);
		align-items: center;
	}

	.status-info {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.info-block {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.info-label {
		font-size: 0.75rem;
		font-weight: 500;
		letter-spacing: 0.1em;
		color: var(--color-text);
		text-transform: uppercase;
	}

	.info-grid {
		display: grid;
		grid-template-columns: repeat(4, auto);
		gap: 1.5rem;
	}

	.info-item {
		display: flex;
		font-size: 0.65rem;
	}

	.info-key {
		color: var(--color-text-dim);
		letter-spacing: 0.05em;
	}

	.info-value {
		color: var(--color-text);
		font-variant-numeric: tabular-nums;
	}

	/* Solar Panel Section */
	.solar-panel-section {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-width: 400px;
	}

	.solar-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.solar-label {
		font-size: 0.7rem;
		color: var(--color-text-dim);
		letter-spacing: 0.1em;
	}

	.solar-value {
		font-size: 0.8rem;
		color: var(--color-text);
		font-variant-numeric: tabular-nums;
	}

	/* Mini Logs */
	.mini-logs {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		min-width: 200px;
		padding: 0.5rem;
		background: rgba(18, 18, 26, 0.6);
		border: 1px solid var(--color-border);
	}

	.mini-logs-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		font-size: 0.65rem;
		color: var(--color-text-dim);
		letter-spacing: 0.1em;
		padding-bottom: 0.25rem;
		border-bottom: 1px solid var(--color-border);
	}

	.expand-btn {
		background: none;
		border: none;
		color: var(--color-text-dim);
		cursor: pointer;
		padding: 0.125rem;
	}

	.expand-btn:hover {
		color: var(--color-primary);
	}

	.mini-log-entry {
		display: flex;
		gap: 0.5rem;
		font-size: 0.6rem;
	}

	.mini-log-time {
		color: var(--color-text-dim);
	}

	.mini-log-source {
		color: var(--color-warning);
	}

	.mini-log-text {
		font-size: 0.6rem;
		color: var(--color-text);
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
</style>
