import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  beforeLoad: async ({ context }) => {
    console.log("**** beforeLoad in index ****");
    const { supabaseClient } = context;
    const { error, data } = await supabaseClient.auth.getSession();
    const isAuthenticated = !error && data?.session;
    if (!isAuthenticated) {
      throw redirect({ to: "/signin" });
    }
    console.log("**** beforeLoad in index end ****");
  },
  component: () => <div>Index Root at /</div>,
});
