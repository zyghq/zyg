import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { addEmailDomain, updateEmailSetting } from "@/db/api";
import { PostmarkMailServerSetting } from "@/db/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { CopyIcon } from "lucide-react";
import * as React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

interface EmailForwardFormProps {
  enabled: boolean;
  refetch?: () => void;
  token: string;
  workspaceId: string;
}

const EmailForwardFormSchema = z.object({
  hasForwardingEnabled: z.boolean().default(false),
});

export function EmailForwardForm({
  enabled,
  refetch = () => {},
  token,
  workspaceId,

}: EmailForwardFormProps) {
  const form = useForm<z.infer<typeof EmailForwardFormSchema>>({
    defaultValues: {
      hasForwardingEnabled: enabled, // Set initial value directly from prop
    },
    resolver: zodResolver(EmailForwardFormSchema),
  });
  const { reset } = form;

  // Update form when enabled prop changes
  React.useEffect(() => {
    reset({ hasForwardingEnabled: enabled });
  }, [enabled, reset]);

  const mutation = useMutation({
    mutationFn: async (inputs: { hasForwardingEnabled: boolean }) => {
      const { hasForwardingEnabled } = inputs;
      const { data, error } = await updateEmailSetting(token, workspaceId, {
        hasForwardingEnabled: hasForwardingEnabled,
      });
      if (error) {
        throw new Error(error.message);
      }
      if (!data) {
        throw new Error("no data returned");
      }
      return data as PostmarkMailServerSetting;
    },
    onError: (error) => {
      console.error(error);
      form.setError("hasForwardingEnabled", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    },
    onSuccess: (data) => {
      console.log("onSuccess");
      console.log(data);
      form.reset({ hasForwardingEnabled: data.hasForwardingEnabled });
      refetch();
    },
  });

  async function onSubmit(values: z.infer<typeof EmailForwardFormSchema>) {
    const { hasForwardingEnabled } = values;
    await mutation.mutateAsync({ hasForwardingEnabled });
  }

  return (
    <Form {...form}>
      <form className="space-y-6" onSubmit={form.handleSubmit(onSubmit)}>
        <FormField
          control={form.control}
          name="hasForwardingEnabled"
          render={({ field }) => (
            <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
              <FormControl>
                <Checkbox
                  checked={field.value}
                  onCheckedChange={field.onChange}
                />
              </FormControl>
              <FormLabel>Email forwarding is set up</FormLabel>
            </FormItem>
          )}
        />
        <Button disabled={mutation.isPending} type="submit" variant="outline">
          Save
        </Button>
      </form>
    </Form>
  );
}

interface ConfirmDomainFormProps {
  domain: string;
  refetch?: () => void;
  token: string;
  workspaceId: string;
}

const ConfirmDomainFormSchema = z.object({
  domain: z.string(),
});

export function ConfirmDomainForm({
  domain,
  refetch = () => {},
  token,
  workspaceId,
}: ConfirmDomainFormProps) {
  const [, copyToClipboard] = useCopyToClipboard();

  const form = useForm<z.infer<typeof ConfirmDomainFormSchema>>({
    defaultValues: {
      domain: domain, // Set initial value directly from prop
    },
    resolver: zodResolver(ConfirmDomainFormSchema),
  });
  const { reset } = form;

  // Update form when enabled prop changes
  React.useEffect(() => {
    reset({ domain });
  }, [domain, reset]);

  const mutation = useMutation({
    mutationFn: async (inputs: { domain: string }) => {
      const { domain } = inputs;
      const { data, error } = await addEmailDomain(token, workspaceId, {
        domain: domain,
      });
      if (error) {
        throw new Error(error.message);
      }
      if (!data) {
        throw new Error("no data returned");
      }
      return data as PostmarkMailServerSetting;
    },
    onError: (error) => {
      console.error(error);
      form.setError("domain", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    },
    onSuccess: (data) => {
      console.log("onSuccess");
      console.log(data);
      form.reset({ domain: data.domain });
      refetch();
    },
  });

  async function onSubmit(values: z.infer<typeof ConfirmDomainFormSchema>) {
    const { domain } = values;
    await mutation.mutateAsync({ domain });
  }

  return (
    <>
      <div className="mb-2 flex items-center">
        <code className="mr-1 w-full rounded-md bg-accent p-2">{domain}</code>
        <Button
          onClick={() => copyToClipboard(domain)}
          size="icon"
          type="button"
          variant="ghost"
        >
          <CopyIcon className="h-4 w-4" />
        </Button>
      </div>
      <Form {...form}>
        <form className="space-y-6" onSubmit={form.handleSubmit(onSubmit)}>
          <FormField
            control={form.control}
            name="domain"
            render={({ field }) => (
              <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                <FormControl>
                  <Checkbox
                    checked={Boolean(field.value)}
                    onCheckedChange={field.onChange}
                  />
                </FormControl>
                <FormLabel>Confirm domain for sending emails</FormLabel>
              </FormItem>
            )}
          />
          <Button disabled={mutation.isPending} type="submit" variant="outline">
            Save
          </Button>
        </form>
      </Form>
    </>
  );
}
