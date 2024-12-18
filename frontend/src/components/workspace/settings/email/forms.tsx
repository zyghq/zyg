import { Spinner } from "@/components/spinner";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  addEmailDomain,
  createEmailSetting,
  updateEmailSetting,
  verifyDNS,
} from "@/db/api";
import { PostmarkMailServerSetting } from "@/db/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { CopyIcon, Power } from "lucide-react";
import * as React from "react";
import { Control, useForm, useWatch } from "react-hook-form";
import { z } from "zod";

interface SupportEmailProps {
  editable: boolean;
  email: string;
  memberName: string;
  refetch: () => void;
  token: string;
  workspaceId: string;
  workspaceName: string;
}

const SupportEmailFormSchema = z.object({
  email: z.string().email(),
});

function EmailWatched({
  control,
  memberName,
  workspaceName,
}: {
  control: Control<z.infer<typeof SupportEmailFormSchema>>;
  memberName: string;
  workspaceName: string;
}) {
  const email = useWatch({
    control,
    defaultValue: "support@example.com",
    name: "email",
  });
  return (
    <p className="mb-4 text-xs text-muted-foreground">
      Outgoing emails will be sent using member's name and workspace name, e.g.{" "}
      <span className="font-semibold">
        {memberName || "User"} at {workspaceName}
      </span>{" "}
      {`(${email || "support@example.com"})`}
    </p>
  );
}

export function SupportEmailForm({
  editable,
  email,
  memberName,
  refetch = () => {},
  token,
  workspaceId,
  workspaceName,
}: SupportEmailProps) {
  const [, copyToClipboard] = useCopyToClipboard();
  const [isEditable, setIsEditable] = React.useState(editable);

  const form = useForm<z.infer<typeof SupportEmailFormSchema>>({
    defaultValues: {
      email: email, // Set initial value directly from prop
    },
    resolver: zodResolver(SupportEmailFormSchema),
  });
  const { control, reset } = form;

  // on props change
  React.useEffect(() => {
    reset({ email });
  }, [email, reset]);

  React.useEffect(() => {
    setIsEditable(editable);
  }, [editable]);

  const mutation = useMutation({
    mutationFn: async (inputs: { email: string }) => {
      const { email } = inputs;
      const { data, error } = await createEmailSetting(token, workspaceId, {
        email: email,
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
      form.setError("email", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    },
    onSuccess: (data) => {
      console.log("onSuccess");
      console.log(data);
      form.reset({ email: data.email });
      refetch();
    },
  });

  async function onSubmit(values: z.infer<typeof SupportEmailFormSchema>) {
    const { email } = values;
    await mutation.mutateAsync({ email });
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormControl className="mb-2">
                <div className="relative flex items-center">
                  <Input
                    {...field}
                    className={!isEditable ? "mr-1" : ""} // Add padding when copy icon is visible
                    disabled={!isEditable}
                    required
                    type="email"
                  />
                  {!isEditable && (
                    <Button
                      onClick={() => copyToClipboard(field.value)}
                      size="icon"
                      type="button"
                      variant="ghost"
                    >
                      <CopyIcon className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              </FormControl>
            </FormItem>
          )}
        />

        <EmailWatched
          control={control}
          memberName={memberName}
          workspaceName={workspaceName}
        />

        {/*<p className="mb-4 text-xs text-muted-foreground">*/}
        {/*  Outgoing emails will be sent using member's name and workspace name,*/}
        {/*  e.g.{" "}*/}
        {/*  <span className="font-semibold">*/}
        {/*    {memberName} at {workspaceName}*/}
        {/*  </span>{" "}*/}
        {/*  ({value || "support@example.com"})*/}
        {/*</p>*/}

        <div className="flex space-x-2">
          {isEditable ? (
            <>
              <Button type="submit">Save</Button>
              <Button
                onClick={() => setIsEditable(false)}
                type="button"
                variant="outline"
              >
                Cancel
              </Button>
            </>
          ) : (
            <Button
              onClick={() => setIsEditable(true)}
              type="button"
              variant={"outline"}
            >
              Edit
            </Button>
          )}
        </div>
      </form>
    </Form>
  );
}

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
  confirm: boolean;
  domain: string;
  refetch?: () => void;
  token: string;
  workspaceId: string;
}

const ConfirmDomainFormSchema = z.object({
  confirm: z.boolean().default(false),
  domain: z.string(),
});

export function ConfirmDomainForm({
  confirm,
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
    reset({ confirm, domain });
  }, [domain, reset, confirm]);

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
            name="confirm"
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
          <FormField
            control={form.control}
            name="domain"
            render={({ field }) => (
              <FormItem className="hidden">
                <FormControl>
                  <input
                    onChange={field.onChange}
                    type="hidden"
                    value={field.value}
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

interface EnableEmailFormProps {
  enabled: boolean;
  refetch?: () => void;
  token: string;
  workspaceId: string;
}

const EnableEmailFormSchema = z.object({
  enabled: z.boolean().default(false),
});

export function EnableEmailForm({
  enabled,
  refetch = () => {},
  token,
  workspaceId,
}: EnableEmailFormProps) {
  const form = useForm<z.infer<typeof EnableEmailFormSchema>>({
    defaultValues: {
      enabled: enabled, // Set initial value directly from prop
    },
    resolver: zodResolver(EnableEmailFormSchema),
  });
  const { reset } = form;

  // Update form when enabled prop changes
  React.useEffect(() => {
    reset({ enabled: enabled });
  }, [enabled, reset]);

  const mutation = useMutation({
    mutationFn: async (inputs: { enabled: boolean }) => {
      const { enabled } = inputs;
      const { data, error } = await updateEmailSetting(token, workspaceId, {
        enabled: enabled,
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
      form.setError("enabled", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    },
    onSuccess: (data) => {
      console.log("onSuccess");
      console.log(data);
      form.reset({ enabled: data.isEnabled });
      refetch();
    },
  });

  return (
    <>
      {enabled ? (
        <Button
          disabled={mutation.isPending}
          onClick={() => mutation.mutateAsync({ enabled: false })}
          type="submit"
          variant={"destructive"}
        >
          <Power className="mr-2 h-5 w-5" /> Disable Email
        </Button>
      ) : (
        <Button
          disabled={mutation.isPending}
          onClick={() => mutation.mutateAsync({ enabled: true })}
          type="submit"
        >
          <Power className="mr-2 h-5 w-5" /> Enable Email
        </Button>
      )}
    </>
  );
}

interface VerifyDNSFormProps {
  refetch?: () => void;
  token: string;
  workspaceId: string;
}

export function VerifyDNS({
  refetch = () => {},
  token,
  workspaceId,
}: VerifyDNSFormProps) {
  const mutation = useMutation({
    mutationFn: async () => {
      const { data, error } = await verifyDNS(token, workspaceId);
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
    },
    onSuccess: (data) => {
      console.log("onSuccess");
      console.log(data);
      refetch();
    },
  });

  return (
    <Button
      disabled={mutation.isPending}
      onClick={() => mutation.mutate()}
      variant="outline"
    >
      {mutation.isPending && <Spinner className="mr-1 h-4 w-4 animate-spin" />}
      Verify DNS
    </Button>
  );
}
