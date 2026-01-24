import path from "path"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
    server: {
        port: 10104,
        host: "0.0.0.0", // Ensure server binds to all interfaces for Docker
        proxy: {
            "/api": {
                target: process.env.API_TARGET || "http://localhost:10101",
                changeOrigin: true,
            },
        },
    },
    build: {
        outDir: "dist",
        emptyOutDir: true,
    },
})
