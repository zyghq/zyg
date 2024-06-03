import { Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { QueryClient } from "@tanstack/react-query";

import { Toaster } from "@/components/ui/toaster";

import { TanStackRouterDevtools } from "@tanstack/router-devtools";

import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { AuthContext } from "@/auth";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
  auth: AuthContext;
}>()({
  component: RootComponent,
});

function RootComponent() {
  return (
    <>
      <Outlet />
      <Toaster />
      <TanStackRouterDevtools position="bottom-right" />
      <ReactQueryDevtools buttonPosition="bottom-left" />
    </>
  );
}
