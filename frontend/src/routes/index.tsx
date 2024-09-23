import { createFileRoute, redirect } from "@tanstack/react-router";

// TODO: redirect to last used workspace.
export const Route = createFileRoute("/")({
  beforeLoad: async ({ context }) => {
    const { supabaseClient } = context;
    const { data, error } = await supabaseClient.auth.getSession();
    const isAuthenticated = !error && data?.session;
    if (!isAuthenticated) {
      throw redirect({ to: "/signin" });
    }
  },
  component: () => <div>Index Root at /</div>,
});
