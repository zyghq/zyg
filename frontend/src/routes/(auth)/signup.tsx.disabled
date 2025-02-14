import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/components/ui/use-toast";
import { getOrCreateZygAccount } from "@/db/api";
import { zodResolver } from "@hookform/resolvers/zod";
import { ExclamationTriangleIcon } from "@radix-ui/react-icons";
import {
  createFileRoute,
  Link,
  redirect,
  useNavigate,
  useRouter,
  useRouterState,
} from "@tanstack/react-router";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";

type FormInputs = {
  email: string;
  name: string;
  password: string;
};

export const Route = createFileRoute("/(auth)/signup")({
  beforeLoad: async ({ context }) => {
    const { supabaseClient } = context;
    const { data, error } = await supabaseClient.auth.getSession();
    const isAuthenticated = !error && data?.session;
    if (isAuthenticated) {
      throw redirect({ to: "/workspaces" });
    }
  },
  component: SignUpComponent,
});

const formSchema = z.object({
  email: z.string().email(),
  name: z.string().min(2),
  password: z.string().min(6),
});

function SignUpComponent() {
  const { supabaseClient } = Route.useRouteContext();
  const router = useRouter();
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const navigate = useNavigate();
  const { toast } = useToast();

  const form = useForm<FormInputs>({
    defaultValues: {
      email: "",
      name: "",
      password: "",
    },
    resolver: zodResolver(formSchema),
  });

  const { formState } = form;
  const { errors, isSubmitSuccessful, isSubmitting } = formState;

  const isSigningUp = isLoading || isSubmitting;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    try {
      const { email, name, password } = inputs;
      const { data, error: errSupabase } = await supabaseClient.auth.signUp({
        email: email,
        options: {
          data: {
            name: name,
          },
        },
        password: password,
      });
      if (errSupabase) {
        console.error(errSupabase);
        const { message, name } = errSupabase;
        if (name === "AuthWeakPasswordError") {
          form.setError("password", {
            message: message,
            type: "weakPassword",
          });
          return;
        }
        if (name === "AuthApiError") {
          form.setError("root", {
            message: message,
            type: "authError",
          });
          return;
        } else {
          form.setError("root", {
            message: "Please try again later.",
            type: "serverError",
          });
          return;
        }
      }
      if (!data || !data?.session) {
        form.setError("root", {
          message: "Please try again later.",
          type: "serverError",
        });
        return;
      }
      const { session } = data;
      const { access_token: token } = session;
      const { data: account, error: errAccount } = await getOrCreateZygAccount(
        token,
        { name }
      );
      if (errAccount || !account) {
        console.error(errAccount);
        form.setError("root", {
          message: "Something went wrong. Please try again later.",
          type: "serverError",
        });
        return;
      }

      toast({
        description: "You are now signed up!",
      });

      await router.invalidate();
      await navigate({ replace: true, to: "/workspaces" });
    } catch (err) {
      console.error(err);
      form.setError("root", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    }
  };

  return (
    <div className="flex min-h-screen flex-col justify-center p-4">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <Card className="mx-auto w-full max-w-sm">
            <CardHeader>
              <CardTitle>Create your Zyg account.</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {errors?.root?.type === "authError" && (
                <Alert variant="destructive">
                  <ExclamationTriangleIcon className="h-4 w-4" />
                  <AlertTitle>Error</AlertTitle>
                  <AlertDescription>
                    {`${errors?.root?.message || "Please try again later."}`}
                  </AlertDescription>
                </Alert>
              )}
              {errors?.root?.type === "serverError" && (
                <Alert variant="destructive">
                  <ExclamationTriangleIcon className="h-4 w-4" />
                  <AlertTitle>Error</AlertTitle>
                  <AlertDescription>
                    {`${errors?.root?.message || "Please try again later."}`}
                  </AlertDescription>
                </Alert>
              )}
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Name" {...field} required />
                    </FormControl>
                    <FormDescription>
                      This is your first and last name.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="name@example.com"
                        type="email"
                        {...field}
                        required
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="VeryS3Cure"
                        type="password"
                        {...field}
                        required
                      />
                    </FormControl>

                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
            <CardFooter className="flex justify-between">
              <Button aria-label="Log In" asChild variant="outline">
                <Link preload={false} to="/signin">
                  Sign In
                </Link>
              </Button>
              <Button
                aria-disabled={isSigningUp || isSubmitSuccessful}
                aria-label="Submit"
                disabled={isSigningUp || isSubmitSuccessful}
                type="submit"
              >
                Submit
              </Button>
            </CardFooter>
          </Card>
        </form>
      </Form>
    </div>
  );
}
