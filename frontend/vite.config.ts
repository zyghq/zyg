import { sentryVitePlugin } from "@sentry/vite-plugin";
import { TanStackRouterVite } from "@tanstack/router-vite-plugin";
import react from "@vitejs/plugin-react-swc";
import path from "path";
import { defineConfig, loadEnv } from "vite";

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  return {
    build: {
      sourcemap: true,
    },

    plugins: [
      react(),
      TanStackRouterVite(),
      sentryVitePlugin({
        org: "zyghq",
        project: "zyg-frontend",
        telemetry: env.SENTRY_TELEMETRY_ENABLED === "1" || false,
      }),
    ],

    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },

    server: {
      port: 3000,
    },
  };
});
