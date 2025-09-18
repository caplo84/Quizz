import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  define: {
    // Make process.env available in the browser
    'process.env': process.env
  },
  server: {
    port: 5173,
    host: '0.0.0.0', // This allows external connections
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL || 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
        configure: (proxy, options) => {
          console.log('🔧 Proxy configured for /api -> ', options.target);
        }
      }
    }
  }
})
