import { createServerClient } from "@supabase/ssr";

export function createClient(cookieStore) {
  return createServerClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY,
    {
      cookies: {
        get(name) {
          return cookieStore.get(name)?.value;
        },
        set(name, value, options) {
          return cookieStore.set({ name, value, ...options });
        },
        remove(name, options) {
          return cookieStore.set({ name, value: "", ...options });
        },
      },
    }
  );
}
