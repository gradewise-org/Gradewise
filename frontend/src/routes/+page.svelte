<script lang="ts">
	import { PUBLIC_BASE_URL } from '$env/static/public';
	import type { PageProps } from './$types';
	import { onMount } from 'svelte';

	let { data }: PageProps = $props();

	let count = $state(0);
	let healthStatus = $state('Loading...');

	onMount(() => {
		const updateHealth = async () => {
			try {
				const response = await fetch(`${PUBLIC_BASE_URL}/api/health`);
				healthStatus = (await response.text()) + ' - ' + new Date().toLocaleTimeString();
			} catch (error) {
				healthStatus = 'Error fetching health status';
			}
		};

		updateHealth();
		const interval = setInterval(updateHealth, 1000);

		return () => clearInterval(interval);
	});
</script>

<div class="m-1">
	<h1>Welcome to SvelteKit</h1>
	<p>
		Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation
	</p>

	<hr class="my-2" />

	<div class="flex flex-col gap-4">
		<button
			onclick={() => {
				count += 1;
			}}
			class="rounded-md bg-blue-500 p-2 text-white"
		>
			I've been clicked {count} times
		</button>

		<div>
			<h2>Data from the server (server-side fetch)</h2>
			<pre>{data.message}</pre>
		</div>

		<div>
			<h2>Health Status (client-side fetch)</h2>
			<pre>{healthStatus}</pre>
		</div>
	</div>
</div>
