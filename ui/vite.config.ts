import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: resolve(__dirname, '../static'),
    emptyOutDir: true,
    manifest: true,
    rollupOptions: {
      input: {
        // Register island entry points here as they are created.
        // Remove placeholder when adding the first real island.
        placeholder: resolve(__dirname, 'src/entries/placeholder.ts'),
      },
      output: {
        entryFileNames: 'js/[name].js',
        chunkFileNames: 'js/[name]-[hash].js',
        assetFileNames: (assetInfo) => {
          const info = assetInfo.name.split('.');
          const ext = info[info.length - 1];
          return `assets/[name]-[hash][extname]`;
        },
      },
    },
  },
  server: {
    port: 8081,
    strictPort: true,
    proxy: {
      '/api': 'http://localhost:8080',
      '^/(?!static/).*': 'http://localhost:8080',
    },
  },
});
