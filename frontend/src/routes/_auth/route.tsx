import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";
import { getOrCreateZygAccount } from "@/db/api";
import { AccountStoreProvider } from "@/providers";

const accountQueryOptions = (token: string) =>
  queryOptions({
    queryKey: ["account", token],
    queryFn: async () => {
      console.log("*** accountQueryOptions start ****");
      console.log("*** token ****", token);
      console.log("*** accountQueryOptions end ****");
      const { error, data } = await getOrCreateZygAccount(token);
      if (error) throw new Error("failed to authenticate user account.");
      return data;
    },
  });

export const Route = createFileRoute("/_auth")({
  beforeLoad: async ({ context }) => {
    console.log("*** beforeLoad start ****");
    const { supabaseClient } = context;
    const { error, data } = await supabaseClient.auth.getSession();
    if (error || !data?.session) {
      throw redirect({ to: "/signin" });
    }

    console.log("*** error ****", error);
    console.log("*** data ****", data);

    const token = data.session.access_token;

    console.log("*** token ****", token);
    console.log("*** beforeLoad end ****");
    return { token };
  },
  loader: async ({ context: { queryClient, supabaseClient } }) => {
    const { error, data } = await supabaseClient.auth.getSession();
    if (error || !data?.session) throw redirect({ to: "/signin" });
    const token = data.session.access_token;
    console.log("*********** Token in loader ***********", token);
    return queryClient.ensureQueryData(accountQueryOptions(token));
  },
  component: AuthLayout,
});

function AuthLayout() {
  const { token } = Route.useRouteContext();
  // const { data } = useSuspenseQuery(accountQueryOptions(token));

  console.log("*********** Token in AuthLayout ***********", token);

  return (
    <Outlet />
    // <AccountStoreProvider
    //   initialValue={{
    //     hasData: data ? true : false,
    //     error: data ? null : new Error("failed to fetch account details"),
    //     token,
    //     account: data
    //       ? {
    //           email: data.email,
    //           accountId: data.accountId,
    //           name: data.name,
    //           provider: data.provider,
    //           createdAt: data.createdAt,
    //           updatedAt: data.updatedAt,
    //         }
    //       : null,
    //   }}
    // >
    //   <Outlet />
    // </AccountStoreProvider>
  );
}
