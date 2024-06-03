import React from "react";
import { z } from "zod";
import {
  createFileRoute,
  Link,
  useRouterState,
  useRouter,
} from "@tanstack/react-router";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { ArrowLeftIcon } from "@radix-ui/react-icons";
import { useAuth } from "@/auth";

const searchSearchSchema = z.object({
  redirect: z.string().optional().catch("/signin"),
});

export const Route = createFileRoute("/signout")({
  validateSearch: searchSearchSchema,
  component: SignOutComponent,
});

function SignOutComponent() {
  const auth = useAuth();
  const router = useRouter();
  const navigate = Route.useNavigate();
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const [isError, setIsError] = React.useState(false);

  async function confirmSignOut() {
    const { error } = await auth.client.auth.signOut();
    if (error) {
      setIsError(true);
      return;
    }
    await router.invalidate();
    await navigate({ to: "/" });
  }

  return (
    <div className="flex min-h-screen flex-col justify-center p-4">
      <Card className="mx-auto w-full max-w-sm">
        <CardHeader>
          <CardTitle>Sign out</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-muted-foreground">
            Are you sure you want to sign out?
          </p>
        </CardContent>
        <CardFooter className="flex justify-between">
          <Button variant="outline" aria-label="Log In" asChild>
            <Link to="/workspaces">
              <ArrowLeftIcon className="mr-1 h-4 w-4" />
              <span>Workspaces</span>
            </Link>
          </Button>
          <Button
            aria-label="Sign Out"
            onClick={() => confirmSignOut()}
            disabled={isLoading}
          >
            Yes, I'll be back
          </Button>
        </CardFooter>
        {isError && (
          <CardFooter className="flex justify-center">
            <p className="text-red-500 text-sm">
              Something went wrong. Please try again later.
            </p>
          </CardFooter>
        )}
      </Card>
    </div>
  );
}
