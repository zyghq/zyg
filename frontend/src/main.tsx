// import React from "react";
import ReactDOM from "react-dom/client";

import { RouterProvider, createRouter } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ThemeProvider } from "@/providers";
// import { AuthProvider, useAuth } from "@/auth";
import { createClient } from "@supabase/supabase-js";

import "./globals.css";

const supabaseClient = createClient(
  import.meta.env.VITE_SUPABASE_URL,
  import.meta.env.VITE_SUPABASE_ANON_KEY
);

// Import the generated route tree
import { routeTree } from "./routeTree.gen";

const queryClient = new QueryClient();

// Set up a Router instance
const router = createRouter({
  routeTree,
  context: {
    queryClient,
    supabaseClient,
    // isAuthenticated?: false,
    // session: undefined!,
    // isLoading: undefined!,
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
// function AuthRouter() {
//   const auth = useAuth();
//   const { isLoading, isAuthenticated, session } = auth;
//   console.log("*** in auth router start ****");
//   console.log("isLoading", isLoading);
//   console.log("isAuthenticated", isAuthenticated);
//   console.log("session", session);
//   console.log("*** in auth router end ****");

//   // const routerContext = React.useMemo(
//   //   () => ({
//   //     isLoading,
//   //   }),
//   //   [isLoading]
//   // );

//   // React.useEffect(() => {
//   //   router.invalidate();
//   // }, [routerContext]);

//   // if (isLoading) return <div>Auth Is Loading...</div>;

//   return <RouterProvider router={router} />;
// }

// Render the app
const rootElement = document.getElementById("app")!;
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        {/* <AuthProvider>
          <AuthRouter />
        </AuthProvider> */}
        <RouterProvider router={router} />
      </QueryClientProvider>
    </ThemeProvider>
  );
}
