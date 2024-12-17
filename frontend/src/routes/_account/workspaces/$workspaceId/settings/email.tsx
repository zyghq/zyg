import { Icons } from "@/components/icons.tsx";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { DNSRecords } from "@/components/workspace/settings/dns-records";
import { EmailForwardForm } from "@/components/workspace/settings/email-forward-form.tsx";
import { EnableEmail } from "@/components/workspace/settings/enable-email";
import { getEmailSetting } from "@/db/api.ts";
import { PostmarkMailServerSetting } from "@/db/schema.ts";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { CopyIcon } from "lucide-react";
import { motion } from "motion/react";
import * as React from "react";
import { useForm } from "react-hook-form";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/email",
)({
  component: EmailSettings,
});

interface DNSRecord {
  hostname: string;
  status: "Pending" | "Verified";
  type: string;
  value: string;
}

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      when: "beforeChildren",
    },
  },
};

const sectionVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    transition: {
      duration: 0.1,
      ease: "easeOut",
    },
    y: 0,
  },
};

interface CustomerServiceEmailProps {
  editable: boolean;
  email: string;
  memberName: string;
  workspaceName: string;
}

interface SupportEmailFormValues {
  email: string;
}

const CustomerServiceEmail = ({
  editable,
  email,
  memberName,
  workspaceName,
}: CustomerServiceEmailProps) => {
  const [propEmail, setPropEmail] = React.useState(email);
  const [, copyToClipboard] = useCopyToClipboard();
  const [isEditable, setIsEditable] = React.useState(editable);

  const onSubmit = (data: SupportEmailFormValues) => {
    console.log(data);
  };

  const form = useForm<SupportEmailFormValues>({
    defaultValues: {
      email: propEmail,
    },
  });

  const { reset } = form;
  const value = form.getValues("email");

  React.useEffect(() => {
    setPropEmail(email);
    setIsEditable(editable);
    reset({ email });
  }, [email, editable, reset]);

  return (
    <motion.section className="flex flex-col" variants={sectionVariants}>
      <h2 className="mb-2 text-xl">1. Your customer service email address</h2>
      <div className="mb-4 text-sm text-muted-foreground">
        The email address where you want customers to contact your company.
        Outgoing emails to customers will also be sent from this email address.
      </div>

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

          <p className="mb-4 text-xs text-muted-foreground">
            Outgoing emails will be sent using member's name and workspace name,
            e.g.{" "}
            <span className="font-semibold">
              {memberName} at {workspaceName}
            </span>{" "}
            ({value || "support@example.com"})
          </p>

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
    </motion.section>
  );
};

interface ReceivingEmailProps {
  inboundEmail: string;
}

const ReceivingEmail = ({ inboundEmail }: ReceivingEmailProps) => {
  return (
    <motion.section className="flex flex-col" variants={sectionVariants}>
      <h2 className="mb-2 text-xl">2. Receiving emails</h2>
      <div className="mb-4 text-sm text-muted-foreground">
        This allows Zyg to process incoming emails from customers.
      </div>
      <div className="mb-2 text-sm font-semibold">
        Forward inbound emails to this address:
      </div>
      <motion.code
        className="mb-2 flex rounded-md bg-accent p-4"
        whileHover={{ scale: 1.01 }}
      >
        {inboundEmail}
      </motion.code>
      <EmailForwardForm />
    </motion.section>
  );
};

interface SendingEmailsProps {
  records: DNSRecord[];
}

const SendingEmails = ({ records }: SendingEmailsProps) => (
  <motion.section className="flex flex-col" variants={sectionVariants}>
    <h2 className="mb-2 text-xl">3. Sending emails</h2>
    <div className="mb-4 text-sm text-muted-foreground">
      Allows Zyg to send emails on your behalf. Verifying your domain gives
      email clients it was sent by Zyg with your permissions.
    </div>
    <div className="mb-2 text-sm font-semibold">
      Add the following DNS records for logly.dev.
    </div>
    <DNSRecords records={records} />
  </motion.section>
);

function EmailSettings() {
  const { token } = Route.useRouteContext();
  const { workspaceId } = Route.useParams();
  const {
    data: setting,
    error,
    isPending,
  } = useQuery({
    enabled: !!token,
    queryFn: async () => {
      const { data, error } = await getEmailSetting(token, workspaceId);
      if (error) throw new Error("failed to fetch email setting");
      return data as PostmarkMailServerSetting;
    },
    queryKey: ["emailSetting", workspaceId, token],
  });

  if (error) {
    return (
      <div className="container py-12">
        <div className="flex items-center gap-1">
          <Icons.oops className="h-5 w-5" />
          <div className="text-md text-red-500">Something went wrong</div>
        </div>
      </div>
    );
  }
  if (isPending) return null;

  const dnsRecords = (): DNSRecord[] => {
    const records: DNSRecord[] = [];
    const dkimRecord = {
      hostname: setting?.dkimHost,
      status: setting?.dkimUpdateStatus || "Pending",
      type: "TXT",
      value: setting?.dkimTextValue,
    } as DNSRecord;
    const cnameRecord = {
      hostname: setting?.returnPathDomainCNAME,
      status: setting?.returnPathDomainVerified ? "Verified" : "Pending",
      type: "CNAME",
      value: setting?.returnPathDomain,
    } as DNSRecord;
    records.push(dkimRecord);
    records.push(cnameRecord);
    return records;
  };

  return (
    <motion.div
      animate="visible"
      className="container py-12"
      initial="hidden"
      variants={containerVariants}
    >
      <div className="max-w-2xl space-y-8">
        <motion.header variants={sectionVariants}>
          <h1 className="mb-2 text-2xl font-semibold">Email</h1>
          <p className="text-sm text-muted-foreground">
            To setup emails for this workspace, you will need to be able to edit
            DNS records for your domain. Feel free to get in touch with us if
            you need help.
          </p>
        </motion.header>
        <Separator />
        <CustomerServiceEmail
          editable={!setting?.email}
          email={setting?.email || ""}
          memberName={"Sanchit"}
          workspaceName={"Zyg"}
        />
        <Separator />

        {setting?.inboundEmail && (
          <ReceivingEmail inboundEmail={setting?.inboundEmail} />
        )}

        {setting?.inboundEmail && <Separator />}

        {setting?.hasDNS && <SendingEmails records={dnsRecords()} />}

        {setting?.hasDNS && <Separator />}

        {setting?.hasForwardingEnabled && (
          <motion.section className="flex flex-col" variants={sectionVariants}>
            <h2 className="mb-2 text-xl">4. Enable email</h2>
            <div className="mb-4 text-sm text-muted-foreground">
              Now that you've completed the configuration process, you can
              enable this email address for use in the workspace.
            </div>
            <div className="flex items-center gap-4">
              <EnableEmail />
              <div className="text-sm font-semibold">Email is Disabled</div>
            </div>
          </motion.section>
        )}
      </div>
    </motion.div>
  );
}
