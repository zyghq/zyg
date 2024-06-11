import {
  createFileRoute,
  Link,
  redirect,
  useRouterState,
} from "@tanstack/react-router";
import { Card, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { ExitIcon } from "@radix-ui/react-icons";
import { Icons } from "@/components/icons";
import { buttonVariants } from "@/components/ui/button";
import { queryOptions } from "@tanstack/react-query";

type Workspace = {
  workspaceId: string;
  accountId: string;
  name: string;
  createdAt: string;
  updatedAt: string;
};

async function fetchWorkspaces(token: string) {
  try {
    const response = await fetch(
      `${import.meta.env.VITE_ZYG_URL}/workspaces/`,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );
    // handle 4xx-5xx errors
    if (!response.ok) {
      const { status, statusText } = response;
      throw new Error(`error fetching workspaces: ${status} ${statusText}`);
    }

    // Ok
    return await response.json();
  } catch (err) {
    console.log("*** error fetching workspaces in fetch API call ***", err);
    console.error(err);
    throw new Error("error fetching workspaces");
  }
}

const workspacesQueryOptions = (token: string) =>
  queryOptions({
    queryKey: ["workpaces", token],
    queryFn: async () => {
      return await fetchWorkspaces(token);
    },
    enabled: !!token,
  });

// TODO: do error handling
// https://tanstack.com/router/latest/docs/framework/react/guide/external-data-loading#error-handling-with-tanstack-query
export const Route = createFileRoute("/_auth/workspaces/")({
  loader: async ({ context: { queryClient, supabaseClient } }) => {
    const { error, data } = await supabaseClient.auth.getSession();
    if (error || !data?.session) throw redirect({ to: "/signin" });
    const token = data.session.access_token as string;
    return queryClient.ensureQueryData(workspacesQueryOptions(token));
  },
  component: Workspaces,
});

function Workspaces() {
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const workspaces: Workspace[] = Route.useLoaderData();

  return (
    <div className="relative flex min-h-screen flex-col bg-background">
      <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 max-w-screen-2xl items-center">
          <div className="flex flex-1 items-end">
            <Icons.logo className="my-auto mr-2 h-5 w-5" />
            <span className="font-semibold">Zyg.</span>
          </div>
          <div className="flex justify-between space-x-2 md:justify-end">
            <Link
              to="/signout"
              preload={false}
              className={buttonVariants({ size: "icon", variant: "outline" })}
            >
              <ExitIcon />
            </Link>
          </div>
        </div>
      </header>
      <div className="container py-4">
        <h1 className="mb-1 text-3xl font-bold">Create a Workspace</h1>
        <p className="mb-4 text-foreground">
          Use workspaces to organize your projects and teams, or separate your
          live and test environments.
          <br /> Customers and team members are specific to a workspace.
        </p>
        <Button variant="default" asChild disabled={isLoading}>
          <Link to={"/workspaces/add"}>Create Workspace</Link>
        </Button>
        {workspaces && workspaces.length > 0 && (
          <>
            <Separator className="my-4 md:w-1/3" />
            <div className="text-lg font-semibold">Open a Workspace</div>
            <div className="mt-4 space-y-2 md:w-3/5 lg:w-2/5">
              {workspaces.map((workspace) => (
                <Card key={workspace.workspaceId}>
                  <CardHeader>
                    <CardTitle>{workspace.name}</CardTitle>
                  </CardHeader>
                  <CardFooter className="justify-end">
                    <Button asChild>
                      <Link
                        to={"/workspaces/$workspaceId"}
                        params={{ workspaceId: workspace.workspaceId }}
                        search={{ status: "todo" }}
                      >
                        Open
                      </Link>
                    </Button>
                  </CardFooter>
                </Card>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
