// import React from "react";
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

import { ArrowLeftIcon, ExclamationTriangleIcon } from "@radix-ui/react-icons";
import { getOrCreateZygAccount } from "@/db/api";
// import { Session } from "@supabase/supabase-js";
// import { useAuth } from "@/auth";

const searchSearchSchema = z.object({
  redirect: z.string().optional().catch(""),
});

type FormInputs = {
  email: string;
  password: string;
};

// const supabaseClient = createClient(
//   import.meta.env.VITE_SUPABASE_URL,
//   import.meta.env.VITE_SUPABASE_ANON_KEY
// );

const fallback = "/workspaces" as const;

export const Route = createFileRoute("/signin")({
  validateSearch: searchSearchSchema,
  beforeLoad: async ({ context, search }) => {
    console.log("**** beforeLoad in signin ****");
    const { supabaseClient } = context;
    const { error, data } = await supabaseClient.auth.getSession();
    console.log("**** error ****", error);
    console.log("**** data ****", data);

    const isAuthenticated = !error && data?.session;
    if (isAuthenticated) {
      throw redirect({ to: search.redirect || fallback });
    }
    // const { isAuthenticated } = context;
    // if (isAuthenticated) throw redirect({ to: search.redirect || fallback });
    // const { auth } = context;
    // const session = await auth?.client.auth.getSession();
    // const { error: errSupa, data } = session;
    // if (!errSupa && data?.session) {
    //   throw redirect({ to: search.redirect || fallback });
    // }
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

  // const [session, setSession] = React.useState<Session | null>(null);

  // React.useEffect(() => {
  //   const setData = async () => {
  //     const {
  //       data: { session },
  //       error,
  //     } = await supabaseClient.auth.getSession();
  //     if (error) {
  //       throw error;
  //     }
  //     setSession(session);
  //   };

  //   setData();
  // }, [supabaseClient]);

  // const [signedIn, setSignedIn] = React.useState(false);
  // const [isSigningIn, setIsSigningIn] = React.useState(false);

  // const [, setSession] = React.useState<Session | null>(null);
  // const [isSigningIn, setIsSigningIn] = React.useState(true);

  // React.useEffect(() => {
  //   const { data: listener } = supabaseClient.auth.onAuthStateChange(
  //     (_event, session) => {
  //       setSession(session);
  //       setIsSigningIn(false);
  //     }
  //   );

  //   const setData = async () => {
  //     const {
  //       data: { session },
  //       error,
  //     } = await supabaseClient.auth.getSession();
  //     if (error) {
  //       throw error;
  //     }

  //     setSession(session);
  //     setIsSigningIn(false);

  //     // if (session) {
  //     //   await router.invalidate();
  //     //   await navigate({ to: search.redirect || fallback });
  //     // }
  //   };

  //   setData();

  //   return () => {
  //     listener?.subscription.unsubscribe();
  //   };
  // }, [supabaseClient, navigate, router, search]);

  // React.useEffect(() => {
  //   const doRedirect = async () => {
  //     if (session) {
  //       await router.invalidate();
  //       await navigate({ to: search.redirect || fallback });
  //     }
  //   };
  //   if (session) {
  //     toast({
  //       description: `Welcome back, ${session.user.email}. You are now signed in.`,
  //     });
  //     doRedirect();
  //   }
  // }, [session, router, navigate, search, toast]);

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

  // async function doRedirect() {
  //   if (signedIn) {
  //     await router.invalidate();
  //     await navigate({ to: search.redirect || fallback });
  //   }
  // }

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
                    {`${errors?.root?.message || "Please try agian later."}`}
                  </AlertDescription>
                </Alert>
              )}
              {errors?.root?.type === "serverError" && (
                <Alert variant="destructive">
                  <ExclamationTriangleIcon className="h-4 w-4" />
                  <AlertTitle>Error</AlertTitle>
                  <AlertDescription>
                    {`${errors?.root?.message || "Please try agian later."}`}
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
                        placeholder="VeryS3curE"
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
                  <ArrowLeftIcon className="mr-1 h-4 w-4 my-auto" />
                  Sign Up
                </Link>
              </Button>
              <Button
                type="submit"
                disabled={isLoggingIn || isSubmitSuccessful}
                aria-disabled={isLoggingIn || isSubmitSuccessful}
                aria-label="Sign In"
              >
                Sign In
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
