import { createFileRoute, redirect } from "@tanstack/react-router";
import { createClient } from "@workos-inc/authkit-js";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuth } from "@workos-inc/authkit-react";

export const Route = createFileRoute("/(auth)/signin")({
  beforeLoad: async () => {
    const authKit = await createClient(import.meta.env.VITE_WORKOS_CLIENT_ID);
    const user = authKit.getUser();
    if (user) {
      throw redirect({ to: "/" });
    }
  },
  component: SignInComponent,
});

function SignInComponent() {
  const { signIn } = useAuth();
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <Card className="w-96 shadow-md">
        <CardHeader>
          <CardTitle className="text-center text-2xl font-semibold">
            You've been logged out
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-balance text-center text-muted-foreground">
            Your session has ended. Please sign in again to continue.
          </p>
        </CardContent>
        <CardFooter className="flex justify-center">
          <Button onClick={() => signIn()}>Sign In</Button>
        </CardFooter>
      </Card>
    </div>
  );
}
