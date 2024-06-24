import {
  createFileRoute,
  redirect,
  Link,
  useNavigate,
  useRouter,
  useRouterState,
} from "@tanstack/react-router";

import { zodResolver } from "@hookform/resolvers/zod";

import { useForm, SubmitHandler } from "react-hook-form";
import { z } from "zod";

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

import { ExclamationTriangleIcon } from "@radix-ui/react-icons";
import { getOrCreateZygAccount } from "@/db/api";

const searchSearchSchema = z.object({
  redirect: z.string().optional().catch(""),
});

type FormInputs = {
  email: string;
  password: string;
};

const fallback = "/workspaces" as const;

export const Route = createFileRoute("/signin")({
  validateSearch: searchSearchSchema,
  beforeLoad: async ({ context, search }) => {
    const { supabaseClient } = context;
    const { error, data } = await supabaseClient.auth.getSession();
    const isAuthenticated = !error && data?.session;
    if (isAuthenticated) {
      throw redirect({ to: search.redirect || fallback });
    }
  },
  component: SignInComponent,
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
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  const isLoggingIn = isLoading || isSubmitting;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    try {
      const { email, password } = inputs;
      const { error: errSupa, data } =
        await supabaseClient.auth.signInWithPassword({ email, password });
      if (errSupa) {
        console.error(errSupa);
        const { name, message } = errSupa;
        if (name === "AuthApiError") {
          form.setError("root", {
            type: "authError",
            message: message,
          });
          return;
        } else {
          form.setError("root", {
            type: "serverError",
            message: "Please try again later.",
          });
          return;
        }
      }
      const { session } = data;
      if (session) {
        const { access_token: token } = session;
        const { error: errAccount, data: account } =
          await getOrCreateZygAccount(token);
        if (errAccount || !account) {
          console.error(errAccount);
          form.setError("root", {
            type: "serverError",
            message: "Something went wrong. Please try again later.",
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
        type: "serverError",
        message: "Something went wrong. Please try again later.",
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
                        type="email"
                        placeholder="name@example.com"
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
                        type="password"
                        placeholder="VeryS3Cure"
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
              <Button variant="outline" aria-label="Sign Up" asChild>
                <Link to="/signup" preload={false}>
                  Sign Up
                </Link>
              </Button>
              <Button
                type="submit"
                disabled={isLoggingIn || isSubmitSuccessful}
                aria-disabled={isLoggingIn || isSubmitSuccessful}
                aria-label="Submit"
              >
                Submit
              </Button>
            </CardFooter>
            <CardFooter className="flex justify-center">
              <Button variant="link" asChild>
                <Link to="/recover" preload={false}>
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
