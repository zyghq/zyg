import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import {
  createFileRoute,
  redirect,
  useRouterState,
} from "@tanstack/react-router";
import { useAuth } from "@workos-inc/authkit-react";
import { createClient } from "@workos-inc/authkit-js";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Spinner } from "@/components/spinner";

export const Route = createFileRoute("/(auth)/signout")({
  beforeLoad: async () => {
    const authKit = await createClient(import.meta.env.VITE_WORKOS_CLIENT_ID);
    const user = authKit.getUser();
    if (!user) {
      throw redirect({ to: "/signin" });
    }
  },
  component: SignOutComponent,
});

function SignOutComponent() {
  const { signOut, isLoading, user, signIn } = useAuth();
  const isRouteLoading = useRouterState({ select: (s) => s.isLoading });
  if (isLoading) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center gap-4">
        <Spinner
          size={32}
          className="h-5 w-5 animate-spin text-muted-foreground"
        />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-4">
      {user ? (
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
            <Button
              disabled={isRouteLoading || isLoading}
              variant="default"
              className="w-full"
              onClick={() => signOut()}
            >
              Sign Out
            </Button>
          </CardFooter>
        </Card>
      ) : (
        <Button variant="ghost" onClick={() => signIn()}>
          Sign In
        </Button>
      )}
    </div>
  );
}
