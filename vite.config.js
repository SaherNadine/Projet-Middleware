import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/config_api': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/config_api/, ''),
			},
			'/timetable_api': {
				target: 'http://localhost:8081',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/timetable_api/, ''),
			},
		},
	},
});
