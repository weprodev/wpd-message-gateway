import path from "node:path"
import { defineConfig } from "vitest/config"
import react from "@vitejs/plugin-react"

const apiTarget = process.env.API_TARGET || "http://127.0.0.1:10101"

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    host: "0.0.0.0",
    port: 10104,
    proxy: {
      "/api": { target: apiTarget, changeOrigin: true },
      "/v1": { target: apiTarget, changeOrigin: true },
    },
  },
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: "./src/test-setup.ts",
    pool: "forks",
  },
})
