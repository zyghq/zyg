"use client";

import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { SendHorizonalIcon } from "lucide-react";
import { sendThreadChatMessage } from "@/app/threads/_actions";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { useRouter } from "next/navigation";

const formSchema = z.object({
  message: z.string().min(1, "Message is required"),
});

function SubmitButton({ isDisabled }) {
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

export default function MessageThreadForm({ threadId }) {
  //   const router = useRouter();
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      message: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors, isSubmitSuccessful } = formState;

  const onEnterPress = (e) => {
    if (e.keyCode === 13 && e.shiftKey === false) {
      e.preventDefault();
      e.target.form.requestSubmit();
    }
  };

  async function onSubmit(values) {
    const response = await sendThreadChatMessage(threadId, values);
    const { error, data } = response;
    if (error) {
      console.log("got error from server....", error);
      const { message } = error;
      form.setError("root.serverError", {
        type: 500,
        message:
          message || "Failed to create the chat. Please try again later.",
      });
      return;
    }
    console.log(data);
    // const { threadId } = data;
    // return router.push(`/threads/${threadId}/`);
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
                  placeholder="Send us a message"
                  title="Send us a message"
                  required
                  {...field}
                  onKeyDown={onEnterPress}
                />
              </FormControl>
              {errors?.root?.serverError?.type === 500 && (
                <FormMessage type="error">
                  {errors?.root?.serverError?.message}
                </FormMessage>
              )}
            </FormItem>
          )}
        />
        <SubmitButton isDisabled={isSubmitting || isSubmitSuccessful} />
      </form>
    </Form>
  );
}
