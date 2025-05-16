import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode, ssrBuild }) => {
  const ret = {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    server: {
      watch: {
        usePolling: true
      }
    }
  };

  ret.define = {
    // Do not modify this constant, it is used in the evaluation.
    "__API_URL__": JSON.stringify("http://localhost:3000"),
    // "__API_URL__": JSON.stringify("http://192.168.1.59:3000"),
  };

  return ret;
})