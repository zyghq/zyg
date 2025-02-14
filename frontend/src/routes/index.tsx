import { createFileRoute, redirect } from "@tanstack/react-router";
import { createClient } from "@workos-inc/authkit-js";
import { useAuth } from "@workos-inc/authkit-react";
import { Button } from "@/components/ui/button";
import { Link } from "@tanstack/react-router";
import { Spinner } from "@/components/spinner";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export const Route = createFileRoute("/")({
  beforeLoad: async () => {
    const authKit = await createClient(import.meta.env.VITE_WORKOS_CLIENT_ID);
    const user = authKit.getUser();
    if (!user) {
      throw redirect({ to: "/signin" });
    }
    const token = await authKit.getAccessToken();
    return {
      user,
      token,
    };
  },
  component: IndexComponent,
});

function IndexComponent() {
  const { user = null } = Route.useRouteContext();
  const { isLoading, signIn } = useAuth();

  if (isLoading) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-4 bg-gray-50">
        <Spinner
          size={32}
          className="h-5 w-5 animate-spin text-muted-foreground"
        />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50 p-4">
      {user ? (
        <>
          <p className="mb-2 text-xl font-medium">Welcome</p>
          <Card className="w-96 shadow-md">
            <CardContent className="pt-6">
              <div className="flex items-center space-x-4">
                <Avatar className="h-12 w-12">
                  <AvatarImage
                    src={user.profilePictureUrl ?? undefined}
                    alt={user.firstName ?? undefined}
                  />
                  <AvatarFallback>{user.firstName}</AvatarFallback>
                </Avatar>
                <div>
                  <h2 className="text-xl font-normal">{user.firstName}</h2>
                  <p className="text-sm text-muted-foreground">{user.email}</p>
                </div>
              </div>
            </CardContent>
            <CardFooter>
              <Button variant="default" asChild className="w-full">
                <Link to="/workspaces">Go to Workspaces</Link>
              </Button>
            </CardFooter>
          </Card>
        </>
      ) : (
        <Button variant="ghost" onClick={() => signIn()}>
          Sign In
        </Button>
      )}
    </div>
  );
}
