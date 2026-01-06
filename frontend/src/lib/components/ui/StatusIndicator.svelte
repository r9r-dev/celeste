<script lang="ts">
	interface Props {
		status: 'running' | 'stopped' | 'warning' | 'error' | 'pending';
		label?: string;
		pulse?: boolean;
		size?: 'sm' | 'md' | 'lg';
		class?: string;
	}

	let {
		status,
		label,
		pulse = true,
		size = 'md',
		class: className = ''
	}: Props = $props();

	const statusLabels: Record<string, string> = {
		running: 'RUNNING',
		stopped: 'STOPPED',
		warning: 'WARNING',
		error: 'ERROR',
		pending: 'PENDING'
	};
</script>

<div class="status-indicator size-{size} {className}">
	<span class="dot status-{status}" class:pulse></span>
	{#if label !== undefined}
		<span class="label">{label || statusLabels[status]}</span>
	{/if}
</div>

<style>
	.status-indicator {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
	}

	.dot {
		border-radius: 50%;
		flex-shrink: 0;
	}

	.size-sm .dot { width: 6px; height: 6px; }
	.size-md .dot { width: 8px; height: 8px; }
	.size-lg .dot { width: 10px; height: 10px; }

	.status-running {
		background: var(--color-success);
		box-shadow: 0 0 8px var(--color-success);
	}

	.status-stopped {
		background: var(--color-text-dim);
	}

	.status-warning {
		background: var(--color-warning);
		box-shadow: 0 0 8px var(--color-warning);
	}

	.status-error {
		background: var(--color-danger);
		box-shadow: 0 0 8px var(--color-danger);
	}

	.status-pending {
		background: var(--color-info);
		box-shadow: 0 0 8px var(--color-info);
	}

	.pulse {
		animation: pulse 2s ease-in-out infinite;
	}

	@keyframes pulse {
		0%, 100% { opacity: 1; }
		50% { opacity: 0.5; }
	}

	.label {
		font-size: 0.75rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: var(--color-text);
	}

	.size-sm .label { font-size: 0.625rem; }
	.size-lg .label { font-size: 0.875rem; }
</style>
