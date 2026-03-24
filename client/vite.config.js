import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    // In development, proxy /api requests to your Go server
    // So you can run `npm run dev` and `go run main.go` side by side
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  },
  build: {
    // Output goes to ../web/dist — Go server will serve this folder
    outDir: '../web/dist',
    emptyOutDir: true,
  }
})
