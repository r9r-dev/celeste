<script lang="ts">
	interface Props {
		value: number;
		max?: number;
		label?: string;
		showValue?: boolean;
		variant?: 'default' | 'segmented' | 'minimal' | 'solar';
		color?: 'primary' | 'success' | 'warning' | 'danger';
		class?: string;
	}

	let {
		value,
		max = 100,
		label,
		showValue = true,
		variant = 'default',
		color = 'primary',
		class: className = ''
	}: Props = $props();

	const percentage = $derived(Math.min(100, Math.max(0, (value / max) * 100)));
	const segments = 20;
	const filledSegments = $derived(Math.round((percentage / 100) * segments));

	// Solar variant - more segments with color gradient
	const solarSegments = 30;
	const filledSolarSegments = $derived(Math.round((percentage / 100) * solarSegments));
</script>

<div class="progress-container {className}">
	{#if label}
		<div class="progress-header">
			<span class="progress-label">{label}</span>
			{#if showValue}
				<span class="progress-value">{percentage.toFixed(0)}%</span>
			{/if}
		</div>
	{/if}

	{#if variant === 'segmented'}
		<div class="progress-segmented">
			{#each Array(segments) as _, i}
				<div
					class="segment"
					class:filled={i < filledSegments}
					class:color-primary={color === 'primary'}
					class:color-success={color === 'success'}
					class:color-warning={color === 'warning'}
					class:color-danger={color === 'danger'}
					style="--segment-index: {i}; --total-segments: {segments}"
				></div>
			{/each}
		</div>
	{:else if variant === 'solar'}
		<div class="progress-solar">
			{#each Array(solarSegments) as _, i}
				<div
					class="solar-segment"
					class:filled={i < filledSolarSegments}
					style="--segment-index: {i}; --total-segments: {solarSegments}; --fill-ratio: {filledSolarSegments / solarSegments}"
				></div>
			{/each}
		</div>
	{:else if variant === 'minimal'}
		<div class="progress-minimal">
			<div
				class="progress-fill color-{color}"
				style="width: {percentage}%"
			></div>
		</div>
	{:else}
		<div class="progress-bar">
			<div class="progress-track">
				<div
					class="progress-fill color-{color}"
					style="width: {percentage}%"
				></div>
				<div class="progress-glow color-{color}" style="width: {percentage}%"></div>
			</div>
			{#if showValue && !label}
				<span class="progress-inline-value">{value} / {max}</span>
			{/if}
		</div>
	{/if}
</div>

<style>
	.progress-container {
		width: 100%;
	}

	.progress-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 0.5rem;
	}

	.progress-label {
		font-size: 0.7rem;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		color: var(--color-text-dim);
	}

	.progress-value {
		font-size: 0.8rem;
		color: var(--color-text);
		font-variant-numeric: tabular-nums;
	}

	/* Segmented variant with color shift */
	.progress-segmented {
		display: flex;
		gap: 3px;
	}

	.segment {
		flex: 1;
		height: 12px;
		background: var(--color-border);
		transition: all 0.3s ease;
		opacity: 0.3;
	}

	.segment.filled {
		opacity: 1;
	}

	.segment.filled.color-primary {
		--hue-shift: calc(var(--segment-index) / var(--total-segments) * 30);
		background: hsl(calc(168 + var(--hue-shift)), 84%, calc(45% + var(--segment-index) / var(--total-segments) * 15%));
		box-shadow: 0 0 10px hsla(calc(168 + var(--hue-shift)), 84%, 50%, 0.6);
	}

	.segment.filled.color-success {
		background: var(--color-success);
		box-shadow: 0 0 8px color-mix(in srgb, var(--color-success) 50%, transparent);
	}

	.segment.filled.color-warning {
		background: var(--color-warning);
		box-shadow: 0 0 8px color-mix(in srgb, var(--color-warning) 50%, transparent);
	}

	.segment.filled.color-danger {
		background: var(--color-danger);
		box-shadow: 0 0 8px color-mix(in srgb, var(--color-danger) 50%, transparent);
	}

	/* Solar variant - like the SOLAR PANEL bar */
	.progress-solar {
		display: flex;
		gap: 2px;
		padding: 4px;
		background: var(--color-surface);
		border: 1px solid var(--color-border);
	}

	.solar-segment {
		flex: 1;
		height: 20px;
		background: var(--color-border);
		opacity: 0.2;
		transition: all 0.2s ease;
	}

	.solar-segment.filled {
		--progress: calc(var(--segment-index) / var(--total-segments));
		--hue: calc(168 + var(--progress) * 40);
		--lightness: calc(40% + var(--progress) * 20%);
		background: hsl(var(--hue), 80%, var(--lightness));
		box-shadow:
			0 0 8px hsla(var(--hue), 80%, 50%, 0.5),
			inset 0 1px 0 hsla(0, 0%, 100%, 0.2);
		opacity: 1;
	}

	/* Default bar variant with glow */
	.progress-bar {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.progress-track {
		flex: 1;
		height: 3px;
		background: var(--color-border);
		overflow: visible;
		position: relative;
	}

	.progress-fill {
		height: 100%;
		transition: width 0.3s ease;
		position: relative;
	}

	.progress-fill.color-primary {
		background: linear-gradient(90deg,
			hsl(168, 84%, 35%) 0%,
			hsl(168, 84%, 42%) 50%,
			hsl(178, 84%, 50%) 100%
		);
	}

	.progress-fill.color-success {
		background: linear-gradient(90deg, hsl(145, 80%, 35%), hsl(145, 100%, 50%));
	}

	.progress-fill.color-warning {
		background: linear-gradient(90deg, hsl(35, 80%, 40%), hsl(40, 100%, 50%));
	}

	.progress-fill.color-danger {
		background: linear-gradient(90deg, hsl(350, 80%, 40%), hsl(355, 100%, 60%));
	}

	.progress-glow {
		position: absolute;
		top: -2px;
		left: 0;
		height: 7px;
		filter: blur(4px);
		opacity: 0.6;
		pointer-events: none;
	}

	.progress-glow.color-primary {
		background: linear-gradient(90deg, transparent, var(--color-primary-glow));
	}

	.progress-glow.color-success {
		background: linear-gradient(90deg, transparent, var(--color-success));
	}

	.progress-glow.color-warning {
		background: linear-gradient(90deg, transparent, var(--color-warning));
	}

	.progress-glow.color-danger {
		background: linear-gradient(90deg, transparent, var(--color-danger));
	}

	.progress-inline-value {
		font-size: 0.7rem;
		color: var(--color-text-dim);
		white-space: nowrap;
		font-variant-numeric: tabular-nums;
	}

	/* Minimal variant */
	.progress-minimal {
		height: 2px;
		background: var(--color-border);
		overflow: hidden;
	}
</style>
