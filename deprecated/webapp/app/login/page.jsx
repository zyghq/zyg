"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import * as React from "react";
import { useFormStatus } from "react-dom";
import { useForm } from "react-hook-form";
import { z } from "zod";

import Link from "next/link";
import { useRouter } from "next/navigation";

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
import { Label } from "@/components/ui/label";
import { useToast } from "@/components/ui/use-toast";

import { ArrowLeftIcon, ExclamationTriangleIcon } from "@radix-ui/react-icons";

import { login } from "./actions";

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
});

function SubmitButton({ isDisabled }) {
  return (
    <Button
      type="submit"
      disabled={isDisabled}
      aria-disabled={isDisabled}
      aria-label="Log In"
    >
      Login
    </Button>
  );
}

export default function LoginPage() {
  const router = useRouter();
  const { toast } = useToast();
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  async function onSubmit(values) {
    const response = await login(values);
    const { error, ok, data } = response;
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
    if (ok) {
      const { email } = data;
      const message = `Welcome back, ${email}. You are now logged in.`;
      toast({
        description: message,
      });
      router.push("/workspaces/");
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Card className="mx-auto w-full max-w-sm">
          <CardHeader>
            <CardTitle>Log in to Zyg.</CardTitle>
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
          </CardContent>
          <CardFooter className="flex justify-between">
            <Button variant="outline" asChild>
              <Link href="/signup/">
                <ArrowLeftIcon className="mr-1 h-4 w-4" />
                <span>Sign Up</span>
              </Link>
            </Button>
            <SubmitButton isDisabled={isSubmitting || isSubmitSuccessful} />
          </CardFooter>
          <CardFooter className="flex justify-center">
            <Button variant="link" asChild>
              <Link href="/recover/">Forgot Password?</Link>
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
