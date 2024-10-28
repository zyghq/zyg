import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { sendThreadChatMessage } from "@/db/api";
import { ThreadChatResponse } from "@/db/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { SendHorizonalIcon } from "lucide-react";
import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

const formSchema = z.object({
  message: z.string().min(1, "Message is required"),
});

type MessageFormProps = {
  customerName: string;
  refetch: () => void;
  threadId: string;
  token: string;
  workspaceId: string;
};

function SubmitButton({ isDisabled }: { isDisabled: boolean }) {
  return (
    <Button
      aria-disabled={isDisabled}
      aria-label="Send Message"
      className="ml-1"
      disabled={isDisabled}
      size="icon"
      type="submit"
    >
      <SendHorizonalIcon className="h-4 w-4" />
    </Button>
  );
}

export function MessageForm({
  customerName,
  refetch,
  threadId,
  token,
  workspaceId,
}: MessageFormProps) {
  const form = useForm({
    defaultValues: {
      message: "",
    },
    resolver: zodResolver(formSchema),
  });

  const mutation = useMutation({
    mutationFn: async (values: { message: string }) => {
      const { message } = values;
      const { data, error } = await sendThreadChatMessage(
        token,
        workspaceId,
        threadId,
        { message },
      );
      if (error) {
        throw new Error(error.message);
      }
      if (!data) {
        throw new Error("no data returned");
      }
      return data as ThreadChatResponse;
    },
    onError: (error) => {
      console.error(error);
      form.setError("message", {
        message: "Something went wrong. Please try again later.",
      });
    },
    onSuccess: (data) => {
      console.log("onSuccess");
      console.log(data);
      refetch();
      form.reset({ message: "" });
    },
  });

  const { formState } = form;
  const { isSubmitting } = formState;

  const onEnterPress = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && e.shiftKey === false) {
      // Capitalize "Enter" correctly
      e.preventDefault();
      (e.target as HTMLTextAreaElement).form?.requestSubmit(); // Cast e.target to HTMLTextAreaElement
    }
  };

  async function onSubmit(values: { message: string }) {
    const { message } = values;
    await mutation.mutateAsync({ message });
  }

  return (
    <Form {...form}>
      <form
        className="flex items-center justify-between"
        onSubmit={form.handleSubmit(onSubmit)}
      >
        <FormField
          control={form.control}
          name="message"
          render={({ field }) => (
            <FormItem className="w-full space-y-2">
              <FormControl>
                <Textarea
                  className="resize-none"
                  placeholder={`Reply to ${customerName}`}
                  required
                  rows={4}
                  {...field}
                  onKeyDown={onEnterPress}
                />
              </FormControl>
            </FormItem>
          )}
        />
        <SubmitButton isDisabled={isSubmitting} />
      </form>
    </Form>
  );
}
