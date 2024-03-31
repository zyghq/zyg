"use client";

import Link from "next/link";
import { useToast } from "@/components/ui/use-toast";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
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
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ExclamationTriangleIcon, ArrowRightIcon } from "@radix-ui/react-icons";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { useRouter } from "next/navigation";
import { createWorkspace } from "./actions";

const formSchema = z.object({
  name: z.string().min(3),
});

function SubmitButton({ isDisabled }) {
  return (
    <Button
      type="submit"
      disabled={isDisabled}
      aria-disabled={isDisabled}
      aria-label="Log In"
    >
      Save & Continue
      <ArrowRightIcon className="h-4 w-4 ml-2" />
    </Button>
  );
}

export default function AddWorkspacePage() {
  const router = useRouter();
  const { toast } = useToast();

  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  async function onSubmit(values) {
    const response = await createWorkspace(values);
    const { error, data } = response;
    if (error) {
      const { message } = error;
      console.log(message);
      form.setError("root.serverError", {
        type: 500,
        message: message || "Something went wrong. Please try again later.",
      });
      return;
    }
    const { slug } = data;
    toast({
      description: "Workspace created successfully.",
    });
    // router.push(`/workspaces/${slug}/`);
    console.log(`successfully created workspace with slug: ${slug}`);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="mt-24">
        <Card className="mx-auto w-full max-w-sm">
          <CardHeader>
            <CardTitle>Create a new Workspace.</CardTitle>
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
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input type="name" {...field} required />
                  </FormControl>
                  <FormDescription>
                    Typically your company name.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </CardContent>
          <CardFooter className="flex justify-between">
            <SubmitButton isDisabled={isSubmitting || isSubmitSuccessful} />
          </CardFooter>
          <CardFooter className="flex justify-center">
            <Button variant="link" asChild>
              <Link href="/workspaces/">Go Back</Link>
            </Button>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
}
