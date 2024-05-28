"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { ExclamationTriangleIcon, ArrowLeftIcon } from "@radix-ui/react-icons";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { reset } from "./actions";

const formSchema = z.object({
  password: z.string().min(6),
  confirm: z.string().min(6),
});

function SubmitButton({ isDisabled }) {
  return (
    <Button
      className="w-full"
      type="submit"
      disabled={isDisabled}
      aria-disabled={isDisabled}
      aria-label="Confirm Reset Password"
    >
      Reset
    </Button>
  );
}

export default function PasswordResetPage() {
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      password: "",
      confirm: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  async function onSubmit(values) {
    const { password, confirm } = values;
    if (password !== confirm) {
      form.setError("confirm", {
        type: "custom",
        message: "Passwords do not match.",
      });
      return;
    }
    const body = {
      password: confirm,
    };
    const response = await reset(body);
    const { error, ok } = response;
    if (error) {
      const { name, message } = error;
      console.log(name, ":", message);
      if (name === "AuthApiError") {
        form.setError("root.serverError", {
          type: 500,
          message: message,
        });
        return;
      }
      if (name === "AuthSessionMissingError") {
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
    if (ok) {
      redirect("/");
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Card className="mx-auto w-full max-w-sm">
          <CardHeader>
            <CardTitle>Reset account password.</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {errors?.root?.serverError?.type === 500 && (
              <Alert variant="destructive">
                <ExclamationTriangleIcon className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>
                  {`${
                    errors?.root?.serverError?.message ||
                    "Please try agian later."
                  }`}
                </AlertDescription>
              </Alert>
            )}
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
                  {/* <FormDescription>A valid email address</FormDescription> */}
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="confirm"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Confirm Password</FormLabel>
                  <FormControl>
                    <Input
                      type="password"
                      placeholder="VeryS3curE Confirm"
                      {...field}
                      required
                    />
                  </FormControl>
                  {/* <FormDescription>A valid email address</FormDescription> */}
                  <FormMessage />
                </FormItem>
              )}
            />
          </CardContent>
          <CardFooter className="flex justify-between">
            <SubmitButton isDisabled={isSubmitting || isSubmitSuccessful} />
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
