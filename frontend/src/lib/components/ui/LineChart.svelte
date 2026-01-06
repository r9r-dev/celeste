<script lang="ts">
	import { onMount } from 'svelte';
	import Chart from 'chart.js/auto';

	interface Props {
		data: number[];
		labels?: string[];
		color?: string;
		fill?: boolean;
		height?: number;
		class?: string;
	}

	let {
		data,
		labels,
		color = '#00d4aa',
		fill = true,
		height = 100,
		class: className = ''
	}: Props = $props();

	let canvas: HTMLCanvasElement;
	let chart: Chart | null = null;

	const chartLabels = $derived(labels || data.map((_, i) => ''));

	onMount(() => {
		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		chart = new Chart(ctx, {
			type: 'line',
			data: {
				labels: chartLabels,
				datasets: [{
					data: data,
					borderColor: color,
					borderWidth: 2,
					backgroundColor: fill
						? `${color}15`
						: 'transparent',
					fill: fill,
					tension: 0.4,
					pointRadius: 0,
					pointHoverRadius: 4,
					pointHoverBackgroundColor: color,
					pointHoverBorderColor: '#fff',
					pointHoverBorderWidth: 2
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				interaction: {
					intersect: false,
					mode: 'index'
				},
				plugins: {
					legend: { display: false },
					tooltip: {
						backgroundColor: '#12121a',
						borderColor: '#2a2a3a',
						borderWidth: 1,
						titleColor: '#808090',
						bodyColor: '#e0e0e0',
						padding: 12,
						displayColors: false,
						callbacks: {
							title: () => ''
						}
					}
				},
				scales: {
					x: {
						display: false,
						grid: { display: false }
					},
					y: {
						display: false,
						grid: { display: false },
						beginAtZero: true
					}
				}
			}
		});

		return () => {
			chart?.destroy();
		};
	});

	$effect(() => {
		if (chart) {
			chart.data.labels = chartLabels;
			chart.data.datasets[0].data = data;
			chart.update('none');
		}
	});
</script>

<div class="chart-container {className}" style="height: {height}px">
	<canvas bind:this={canvas}></canvas>
</div>

<style>
	.chart-container {
		position: relative;
		width: 100%;
	}
</style>
