<script lang="ts">
	import { PUBLIC_BASE_URL } from '$env/static/public';
	import type { PageData } from './$types';
	import { onMount } from 'svelte';
  
	export let data: PageData;   
	let count = 0;               
	let healthStatus = 'Loading...'; 
  
	onMount(() => {
	  const updateHealth = async () => {
		try {
		  const r = await fetch(`${PUBLIC_BASE_URL}/api/health`);
		  healthStatus = (await r.text()) + ' - ' + new Date().toLocaleTimeString();
		} catch {
		  healthStatus = 'Error fetching health status';
		}
	  };
	  updateHealth();
	  const t = setInterval(updateHealth, 1000);
	  return () => clearInterval(t);
	});
  </script>
  
  <div class="m-1">
	<h1>Welcome to SvelteKit</h1>
	<button class="rounded-md bg-blue-500 p-2 text-white" on:click={() => (count += 1)}>
	  I've been clicked {count} times
	</button>
  
	<h2>Data from the server</h2>
	<pre>{data.message}</pre>
  
	<h2>Health Status</h2>
	<pre>{healthStatus}</pre>
  </div>