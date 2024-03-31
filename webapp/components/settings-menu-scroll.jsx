import { cn } from "@/lib/utils";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ChatBubbleIcon, CodeIcon } from "@radix-ui/react-icons";
import { BlocksIcon } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Icons } from "@/components/icons";

export default function SettingsMenuScroll({ className }) {
  return (
    <ScrollArea className={cn("pl-2 pr-1", className)}>
      {/* TODO: refactor this when you can items can groupped inside. */}
      <div className="px-2 pb-8">
        {/* G1 */}
        <div className="my-4 flex items-center gap-1">
          <Avatar className="h-8 w-8">
            <AvatarImage src="https://github.com/shadcn.png" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
          <div>
            <div className="text-xs font-medium">Sanchit Rk</div>
            <div className="text-xs text-foreground">Account</div>
          </div>
        </div>
        {/* G1 Items */}
        <div className="flex flex-col gap-1">
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">Appearance</div>
              </div>
            </Link>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">Personal Notifications</div>
              </div>
            </Link>
          </Button>
        </div>
        {/* G2 */}
        <div className="my-4 flex items-center gap-1">
          <Icons.logo className="mx-1 h-5 w-5" />
          <div>
            <div className="text-xs font-medium">ZygHQ</div>
            <div className="text-xs text-foreground">Workspace</div>
          </div>
        </div>
        {/* G2 Items */}
        <div className="flex flex-col gap-1">
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">General</div>
              </div>
            </Link>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">Members</div>
              </div>
            </Link>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">AI</div>
              </div>
            </Link>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">Labels</div>
              </div>
            </Link>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">Billing</div>
              </div>
            </Link>
          </Button>
        </div>
        {/* G3 */}
        <div className="my-4 flex items-center gap-1">
          <ChatBubbleIcon className="mx-1 h-4 w-4" />
          <div className="text-xs font-medium text-foreground">Channels</div>
        </div>
        {/* G3 Items */}
        <div className="flex flex-col gap-1">
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">Slack</div>
              </div>
            </Link>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">Email</div>
              </div>
            </Link>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <Link href={`/`}>
              <div className="flex">
                <div className="my-auto">Chat</div>
              </div>
            </Link>
          </Button>
        </div>
        {/* G4 */}
        <div className="my-4 flex items-center gap-1">
          <BlocksIcon className="mx-1 h-4 w-4" />
          <div className="text-xs font-medium text-foreground">
            Integrations
          </div>
        </div>
        {/* G4 Items */}
        <div className="flex flex-col gap-1">
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <div className="flex">
              <div className="my-auto">GitHub</div>
            </div>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <div className="flex">
              <div className="my-auto">Linear</div>
            </div>
          </Button>
        </div>
        {/* G5 */}
        <div className="my-4 flex items-center gap-1">
          <CodeIcon className="mx-1 h-4 w-4" />
          <div className="text-xs font-medium text-foreground">Build</div>
        </div>
        {/* G5 Items */}
        <div className="flex flex-col gap-1">
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <div className="flex">
              <div className="my-auto">Webhooks</div>
            </div>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <div className="flex">
              <div className="my-auto">Events</div>
            </div>
          </Button>
          <Button
            variant="ghost"
            asChild
            className="flex w-full justify-between"
          >
            <div className="flex">
              <div className="my-auto">SDKs</div>
            </div>
          </Button>
        </div>
      </div>
    </ScrollArea>
  );
}
