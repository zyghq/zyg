import ReactDOM from "react-dom/client";
import { Suspense } from "react";

import { RouterProvider, createRouter } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider, AccountStore } from "@/providers";

import { createClient } from "@supabase/supabase-js";
import { createAuthContext } from "@/auth";
import Loading from "@/components/loading";

import "./globals.css";

const supabase = createClient(
  import.meta.env.VITE_SUPABASE_URL,
  import.meta.env.VITE_SUPABASE_ANON_KEY
);

// eslint-disable-next-line react-refresh/only-export-components
const Auth = createAuthContext(supabase);

// Import the generated route tree
import { routeTree } from "./routeTree.gen";

const queryClient = new QueryClient();

// Set up a Router instance
const router = createRouter({
  routeTree,
  context: {
    queryClient,
    session: undefined!,
    account: undefined!,
    supaClient: undefined!,
    auth: undefined!,
    AccountStore: undefined!,
  },
  defaultPreload: "intent",
  // Since we're using React Query, we don't want loader calls to ever be stale
  // This will ensure that the loader is always called when the route is preloaded or visited
  defaultPreloadStaleTime: 0,
});

// Register the router instance for type safety
declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

// eslint-disable-next-line react-refresh/only-export-components
function AuthRouter() {
  const { isLoading, session, user, client } = Auth.useContext() || {};

  if (isLoading) {
    return <Loading />;
  }

  return (
    <Suspense fallback={<Loading />}>
      <RouterProvider
        router={router}
        context={{
          auth: Auth,
          session,
          account: user,
          supaClient: client,
          AccountStore,
        }}
      />
    </Suspense>
  );
}

// Render the app
const rootElement = document.getElementById("app")!;
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <Auth.Provider>
          <AuthRouter />
        </Auth.Provider>
      </QueryClientProvider>
    </ThemeProvider>
  );
}
