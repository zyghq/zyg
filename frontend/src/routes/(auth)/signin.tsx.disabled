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

const searchSearchSchema = z.object({
  redirect: z.string().optional().catch(""),
});

type FormInputs = {
  email: string;
  password: string;
};

const fallback = "/workspaces" as const;

export const Route = createFileRoute("/(auth)/signin")({
  beforeLoad: async ({ context, search }) => {
    const { supabaseClient } = context;
    const { data, error } = await supabaseClient.auth.getSession();
    const isAuthenticated = !error && data?.session;
    if (isAuthenticated) {
      throw redirect({ to: search.redirect || fallback });
    }
  },
  component: SignInComponent,
  validateSearch: searchSearchSchema,
});

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
});

function SignInComponent() {
  const { supabaseClient } = Route.useRouteContext();
  const router = useRouter();
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const navigate = useNavigate();
  const search = Route.useSearch();
  const { toast } = useToast();

  const form = useForm<FormInputs>({
    defaultValues: {
      email: "",
      password: "",
    },
    resolver: zodResolver(formSchema),
  });

  const { formState } = form;
  const { errors, isSubmitSuccessful, isSubmitting } = formState;

  const isLoggingIn = isLoading || isSubmitting;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    try {
      const { email, password } = inputs;
      const { data, error: errSupabase } =
        await supabaseClient.auth.signInWithPassword({ email, password });
      if (errSupabase) {
        console.error(errSupabase);
        const { message, name } = errSupabase;
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
      const { session } = data;
      if (session) {
        const { access_token: token } = session;
        const { data: account, error: errAccount } =
          await getOrCreateZygAccount(token);
        if (errAccount || !account) {
          console.error(errAccount);
          form.setError("root", {
            message: "Something went wrong. Please try again later.",
            type: "serverError",
          });
          return;
        }
        toast({
          description: `Welcome back, ${account.email}. You are now signed in.`,
        });
        await router.invalidate();
        await navigate({ to: search.redirect || fallback });
      }
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
              <CardTitle>Sign in to Zyg.</CardTitle>
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
              <Button aria-label="Sign Up" asChild variant="outline">
                <Link preload={false} to="/signup">
                  Sign Up
                </Link>
              </Button>
              <Button
                aria-disabled={isLoggingIn || isSubmitSuccessful}
                aria-label="Submit"
                disabled={isLoggingIn || isSubmitSuccessful}
                type="submit"
              >
                Submit
              </Button>
            </CardFooter>
            <CardFooter className="flex justify-center">
              <Button asChild variant="link">
                <Link preload={false} to="/recover">
                  Forgot Password?
                </Link>
              </Button>
            </CardFooter>
          </Card>
        </form>
      </Form>
    </div>
  );
}
