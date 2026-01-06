<script lang="ts">
	interface Props {
		title?: string;
		expandable?: boolean;
		variant?: 'default' | 'minimal' | 'accent';
		class?: string;
		children: import('svelte').Snippet;
	}

	let { title, expandable = false, variant = 'default', class: className = '', children }: Props = $props();
</script>

<div class="panel variant-{variant} {className}">
	<!-- Corner accents -->
	<div class="corner corner-tl"></div>
	<div class="corner corner-tr"></div>
	<div class="corner corner-bl"></div>
	<div class="corner corner-br"></div>

	{#if title}
		<div class="panel-header">
			<span class="panel-title">{title}</span>
			{#if expandable}
				<button class="expand-btn" aria-label="Expand">
					<svg width="10" height="10" viewBox="0 0 10 10" fill="none">
						<path d="M1 9L9 1M9 1H3M9 1V7" stroke="currentColor" stroke-width="1"/>
					</svg>
				</button>
			{/if}
		</div>
	{/if}
	<div class="panel-content">
		{@render children()}
	</div>
</div>

<style>
	.panel {
		background: rgba(10, 10, 15, 0.8);
		border: 1px solid var(--color-border);
		position: relative;
	}

	/* Corner decorations - Minority Report style */
	.corner {
		position: absolute;
		width: 8px;
		height: 8px;
		pointer-events: none;
	}

	.corner::before,
	.corner::after {
		content: '';
		position: absolute;
		background: var(--color-primary-dim);
	}

	.corner-tl { top: -1px; left: -1px; }
	.corner-tl::before { width: 8px; height: 1px; top: 0; left: 0; }
	.corner-tl::after { width: 1px; height: 8px; top: 0; left: 0; }

	.corner-tr { top: -1px; right: -1px; }
	.corner-tr::before { width: 8px; height: 1px; top: 0; right: 0; }
	.corner-tr::after { width: 1px; height: 8px; top: 0; right: 0; }

	.corner-bl { bottom: -1px; left: -1px; }
	.corner-bl::before { width: 8px; height: 1px; bottom: 0; left: 0; }
	.corner-bl::after { width: 1px; height: 8px; bottom: 0; left: 0; }

	.corner-br { bottom: -1px; right: -1px; }
	.corner-br::before { width: 8px; height: 1px; bottom: 0; right: 0; }
	.corner-br::after { width: 1px; height: 8px; bottom: 0; right: 0; }

	/* Accent variant has glowing corners */
	.variant-accent .corner::before,
	.variant-accent .corner::after {
		background: var(--color-primary);
		box-shadow: 0 0 4px var(--color-primary);
	}

	/* Top line accent */
	.panel::before {
		content: '';
		position: absolute;
		top: 0;
		left: 12px;
		right: 12px;
		height: 1px;
		background: linear-gradient(90deg, transparent, var(--color-primary-dim) 20%, var(--color-primary-dim) 80%, transparent);
		opacity: 0.4;
	}

	.variant-accent::before {
		opacity: 0.8;
		background: linear-gradient(90deg, transparent, var(--color-primary) 20%, var(--color-primary) 80%, transparent);
	}

	.panel-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.625rem 0.875rem;
		border-bottom: 1px solid var(--color-border);
	}

	.panel-title {
		font-size: 0.7rem;
		font-weight: 500;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		color: var(--color-text);
	}

	.expand-btn {
		background: none;
		border: none;
		color: var(--color-text-dim);
		cursor: pointer;
		padding: 0.25rem;
		transition: color 0.2s, transform 0.2s;
		line-height: 0;
	}

	.expand-btn:hover {
		color: var(--color-primary);
		transform: translate(1px, -1px);
	}

	.panel-content {
		padding: 0.875rem;
	}

	/* Minimal variant */
	.variant-minimal {
		background: transparent;
		border-color: transparent;
	}

	.variant-minimal .corner::before,
	.variant-minimal .corner::after {
		display: none;
	}

	.variant-minimal::before {
		display: none;
	}

	.variant-minimal .panel-header {
		border-bottom-color: var(--color-border);
	}
</style>
