import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";

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

const searchSearchSchema = z.object({
  redirect: z.string().catch("/workspaces"),
});

type FormInputs = {
  email: string;
  password: string;
};

export const Route = createFileRoute("/login")({
  validateSearch: searchSearchSchema,
  component: () => <LoginComponent />,
});

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
});

function LoginComponent() {
  const { auth, AccountStore } = Route.useRouteContext();
  const navigate = useNavigate();
  const { redirect } = Route.useSearch();
  const { client } = auth;
  const { toast } = useToast();
  const useStore = AccountStore.useContext();

  const form = useForm<FormInputs>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    const { email, password } = inputs;
    try {
      const { data, error: errSupa } = await client.auth.signInWithPassword({
        email: email,
        password: password,
      });
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
      console.log("account", account);
      useStore.setState((prev) => ({
        ...prev,
        hasData: true,
        error: null,
        account: account,
      }));
      toast({
        description: `Welcome back, ${account.email}. You are now logged in.`,
      });
      return navigate({ to: redirect });
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
              <CardTitle>Log in to Zyg.</CardTitle>
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
              <Button variant="outline" asChild>
                <Link to="/signup">
                  <ArrowLeftIcon className="mr-1 h-4 w-4" />
                  <span>Sign Up</span>
                </Link>
              </Button>
              <Button
                type="submit"
                disabled={isSubmitting || isSubmitSuccessful}
                aria-disabled={isSubmitting || isSubmitSuccessful}
                aria-label="Log In"
              >
                Login
              </Button>
            </CardFooter>
            <CardFooter className="flex justify-center">
              <Button variant="link" asChild>
                <Link href="/recover/">Forgot Password?</Link>
              </Button>
            </CardFooter>
          </Card>
        </form>
      </Form>
    </div>
  );
}
