import { getSession, isAuthenticated } from "@/utils/supabase/helpers";
import { createClient } from "@/utils/supabase/server";

import Link from "next/link";
import { redirect } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Card, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

async function getAccountWorkspacesAPI(authToken) {
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
      }
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

  const { token, error } = await getSession(supabase);
  if (error) {
    return redirect("/login/");
  }
  const [err, workspaces] = await getAccountWorkspacesAPI(token);

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
      <p className="mb-4 text-foreground">
        Use workspaces to organize your projects and teams, or separate your
        live and test environments.
        <br /> Customers and team members are specific to a workspace.
      </p>
      <Button variant="default" asChild>
        <Link href="/workspaces/add/">Create Workspace</Link>
      </Button>
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
                <Link href={`/${workspace.workspaceId}/`}>Open</Link>
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  );
}
