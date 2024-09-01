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
import { updateEmailAction } from "@/app/threads/_actions";
import { CustomerRefreshable } from "@/lib/customer";
import { customerSchema } from "@/lib/customer";

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
  jwt,
  setUpdates,
}: {
  widgetId: string;
  jwt: string;
  setUpdates: (updates: CustomerRefreshable) => void;
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
    const { email } = values;
    const response = await updateEmailAction(widgetId, jwt, {
      email,
    });
    const { error, data } = response;
    if (error) {
      const { message } = error;
      form.setError("root.serverError", {
        message: message || "Please try again later.",
      });
      return;
    }
    if (data) {
      try {
        const updates = customerSchema.parse(data);
        setUpdates({
          externalId: updates.externalId,
          email: updates.email,
          phone: updates.phone,
          name: updates.name,
          avatarUrl: updates.avatarUrl,
          isVerified: updates.isVerified,
          role: updates.role,
          requireIdentities: updates.requireIdentities,
          createdAt: updates.createdAt,
          updatedAt: updates.updatedAt,
        });
      } catch (err) {
        console.error("Failed to parse customer", err);
        if (err instanceof z.ZodError) {
          form.setError("root.serverError", {
            message: "Please try again later.",
          });
        }
      }
    }
    form.reset({ email: "" });
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col">
        <div className="flex w-full justify-between items-center">
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
              </FormItem>
            )}
          />
          <SubmitButton isDisabled={isSubmitting} />
        </div>
        {errors?.root?.serverError && (
          <FormMessage>{errors?.root?.serverError?.message}</FormMessage>
        )}
      </form>
    </Form>
  );
}
