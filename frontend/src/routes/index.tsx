import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  beforeLoad: async ({ context }) => {
    const { supabaseClient } = context;
    const { error, data } = await supabaseClient.auth.getSession();
    const isAuthenticated = !error && data?.session;
    if (!isAuthenticated) {
      throw redirect({ to: "/signin" });
    }
  },
  component: () => <div>Index Root at /</div>,
});
