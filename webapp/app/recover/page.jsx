"use client";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  EnvelopeOpenIcon,
  ExclamationTriangleIcon,
} from "@radix-ui/react-icons";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import Link from "next/link";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { recover } from "./actions";

const formSchema = z.object({
  email: z.string(),
});

function SubmitButton({ isDisabled }) {
  return (
    <Button
      className="w-full"
      type="submit"
      disabled={isDisabled}
      aria-disabled={isDisabled}
      aria-label="Reset Password"
    >
      Reset Password
    </Button>
  );
}

export default function RecoverPasswordPage() {
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  async function onSubmit(values) {
    const { error } = await recover(values);
    if (error) {
      const { name, message } = error;
      console.log(name, message);
      if (name === "AuthApiError") {
        form.setError("root.serverError", {
          type: 500,
          message: message,
        });
        return;
      } else {
        form.setError("root.serverError", {
          type: 500,
          message: "Please try again later.",
        });
        return;
      }
    }
  }

  return (
    <div className="mx-auto max-w-sm space-y-6">
      <div className="space-y-2">
        <h1 className="text-3xl font-bold">Forgot Your Password?</h1>
        <p className="text-gray-500 dark:text-gray-400">
          Enter your email below to reset your password
        </p>
      </div>
      {isSubmitSuccessful && (
        <Alert>
          <EnvelopeOpenIcon className="h-4 w-4" />
          <AlertTitle>Heads up!</AlertTitle>
          <AlertDescription>
            We have sent you an email with a link to reset your password.
          </AlertDescription>
        </Alert>
      )}
      {errors?.root?.serverError?.type === 500 && (
        <Alert variant="destructive">
          <ExclamationTriangleIcon className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>
            {`${
              errors?.root?.serverError?.message || "Please try agian later."
            }`}
          </AlertDescription>
        </Alert>
      )}
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <div className="space-y-4">
            <div className="space-y-2">
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
                    {/* <FormDescription>A valid email address</FormDescription> */}
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <SubmitButton isDisabled={isSubmitting || isSubmitSuccessful} />
            <Link
              className="inline-block w-full text-center underline"
              href="/login/"
            >
              Go back to login
            </Link>
          </div>
        </form>
      </Form>
    </div>
  );
}
