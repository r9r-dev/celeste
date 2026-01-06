<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		class?: string;
		onclick?: () => void;
		children: Snippet;
	}

	let {
		variant = 'primary',
		size = 'md',
		disabled = false,
		class: className = '',
		onclick,
		children
	}: Props = $props();
</script>

<button
	class="btn variant-{variant} size-{size} {className}"
	{disabled}
	{onclick}
>
	{@render children()}
</button>

<style>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 0.5rem;
		font-family: var(--font-mono);
		font-size: 0.75rem;
		font-weight: 500;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		border: 1px solid transparent;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	/* Sizes */
	.size-sm { padding: 0.375rem 0.75rem; }
	.size-md { padding: 0.5rem 1rem; }
	.size-lg { padding: 0.75rem 1.5rem; font-size: 0.875rem; }

	/* Variants */
	.variant-primary {
		background: var(--color-primary);
		color: var(--color-void);
		border-color: var(--color-primary);
	}

	.variant-primary:hover:not(:disabled) {
		background: var(--color-primary-glow);
		box-shadow: 0 0 20px var(--color-primary-dim);
	}

	.variant-secondary {
		background: transparent;
		color: var(--color-primary);
		border-color: var(--color-primary);
	}

	.variant-secondary:hover:not(:disabled) {
		background: var(--color-primary);
		color: var(--color-void);
	}

	.variant-ghost {
		background: transparent;
		color: var(--color-text-dim);
		border-color: var(--color-border);
	}

	.variant-ghost:hover:not(:disabled) {
		color: var(--color-text);
		border-color: var(--color-border-bright);
	}

	.variant-danger {
		background: transparent;
		color: var(--color-danger);
		border-color: var(--color-danger);
	}

	.variant-danger:hover:not(:disabled) {
		background: var(--color-danger);
		color: var(--color-void);
	}
</style>
