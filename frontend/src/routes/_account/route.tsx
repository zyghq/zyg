import { getOrCreateZygAccount } from "@/db/api";
import { AccountStoreProvider } from "@/providers";
import { queryOptions, useQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

const accountQueryOptions = (token: string) =>
  queryOptions({
    queryFn: async () => {
      const { data, error } = await getOrCreateZygAccount(token);
      if (error) throw new Error("failed to authenticate user account.");
      return data;
    },
    queryKey: ["account", token],
  });

export const Route = createFileRoute("/_account")({
  beforeLoad: async ({ context }) => {
    const { supabaseClient } = context;
    const { data, error } = await supabaseClient.auth.getSession();
    if (error || !data?.session) {
      throw redirect({ to: "/signin" });
    }

    const token = data.session.access_token as string;
    return { token };
  },
  component: AuthLayout,
  loader: async ({ context: { queryClient, supabaseClient } }) => {
    const { data, error } = await supabaseClient.auth.getSession();
    if (error || !data?.session) throw redirect({ to: "/signin" });
    const token = data.session.access_token;
    return queryClient.ensureQueryData(accountQueryOptions(token));
  },
});

function AuthLayout() {
  const { token } = Route.useRouteContext();
  const initialData = Route.useLoaderData();

  const response = useQuery({
    enabled: !!token,
    initialData: initialData,
    queryFn: async () => {
      const { data, error } = await getOrCreateZygAccount(token);
      if (error) throw new Error("failed to authenticate user account.");
      return data;
    },
    queryKey: ["account", token],
    staleTime: 1000 * 60,
  });

  const { data } = response;

  return (
    <AccountStoreProvider
      initialValue={{
        account: data
          ? {
              accountId: data.accountId,
              createdAt: data.createdAt,
              email: data.email,
              name: data.name,
              provider: data.provider,
              updatedAt: data.updatedAt,
            }
          : null,
        error: data ? null : new Error("failed to fetch account details"),
        hasData: !!data,
      }}
    >
      <Outlet />
    </AccountStoreProvider>
  );
}
