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
  const { slug } = params;
  return (
    <div className="mt-24">
      <Card className="mx-auto w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-3xl text-center">
            Where do you talk to your customers?
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-zinc-500 text-center">
            {`Link up your first favorite channel to Zyg. Don't worry, you can always add more channels later.`}
          </p>
        </CardContent>
        <CardFooter className="space-y-2 flex-col">
          <div className="px-2 py-4 outline outline-1 outline-zinc-200 w-full rounded-xl flex">
            <div className="my-auto flex-1 flex justify-start">
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
          <div className="px-2 py-4 outline outline-1 outline-zinc-200 w-full rounded-xl flex">
            <div className="my-auto flex-1 flex justify-start">
              <EnvelopeClosedIcon
                width={24}
                height={24}
                className="text-zinc-300 my-auto mx-4"
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
            <Link href={`/workspaces/${slug}/`}>skip for now</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
