import {
  Outlet,
  createRootRouteWithContext,
  redirect as redirectRoute,
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

import { Session, User, SupabaseClient } from "@supabase/supabase-js";

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

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
  auth: ReturnType<typeof createAuthContext>;
  session: Session | null;
  account: User | null;
  supaClient: SupabaseClient;
  AccountStore: typeof AccountStore;
}>()({
  beforeLoad: async ({ context, location }) => {
    const unprotectedRoutes = /^\/(signup|signin|recover|favicon\.ico)(\?.*)?$/;
    if (unprotectedRoutes.test(location.pathname)) return;

    // protected routes
    const { session } = context;
    if (!session) {
      throw redirectRoute({
        to: "/signin",
        search: {
          redirect: location.href,
        },
      });
    }
  },
  loader: async ({ context: { queryClient, session } }) => {
    const token = session?.access_token || "";
    return queryClient.ensureQueryData(accountQueryOptions(token));
  },
  component: RootComponent,
});

function RootComponent() {
  const { session } = Route.useRouteContext();
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
