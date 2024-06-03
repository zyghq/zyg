import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";
import { getOrCreateZygAccount } from "@/db/api";
import { AccountStoreProvider } from "@/providers";

const accountQueryOptions = (token: string) =>
  queryOptions({
    queryKey: ["account", token],
    queryFn: async () => {
      if (!token || token === "") return null;
      const { error, data } = await getOrCreateZygAccount(token);
      if (error) throw new Error("failed to authenticate user account.");
      return data;
    },
  });

export const Route = createFileRoute("/_auth")({
  beforeLoad: async ({ context }) => {
    const { auth } = context;
    const session = await auth?.client.auth.getSession();
    const { error, data } = session;
    if (error || !data?.session) {
      throw redirect({ to: "/signin" });
    }
    const token = data.session.access_token;
    return { token };
  },
  loader: async ({ context: { queryClient, token } }) => {
    return queryClient.ensureQueryData(accountQueryOptions(token));
  },
  component: AuthLayout,
});

function AuthLayout() {
  const { token } = Route.useRouteContext();
  const { data } = useSuspenseQuery(accountQueryOptions(token));

  return (
    <AccountStoreProvider
      initialValue={{
        hasData: data ? true : false,
        error: data ? null : new Error("failed to fetch account details"),
        token,
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
