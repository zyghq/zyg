import { Icons } from "@/components/icons";
import { Button, buttonVariants } from "@/components/ui/button";
import { Card, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { defaultSortKey } from "@/db/store";
import { ExitIcon } from "@radix-ui/react-icons";
import { queryOptions } from "@tanstack/react-query";
import {
  createFileRoute,
  Link,
  redirect,
  useRouterState,
} from "@tanstack/react-router";
import { MoveUpRightIcon } from "lucide-react";

type Workspace = {
  accountId: string;
  createdAt: string;
  name: string;
  updatedAt: string;
  workspaceId: string;
};

const fetchWorkspaces = async (token: string): Promise<Workspace[]> => {
  const response = await fetch(`${import.meta.env.VITE_ZYG_URL}/workspaces/`, {
    headers: { Authorization: `Bearer ${token}` },
    method: "GET",
  });

  if (!response.ok) {
    throw new Error(
      `Error fetching workspaces: ${response.status} ${response.statusText}`,
    );
  }

  return response.json();
};

const workspacesQueryOptions = (token: string) =>
  queryOptions({
    enabled: !!token,
    queryFn: () => fetchWorkspaces(token),
    queryKey: ["workspaces", token],
  });

export const Route = createFileRoute("/_account/workspaces/")({
  component: Workspaces,
  loader: async ({ context: { queryClient, supabaseClient } }) => {
    const { data, error } = await supabaseClient.auth.getSession();
    if (error || !data?.session) throw redirect({ to: "/signin" });
    const token = data.session.access_token;
    return queryClient.ensureQueryData(workspacesQueryOptions(token));
  },
});

function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="flex h-14 items-center justify-between px-8">
        <div className="flex items-center">
          <Icons.logo className="mr-2 h-5 w-5" />
          <span className="font-semibold">Zyg.</span>
        </div>
        <Link
          className={buttonVariants({ size: "icon", variant: "outline" })}
          preload={false}
          to="/signout"
        >
          <ExitIcon />
        </Link>
      </div>
    </header>
  );
}

function WorkspaceCard({ workspace }: { workspace: Workspace }) {
  const isLoading = useRouterState({ select: (s) => s.isLoading });

  return (
    <Card className="shadow-none" key={workspace.workspaceId}>
      <CardHeader>
        <CardTitle className={"font-serif"}>{workspace.name}</CardTitle>
      </CardHeader>
      <CardFooter className="justify-end">
        <Button asChild disabled={isLoading} size="icon" variant="outline">
          <Link
            params={{ workspaceId: workspace.workspaceId }}
            search={{ sort: defaultSortKey }}
            to={"/workspaces/$workspaceId/threads/todo"}
          >
            <MoveUpRightIcon className="mr-1 h-4 w-4" />
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}

function Workspaces() {
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const workspaces = Route.useLoaderData();

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <Header />
      <main className="flex-1 overflow-y-auto">
        <div className="container mx-auto flex min-h-[calc(100vh-3.5rem)] flex-col items-center justify-center px-4 py-8">
          <div className="w-full max-w-2xl space-y-6">
            <h1 className="font-serif text-3xl font-bold">
              Create a Workspace
            </h1>
            <p className="text-muted-foreground">
              Use workspaces to organize your projects and teams, or separate
              your live and test environments.
              <br /> Your customers and team members are specific to a
              workspace.
            </p>
            <div className="flex">
              <Button asChild disabled={isLoading} variant="default">
                <Link to={"/workspaces/add"}>Create Workspace</Link>
              </Button>
            </div>
            {workspaces && workspaces.length > 0 && (
              <>
                <Separator className="my-6" />
                <div className="text-xl font-semibold font-serif">Open a Workspace</div>
                <div className="space-y-4">
                  {workspaces.map((workspace) => (
                    <WorkspaceCard
                      key={workspace.workspaceId}
                      workspace={workspace}
                    />
                  ))}
                </div>
              </>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
