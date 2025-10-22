<script>
	export let data;
  </script>
  
  {#if !data.launched}
	<p>Open this app from Canvas (no LTI launch detected).</p>
  {:else}
	<slot />
  {/if}