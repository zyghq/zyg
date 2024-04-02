"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

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
import { Input } from "@/components/ui/input";
import { ClipboardCopyIcon } from "@radix-ui/react-icons";

const formSchema = z.object({
  username: z.string().min(2, {
    message: "Username must be at least 2 characters.",
  }),
});

export function WorkspaceEditForm({ workspaceId }) {
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
    },
  });

  function onSubmit(values) {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    console.log(values);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="Zyg" {...field} />
              </FormControl>
              <FormDescription>
                Typically your company or team name.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <div>
          <div className="text-sm font-normal">Workspace ID</div>
          <div className="text-xs text-muted-foreground">
            We might ask you for support inquiries.
          </div>
          <div className="mt-4 flex w-full items-center">
            <code className="mr-2 rounded-lg border bg-muted p-2">
              {workspaceId}
            </code>
            <Button variant="ghost" size="sm">
              <ClipboardCopyIcon className="h-4 w-4" />
            </Button>
          </div>
        </div>
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
}

const labelFormSchema = z.object({
  label: z.string().min(2, {
    label: "Label must be at least 2 characters.",
  }),
});

// We can use the same form for adding and editing a label.
// but edit could be displayed inside as card so it can be easily idenfied
// that it is a edit rather than adding a new label.
export function WorkspaceLabelAddOrEditForm({ workspaceId }) {
  const form = useForm({
    resolver: zodResolver(labelFormSchema),
    defaultValues: {
      username: "",
    },
  });

  function onSubmit(values) {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    console.log(values);
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex items-center space-x-2"
      >
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem className="w-full">
              <FormControl>
                <Input placeholder="Label" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button size="sm" type="submit">
          Submit
        </Button>
      </form>
    </Form>
  );
}
