import {
  Outlet,
  createRootRouteWithContext,
  redirect,
} from "@tanstack/react-router";
import {
  QueryClient,
  queryOptions,
  useSuspenseQuery,
} from "@tanstack/react-query";

import { Toaster } from "@/components/ui/toaster";

import { TanStackRouterDevtools } from "@tanstack/router-devtools";

import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createAuthContext } from "@/auth";
import { AccountStore } from "@/providers";
import { getOrCreateZygAccount } from "@/db/api";

const accountQueryOptions = (token: string) =>
  queryOptions({
    queryKey: ["account", token],
    queryFn: async () => {
      const { error, data } = await getOrCreateZygAccount(token);
      if (error) throw new Error("failed to fetch account details.");
      return data;
    },
  });

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
  auth: ReturnType<typeof createAuthContext>;
  AccountStore: typeof AccountStore;
}>()({
  beforeLoad: async ({ context, location }) => {
    const unprotectedRoutes = /^\/(login|signin|favicon\.ico)(\?.*)?$/;
    if (unprotectedRoutes.test(location.pathname)) return;
    const { auth } = context;
    const { error: errSession, data } = await auth.client.auth.getSession();
    if (errSession || !data?.session) {
      throw redirect({
        to: "/login",
        search: {
          redirect: location.href,
        },
      });
    }
  },
  loader: async ({ context: { queryClient, auth } }) => {
    const { client } = auth;
    const { data } = await client.auth.getSession();
    const { session } = data;
    const token = session?.access_token || "";
    return queryClient.ensureQueryData(accountQueryOptions(token));
  },
  component: RootComponent,
});

function RootComponent() {
  const { auth } = Route.useRouteContext();
  const useAuth = auth.useContext();
  const session = useAuth?.session;
  const token = session?.access_token || "";
  const { data } = useSuspenseQuery(accountQueryOptions(token));

  return (
    <>
      <AccountStore.Provider
        initialValue={{
          hasData: data ? true : false,
          error: data ? null : new Error("failed to fetch account details."),
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
      </AccountStore.Provider>
      <Toaster />
      <TanStackRouterDevtools position="bottom-right" />
      <ReactQueryDevtools buttonPosition="bottom-left" />
    </>
  );
}
