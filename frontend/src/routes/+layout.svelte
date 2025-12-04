<script lang="ts">
  import { onMount } from 'svelte';
  import { base } from '$app/paths';
  import '../app.css';
  import type { LayoutData } from './$types';
  import {ltiSession} from "../lib/stores/ltiSession";

  const COOKIE = '__Host-gw_sid';
  const hasCookie = (n: string) => document.cookie.split('; ').some(c => c.startsWith(n + '='));


  export let data: LayoutData;

  // Initialize the ltiSession store from server-side data
  ltiSession.set({
      facultyId: data.lti?.facultyId ?? null,
      user: data.lti?.user ?? {},
  });
  onMount(async () => {
    if (self !== top && 'hasStorageAccess' in document) {
      try {
        // @ts-ignore
        const has = await document.hasStorageAccess();
        if (!has) {
          // @ts-ignore
          await document.requestStorageAccess();
          location.reload();
          return;
        }
      } catch {}
    }
    if (self !== top && !hasCookie(COOKIE)) {
      const ret = encodeURIComponent(location.href);
      top!.location.href = `https://dev.gradewise.org/lti/cookie-init?return=${ret}`;
      return;
    }
  });
</script>

<div class="min-h-screen bg-slate-50 text-slate-900">
  <header class="border-b bg-white">
    <div class="mx-auto flex max-w-4xl items-center justify-between px-4 py-3">
      <a href="{base}/" class="font-semibold">Gradewise</a>
      <nav class="flex gap-4 text-sm">
        <a href="{base}/instructor" class="hover:underline">Instructor</a>
        <a href="{base}/submit/demo-assignment" class="hover:underline">Student demo</a>
        <a href="{base}/instructor/courses/new" class="hover:underline">New Course</a>
          <a href="{base}/instructor/courses/list" class="hover:underline">Show Courses</a>
      </nav>
    </div>
  </header>

  <main class="mx-auto max-w-4xl px-4 py-6">
    <slot />
  </main>
</div>