import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import Link from "next/link";
import { EnvelopeClosedIcon } from "@radix-ui/react-icons";

export default function WorkspaceSetupPage({ params }) {
  const { workspaceId } = params;
  return (
    <div className="mt-24">
      <Card className="mx-auto w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-center text-3xl">
            Where do you talk to your customers?
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-center text-zinc-500">
            {`Link up your first favorite channel to Zyg. Don't worry, you can always add more channels later.`}
          </p>
        </CardContent>
        <CardFooter className="flex-col space-y-2">
          <div className="flex w-full rounded-xl px-2 py-4 outline outline-1 outline-zinc-200">
            <div className="my-auto flex flex-1 justify-start">
              <Image
                src="/slack.svg"
                alt="Slack Logo"
                width={48}
                height={48}
                priority
              />
              <div>
                <div className="my-auto font-semibold">Slack</div>
                <p className="text-xs text-zinc-500">
                  Connect Zyg to your Slack Workspace.
                </p>
              </div>
            </div>
            <div className="my-auto">
              <Button variant="outline">Connect</Button>
            </div>
          </div>
          <div className="flex w-full rounded-xl px-2 py-4 outline outline-1 outline-zinc-200">
            <div className="my-auto flex flex-1 justify-start">
              <EnvelopeClosedIcon
                width={24}
                height={24}
                className="mx-4 my-auto text-zinc-300"
              />
              <div>
                <div className="my-auto font-semibold">Email</div>
                <p className="text-xs text-zinc-500">
                  Setup email-forwarding - requires DNS updates.
                </p>
              </div>
            </div>
            <div className="my-auto">
              <Button variant="outline">Setup</Button>
            </div>
          </div>
        </CardFooter>
        <CardFooter className="flex justify-center">
          <Button variant="link" asChild>
            <Link href={`/workspaces/${workspaceId}/`}>skip for now</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
