import { Icons } from "@/components/icons.tsx";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Dns } from "@/components/workspace/settings/email/dns.tsx";
import {
  ConfirmDomainForm,
  EmailForwardForm,
  EnableEmailForm,
  SupportEmailForm,
  VerifyDNS,
} from "@/components/workspace/settings/email/forms.tsx";
import { getEmailSetting } from "@/db/api";
import { PostmarkMailServerSetting } from "@/db/models";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { CopyIcon } from "lucide-react";
import { motion } from "motion/react";
import { useStore } from "zustand";

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

interface ReceivingEmailProps {
  hasForwardingEnabled: boolean;
  inboundEmail: string;
  refetch: () => void;
  token: string;
  workspaceId: string;
}

const ReceivingEmail = ({
  hasForwardingEnabled,
  inboundEmail,
  refetch = () => {},
  token,
  workspaceId,
}: ReceivingEmailProps) => {
  const [, copyToClipboard] = useCopyToClipboard();
  return (
    <motion.section className="flex flex-col" variants={sectionVariants}>
      <h2 className="mb-2 text-xl">2. Receiving emails</h2>
      <div className="mb-4 text-sm text-muted-foreground">
        This allows Zyg to process incoming emails from customers.
      </div>
      <div className="mb-2 text-sm font-semibold">
        Forward inbound emails to this address:
      </div>
      <div className="mb-2 flex items-center">
        <code className="mr-1 w-full rounded-md bg-accent p-2">
          {inboundEmail}
        </code>
        <Button
          onClick={() => copyToClipboard(inboundEmail)}
          size="icon"
          type="button"
          variant="ghost"
        >
          <CopyIcon className="h-4 w-4" />
        </Button>
      </div>
      <EmailForwardForm
        enabled={hasForwardingEnabled}
        refetch={refetch}
        token={token}
        workspaceId={workspaceId}
      />
    </motion.section>
  );
};

function EmailSettings() {
  const { token } = Route.useRouteContext();
  const { workspaceId } = Route.useParams();

  const workspaceStore = useWorkspaceStore();
  const workspaceName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceName(state),
  );
  const memberName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getMemberName(state),
  );

  const {
    data: setting,
    error,
    isPending,
    refetch: refetchSetting,
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
      hostname: setting?.returnPathDomain,
      status: setting?.returnPathDomainVerified ? "Verified" : "Pending",
      type: "CNAME",
      value: setting?.returnPathDomainCNAME,
    } as DNSRecord;
    records.push(dkimRecord);
    records.push(cnameRecord);
    return records;
  };

  console.log(setting);

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

        <motion.section className="flex flex-col" variants={sectionVariants}>
          <h2 className="mb-2 text-xl">
            1. Your customer service email address
          </h2>
          <div className="mb-4 text-sm text-muted-foreground">
            The email address where you want customers to contact your company.
            Outgoing emails to customers will also be sent from this email
            address.
          </div>
          <SupportEmailForm
            editable={!setting?.email}
            email={setting?.email || ""}
            memberName={memberName}
            refetch={refetchSetting}
            token={token}
            workspaceId={workspaceId}
            workspaceName={workspaceName}
          />
        </motion.section>

        <Separator />

        {setting?.inboundEmail && (
          <ReceivingEmail
            hasForwardingEnabled={setting?.hasForwardingEnabled}
            inboundEmail={setting?.inboundEmail}
            refetch={refetchSetting}
            token={token}
            workspaceId={workspaceId}
          />
        )}
        {setting?.inboundEmail && <Separator />}

        {setting?.hasForwardingEnabled && (
          <motion.section className="flex flex-col" variants={sectionVariants}>
            <h2 className="mb-2 text-xl">3. Sending emails</h2>
            <div className="mb-4 text-sm text-muted-foreground">
              Allows Zyg to send emails on your behalf. Verifying your domain
              gives email clients it was sent by Zyg with your permissions.
            </div>
            {setting?.hasDNS ? (
              <motion.div
                className="flex flex-col space-y-4"
                variants={sectionVariants}
              >
                <div className="mb-2 text-sm font-semibold">
                  Add the following DNS records for logly.dev.
                </div>
                <Dns records={dnsRecords()} />
                <div>
                  <VerifyDNS token={token} workspaceId={workspaceId} />
                </div>
                <div className={"text-sm text-muted-foreground"}>
                  It may take a few minutes for your DNS changes to propagate
                  throughout the internet. In some cases, this process can take
                  up to 24 hours.
                </div>
              </motion.div>
            ) : (
              <ConfirmDomainForm
                confirm={setting?.hasDNS}
                domain={setting?.domain}
                refetch={refetchSetting}
                token={token}
                workspaceId={workspaceId}
              />
            )}
          </motion.section>
        )}
        {setting?.hasForwardingEnabled && <Separator />}

        {setting?.hasForwardingEnabled && (
          <motion.section className="flex flex-col" variants={sectionVariants}>
            <h2 className="mb-2 text-xl">4. Enable email</h2>
            <div className="mb-4 text-sm text-muted-foreground">
              Now that you've completed the configuration process, you can
              enable this email address for use in the workspace.
            </div>
            <div className="flex items-center gap-4">
              <EnableEmailForm
                enabled={setting?.isEnabled || false}
                refetch={refetchSetting}
                token={token}
                workspaceId={workspaceId}
              />
            </div>
          </motion.section>
        )}
      </div>
    </motion.div>
  );
}
