import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { queryOptions } from "@tanstack/react-query";
import { getOrCreateZygAccount } from "@/db/api";
import { AccountStoreProvider } from "@/providers";

const accountQueryOptions = (token: string) =>
  queryOptions({
    queryKey: ["account", token],
    queryFn: async () => {
      const { error, data } = await getOrCreateZygAccount(token);
      if (error) throw new Error("failed to authenticate user account.");
      return data;
    },
  });

export const Route = createFileRoute("/_auth")({
  beforeLoad: async ({ context }) => {
    const { supabaseClient } = context;
    const { error, data } = await supabaseClient.auth.getSession();
    if (error || !data?.session) {
      throw redirect({ to: "/signin" });
    }

    const token = data.session.access_token as string;
    return { token };
  },
  loader: async ({ context: { queryClient, supabaseClient } }) => {
    const { error, data } = await supabaseClient.auth.getSession();
    if (error || !data?.session) throw redirect({ to: "/signin" });
    const token = data.session.access_token;
    return queryClient.ensureQueryData(accountQueryOptions(token));
  },
  component: AuthLayout,
});

function AuthLayout() {
  const data = Route.useLoaderData();
  return (
    <AccountStoreProvider
      initialValue={{
        hasData: data ? true : false,
        error: data ? null : new Error("failed to fetch account details"),
        account: data
          ? {
              email: data.email,
              accountId: data.accountId,
              name: data.name,
              provider: data.provider,
              createdAt: data.createdAt,
              updatedAt: data.updatedAt,
            }
          : null,
      }}
    >
      <Outlet />
    </AccountStoreProvider>
  );
}
