import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { SendHorizonalIcon } from "lucide-react";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { useMutation } from "@tanstack/react-query";
import { sendThreadChatMessage } from "@/db/api";
import { ThreadChatWithMessages } from "@/db/entities";

const formSchema = z.object({
  message: z.string().min(1, "Message is required"),
});

type MessageFormProps = {
  token: string;
  workspaceId: string;
  threadId: string;
  customerName: string;
  refetch: () => void;
};

function SubmitButton({ isDisabled }: { isDisabled: boolean }) {
  return (
    <Button
      className="ml-1"
      size="icon"
      type="submit"
      disabled={isDisabled}
      aria-disabled={isDisabled}
      aria-label="Send Message"
    >
      <SendHorizonalIcon className="h-4 w-4" />
    </Button>
  );
}

export function MessageForm({
  token,
  workspaceId,
  threadId,
  customerName,
  refetch,
}: MessageFormProps) {
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      message: "",
    },
  });

  const mutation = useMutation({
    mutationFn: async (values: { message: string }) => {
      const { message } = values;
      const { error, data } = await sendThreadChatMessage(
        token,
        workspaceId,
        threadId,
        { message }
      );
      if (error) {
        throw new Error(error.message);
      }
      if (!data) {
        throw new Error("no data returned");
      }
      return data as ThreadChatWithMessages;
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
    mutation.mutateAsync({ message });
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex justify-between items-center"
      >
        <FormField
          control={form.control}
          name="message"
          render={({ field }) => (
            <FormItem className="space-y-2 w-full">
              <FormControl>
                <Textarea
                  className="resize-none"
                  rows={4}
                  placeholder={`Reply to ${customerName}`}
                  required
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
