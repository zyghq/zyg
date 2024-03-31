"use client";

import * as React from "react";
import { useFormStatus } from "react-dom";
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
import { useToast } from "@/components/ui/use-toast";
import { ExclamationTriangleIcon, ArrowLeftIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { signup } from "./actions";
import { useRouter } from "next/navigation";

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6), // as per supabase
});

function SubmitButton({ isDisabled }) {
  return (
    <Button
      type="submit"
      disabled={isDisabled}
      aria-disabled={isDisabled}
      aria-label="Sign Up"
    >
      Sign Up
    </Button>
  );
}

export default function SignUpPage() {
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
  const { isSubmitting, errors, isSubmitted } = formState;

  async function onSubmit(values) {
    const response = await signup(values);
    const { error, ok, data } = response;
    if (error) {
      const { name, message } = error;
      console.error(name, message);
      if (name === "AuthWeakPasswordError") {
        form.setError("password", {
          type: 500,
          message: message,
        });
        return;
      }
      if (name === "AuthApiError") {
        form.setError("root.serverError", {
          type: 500,
          message: message,
        });
        return;
      } else {
        // default error message
        form.setError("root.serverError", {
          type: 500,
          message: "Please try again later.",
        });
      }
    }
    if (ok) {
      const { email } = data;
      const message = `Please verify your email address - ${email}.`;
      toast({
        description: message,
      });
      router.push("/login/");
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <Card className="mx-auto w-full max-w-sm">
          <CardHeader>
            <CardTitle>Create your Zyg account.</CardTitle>
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
                    <Input placeholder="name@example.com" {...field} />
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
              <Link href="/login/">
                <ArrowLeftIcon className="h-4 w-4 mr-1" />
                <span>Log in</span>
              </Link>
            </Button>
            <SubmitButton isDisabled={isSubmitting || isSubmitted} />
          </CardFooter>
        </Card>
      </form>
    </Form>

    // <form action={login}>
    //   <Card className="mx-auto w-full max-w-sm">
    //     <CardHeader>
    //       <CardTitle>Log in to Zyg.</CardTitle>
    //       {/* <CardDescription>Deploy your new project in one-click.</CardDescription> */}
    //     </CardHeader>
    //     <CardContent>
    //       <div className="grid w-full items-center gap-4">
    //         <div className="flex flex-col space-y-1.5">
    //           <Label htmlFor="name">Email</Label>
    //           <Input
    //             id="email"
    //             type="email"
    //             name="email"
    //             placeholder="name@example.com"
    //             required
    //             className=""
    //           />
    //         </div>
    //         <div className="flex flex-col space-y-1.5">
    //           <Label htmlFor="password">Password</Label>
    //           <Input
    //             id="password"
    //             type="password"
    //             name="password"
    //             placeholder="V3ryS3curE"
    //             required
    //           />
    //         </div>
    //       </div>
    //     </CardContent>
    //     <CardFooter className="flex justify-between">
    //       <Button variant="outline" asChild>
    //         <Link href="/signup/">Sign up</Link>
    //       </Button>
    //       <SubmitButton />
    //     </CardFooter>
    //   </Card>
    // </form>
  );
}
