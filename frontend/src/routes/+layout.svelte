<script lang="ts">
	import '../app.css';
	import { api } from '$lib/api/client';

	let { children } = $props();

	let currentTime = $state(new Date());
	let version = $state('...');

	$effect(() => {
		const interval = setInterval(() => {
			currentTime = new Date();
		}, 1000);
		return () => clearInterval(interval);
	});

	$effect(() => {
		api.getHealth().then((health) => {
			version = health.version;
		}).catch(() => {
			version = '?';
		});
	});

	const formatTime = (date: Date) => {
		return date.toLocaleTimeString('en-US', {
			hour12: false,
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});
	};

	const formatDate = (date: Date) => {
		return date.toLocaleDateString('en-US', {
			year: 'numeric',
			month: '2-digit',
			day: '2-digit'
		}).replace(/\//g, '-');
	};
</script>

<svelte:head>
	<title>Aperture Science Network</title>
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous">
	<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@300;400;500;600&family=Orbitron:wght@400;500;600;700&display=swap" rel="stylesheet">
</svelte:head>

<div class="app">
	<header class="header">
		<div class="header-left">
			<div class="logo">
				<svg class="logo-icon" width="24" height="24" viewBox="0 0 24 24" fill="none">
					<circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
					<circle cx="12" cy="12" r="4" stroke="currentColor" stroke-width="1.5"/>
					<path d="M12 2V6M12 18V22M2 12H6M18 12H22" stroke="currentColor" stroke-width="1.5"/>
				</svg>
				<span class="logo-text">APERTURE SCIENCE NETWORK</span>
			</div>
			<span class="version">v{version}</span>
		</div>
		<div class="header-right">
			<div class="system-time">
				<span class="time">{formatTime(currentTime)}</span>
				<span class="date">{formatDate(currentTime)}</span>
			</div>
		</div>
	</header>

	<main class="main">
		{@render children()}
	</main>

	<footer class="footer">
		<div class="footer-left">
			<span class="footer-item">SYS: NOMINAL</span>
			<span class="footer-item">CONN: ACTIVE</span>
		</div>
		<div class="footer-center">
			<span class="footer-text">APERTURE SCIENCE - "WE DO WHAT WE MUST BECAUSE WE CAN"</span>
		</div>
		<div class="footer-right">
			<span class="footer-item">LAT: 47.6062</span>
			<span class="footer-item">LONG: -122.3321</span>
		</div>
	</footer>
</div>

<style>
	.app {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		background: var(--color-void);
	}

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
		border-bottom: 1px solid var(--color-border);
		background: var(--color-surface);
	}

	.header-left {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.logo {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.logo-icon {
		color: var(--color-primary);
		filter: drop-shadow(0 0 8px var(--color-primary-dim));
	}

	.logo-text {
		font-size: 0.875rem;
		font-weight: 600;
		letter-spacing: 0.2em;
		color: var(--color-text-bright);
	}

	.version {
		font-size: 0.625rem;
		padding: 0.25rem 0.5rem;
		background: var(--color-border);
		color: var(--color-text-dim);
		letter-spacing: 0.1em;
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 2rem;
	}

	.system-time {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 0.125rem;
	}

	.time {
		font-size: 1rem;
		font-variant-numeric: tabular-nums;
		color: var(--color-primary);
		letter-spacing: 0.1em;
	}

	.date {
		font-size: 0.625rem;
		color: var(--color-text-dim);
		letter-spacing: 0.05em;
	}

	.main {
		flex: 1;
		padding: 1.5rem;
		overflow: auto;
	}

	.footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 0.75rem 1.5rem;
		border-top: 1px solid var(--color-border);
		background: var(--color-surface);
		font-size: 0.625rem;
		letter-spacing: 0.1em;
		color: var(--color-text-dim);
		position: relative;
	}

	.footer-left,
	.footer-right {
		display: flex;
		gap: 1.5rem;
	}

	.footer-center {
		position: absolute;
		left: 50%;
		transform: translateX(-50%);
	}

	.footer-text {
		opacity: 0.5;
	}
</style>
