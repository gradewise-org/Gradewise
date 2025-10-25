import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

export default {
  preprocess: vitePreprocess(),
  compilerOptions: { runes: false },
  kit: {
    adapter: adapter(),
    paths: { base: '/app' }   // always serve under /app
  }
};