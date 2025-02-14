/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_WORKOS_CLIENT_ID: string;
  readonly VITE_SENTRY_DSN: string;
  readonly VITE_SENTRY_ENABLED: string;
  readonly VITE_SENTRY_ENV: string;
}
