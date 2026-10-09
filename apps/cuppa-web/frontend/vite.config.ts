import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  // Relative URLs, so the page works from any sub-path (GitHub Pages serves /cuppa/).
  base: './',
  plugins: [react()]
})
