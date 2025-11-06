<script lang="ts">
  import { onMount } from 'svelte';
  import '../app.css';

  const COOKIE = '__Host-gw_sid';
  const hasCookie = (n: string) => document.cookie.split('; ').some(c => c.startsWith(n + '='));

  onMount(async () => {
    if (self !== top && 'hasStorageAccess' in document) {
      try {
        // @ts-ignore
        const has = await document.hasStorageAccess();
        if (!has) { // @ts-ignore
          await document.requestStorageAccess(); location.reload(); return;
        }
      } catch {}
    }
    if (self !== top && !hasCookie(COOKIE)) {
      const ret = encodeURIComponent(location.href);
      top!.location.href = `https://dev.gradewise.org/lti/cookie-init?return=${ret}`;
    }
  });
</script>

<slot />