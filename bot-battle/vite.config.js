import { defineConfig } from "vite";

// Sets up a proxy for the API requests to avoid CORS issues during development. 
// All requests to /api will be proxied to http://localhost:8080.
export default defineConfig({
  server: {
    proxy: {
      "/api": "http://localhost:8080",
      "/ws": {
        target: "ws://localhost:8080",
        ws: true
      },
    }
  }
});