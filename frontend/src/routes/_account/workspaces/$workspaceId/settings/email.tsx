import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { DNSRecords } from "@/components/workspace/settings/dns-records";
import { EmailForwardEnableForm } from "@/components/workspace/settings/email-forward-enable-form";
import { EnableEmail } from "@/components/workspace/settings/enable-email";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/email",
)({
  component: ComingSoonSettings,
});

function ComingSoonSettings() {
  return (
    <div className="container py-12">
      <div className="max-w-2xl space-y-8">
        <header>
          <h1 className="mb-2 text-2xl font-semibold">Email</h1>
          <p className="text-sm text-muted-foreground">
            To setup emails for this workspace, you will need to be able to edit
            DNS records for your domain. Feel free to get in touch with us if
            you need help.
          </p>
        </header>

        <Separator />

        <section className="flex flex-col">
          <h2 className="mb-2 text-xl">
            1. Your customer service email address
          </h2>
          <div className="mb-4 text-sm text-muted-foreground">
            The email address where you want customers to contact your company.
            Outgoing emails to customers will also be sent from this email
            address.
          </div>
          <code className="mb-2 flex rounded-md bg-accent p-4">
            support@logly.dev
          </code>
          <p className="text-xs text-muted-foreground">
            Outgoing emails will be sent using member's name and workspace name,
            e.g. <span className="font-semibold">Sanchit at Zyg</span>{" "}
            (support@logly.dev)
          </p>
        </section>
        <div>
          <Button variant="outline">Edit</Button>
        </div>
        <Separator />
        <section className="flex flex-col">
          <h2 className="mb-2 text-xl">2. Receiving emails</h2>
          <div className="mb-4 text-sm text-muted-foreground">
            This allows Zyg to process incoming emails from customers.
          </div>
          <div className="mb-2 text-sm font-semibold">
            Forward inbound emails to this address:
          </div>
          <code className="mb-2 flex rounded-md bg-accent p-4">
            inbound@logly.dev
          </code>
          <EmailForwardEnableForm />
        </section>
        <Separator />
        <section className="flex flex-col">
          <h2 className="mb-2 text-xl">3. Sending emails</h2>
          <div className="mb-4 text-sm text-muted-foreground">
            Allows Zyg to send emails on your behalf. Verifying your domain
            gives email clients it was sent by Zyg with your permissions.
          </div>
          <div className="mb-2 text-sm font-semibold">
            Add the following DNS records for logly.dev.
          </div>
          <DNSRecords />
        </section>
        <Separator />
        <section className="flex flex-col">
          <h2 className="mb-2 text-xl">4. Enable email</h2>
          <div className="mb-4 text-sm text-muted-foreground">
            Now that you've completed the configuration process, you can enable
            this email address for use in the workspace.
          </div>
        </section>
        <EnableEmail />
      </div>
    </div>
  );
}
