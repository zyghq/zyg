import { Outlet, createRootRouteWithContext } from "@tanstack/react-router";
import { QueryClient } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/toaster";
// @sanchitrk: enable these post demo.
// import { TanStackRouterDevtools } from "@tanstack/router-devtools";
// import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { SupabaseClient } from "@supabase/supabase-js";

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
  supabaseClient: SupabaseClient;
}>()({
  component: RootComponent,
});

function RootComponent() {
  return (
    <>
      <Outlet />
      <Toaster />
      {/* <TanStackRouterDevtools position="bottom-right" />
      <ReactQueryDevtools buttonPosition="bottom-left" /> */}
    </>
  );
}
