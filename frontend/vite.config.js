import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export default defineConfig({
	// add allowd hosts
	server: {
		allowedHosts: ['especially-strong-piglet.ngrok-free.app'],
	},

	plugins: [vue()],
});
