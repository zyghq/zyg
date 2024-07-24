import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { SubmitHandler } from "react-hook-form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { SendHorizonalIcon } from "lucide-react";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";

const formSchema = z.object({
  email: z.string().email(),
});

type FormValues = z.infer<typeof formSchema>;

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

export default function AskEmailForm({
  widgetId,
  threadId,
  jwt,
}: {
  widgetId: string;
  threadId: string;
  jwt: string;
}) {
  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
    },
  });

  const { formState } = form;
  const { isSubmitting, errors } = formState;

  const onSubmit: SubmitHandler<FormValues> = async (values) => {
    console.log(values);
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex justify-between items-center"
      >
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem className="space-y-2 w-full">
              <FormControl>
                <Input
                  placeholder="you@example.com"
                  title="Email"
                  required
                  {...field}
                />
              </FormControl>
              {errors?.root?.serverError && (
                <FormMessage>{errors?.root?.serverError?.message}</FormMessage>
              )}
            </FormItem>
          )}
        />
        <SubmitButton isDisabled={isSubmitting} />
      </form>
    </Form>
  );
}
