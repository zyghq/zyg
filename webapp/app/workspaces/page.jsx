import { redirect } from "next/navigation";
import { createClient } from "@/utils/supabase/server";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Card, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { isAuthenticated, getAuthToken } from "@/utils/supabase/helpers";

async function getAccountWorkspaces(authToken) {
  console.log(authToken);
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
      },
    );
    // handle 4xx-5xx errors
    if (!response.ok) {
      const { status, statusText } = response;
      throw new Error(`Error fetching workspaces: ${status} ${statusText}`);
    }
    const data = await response.json();
    return [null, data];
  } catch (err) {
    console.error(err);
    return [err, null];
  }
}

export default async function SelectWorkspacePage() {
  const supabase = createClient();

  if (!(await isAuthenticated(supabase))) {
    return redirect("/login/");
  }

  const authToken = await getAuthToken(supabase);
  const [err, workspaces] = await getAccountWorkspaces(authToken);

  if (err) {
    return (
      <div className="container">
        <h1 className="mb-1 text-3xl font-bold">Error</h1>
        <p className="mb-4 text-red-500">
          There was an error fetching your workspaces. Please try again later.
        </p>
      </div>
    );
  }

  return (
    <div className="container">
      <h1 className="mb-1 text-3xl font-bold">Create a Workspace</h1>
      <p className="mb-4 text-zinc-500">
        Use workspaces to organize your projects and teams, or separate your
        live and test environments.
        <br /> Customers and team members are specific to a workspace.
      </p>
      <Button asChild>
        <Link href="/workspaces/add/">Create Workspace</Link>
      </Button>
      <Separator className="my-4 md:w-1/3" />
      <div className="font-semibold text-zinc-800">Open a Workspace</div>
      <div className="mt-4 space-y-2 md:w-3/5 lg:w-2/5">
        {workspaces.map((workspace) => (
          <Card key={workspace.workspaceId}>
            <CardHeader>
              <CardTitle>{workspace.name}</CardTitle>
            </CardHeader>
            <CardFooter className="justify-end">
              <Button asChild>
                <Link href={`/${workspace.workspaceId}/`}>Open</Link>
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  );
}
