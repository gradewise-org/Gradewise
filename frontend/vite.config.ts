import tailwindcss from '@tailwindcss/vite';
import { svelteTesting } from '@testing-library/svelte/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

//TODO: Read data from env variables during build (add variable in Tiltfile)
export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		host: true, port: 8080, strictPort: true,
		allowedHosts: ['dev.gradewise.org','jon-dev.gradewise.org'],
		origin: 'https://dev.gradewise.org',
		hmr: { protocol: 'wss', host: 'localhost', port: 8080 }
	  },
	preview: { port: 8080},

	test: {
		workspace: [
			{
				extends: './vite.config.ts',
				plugins: [svelteTesting()],
				test: {
					name: 'client',
					environment: 'jsdom',
					clearMocks: true,
					include: ['src/**/*.svelte.{test,spec}.{js,ts}'],
					exclude: ['src/lib/server/**'],
					setupFiles: ['./vitest-setup-client.ts']
				}
			},
			{
				extends: './vite.config.ts',
				test: {
					name: 'server',
					environment: 'node',
					include: ['src/**/*.{test,spec}.{js,ts}'],
					exclude: ['src/**/*.svelte.{test,spec}.{js,ts}']
				}
			}
		]
	}
});
