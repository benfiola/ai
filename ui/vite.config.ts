import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react-swc";

import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig(() => {
  const backend = process.env.BACKEND || "localhost:8080";
  return {
    plugins: [
      react(),
      tailwindcss(),
      tanstackRouter({ target: "react", autoCodeSplitting: true }),
    ],
    server: {
      proxy: {
        "/api": {
          target: `http://${backend}`,
          changeOrigin: true,
        },
      },
    },
  };
});
