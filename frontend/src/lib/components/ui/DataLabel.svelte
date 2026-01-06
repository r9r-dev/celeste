<script lang="ts">
	interface Props {
		label: string;
		value: string | number;
		unit?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		highlight?: boolean;
		class?: string;
	}

	let {
		label,
		value,
		unit,
		size = 'md',
		highlight = false,
		class: className = ''
	}: Props = $props();
</script>

<div class="data-label size-{size} {className}" class:highlight>
	<span class="label">{label}</span>
	<div class="value-container">
		<span class="value">{value}</span>
		{#if unit}
			<span class="unit">{unit}</span>
		{/if}
	</div>
</div>

<style>
	.data-label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}

	.label {
		font-size: 0.625rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		color: var(--color-text-dim);
	}

	.value-container {
		display: flex;
		align-items: baseline;
		gap: 0.25rem;
	}

	.value {
		font-variant-numeric: tabular-nums;
		color: var(--color-text);
	}

	.unit {
		font-size: 0.75em;
		color: var(--color-text-dim);
	}

	/* Sizes */
	.size-sm .value { font-size: 0.875rem; }
	.size-md .value { font-size: 1.25rem; }
	.size-lg .value { font-size: 1.75rem; }
	.size-xl .value { font-size: 2.5rem; font-weight: 300; }

	/* Highlight */
	.highlight .value {
		color: var(--color-primary);
		text-shadow: 0 0 20px var(--color-primary-dim);
	}
</style>
