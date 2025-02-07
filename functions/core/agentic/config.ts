export const DEFAULT_LLM_PROVIDER = Deno.env.get("DEFAULT_LLM_PROVIDER") ??
  "gemini-2.0-flash-lite-preview-02-05";

export const LLM_RETRY_CONFIG = {
  initialRetryIntervalMillis: 1000, // Start with 1s delay
  retryIntervalFactor: 1.5, // Gentler exponential backoff with some jitter
  maxRetryIntervalMillis: 5000, // Cap at 5s between retries
  maxRetryAttempts: 5, // Keep the 5 retry attempts
  maxRetryDurationMillis: 30000, // Allow up to 30s total for retries
} as const;

export const ZYG_SRV_BASE_URL = Deno.env.get("ZYG_SRV_BASE_URL") ??
  "http://localhost:8080";
