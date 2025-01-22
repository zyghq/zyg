import { RecentCustomerEvents } from "@/components/thread/recent-customer-events.tsx";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import {
  SetThreadAssigneeForm,
  SetThreadPriorityForm,
  SetThreadStatusForm,
  ThreadLabels,
} from "@/components/workspace/thread/thread-properties-forms.tsx";
import { customerRoleVerboseName, getInitials } from "@/db/helpers.ts";
import { ThreadShape } from "@/db/shapes.ts";
import { WorkspaceStoreState } from "@/db/store.ts";
import { cn } from "@/lib/utils";
import { useWorkspaceStore } from "@/providers";
import { DotsHorizontalIcon } from "@radix-ui/react-icons";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { CopyIcon, PanelRight } from "lucide-react";
import React from "react";
import { useStore } from "zustand";

interface ThreadDetailSheetProps {
  activeThread: ThreadShape;
  token: string;
  workspaceId: string;
}

export function ThreadDetailsSheet({
  activeThread,
  token,
  workspaceId,
}: ThreadDetailSheetProps) {
  const [open, setOpen] = React.useState(false);

  const workspaceStore = useWorkspaceStore();
  const workspaceLabels = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) => state.viewLabels(state),
  );

  const customerName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewCustomerName(state, activeThread?.customerId || ""),
  );
  const customerEmail = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewCustomerEmail(state, activeThread?.customerId || ""),
  );
  const customerExternalId = useStore(
    workspaceStore,
    (state: WorkspaceStoreState) =>
      state.viewCustomerExternalId(state, activeThread?.customerId || ""),
  );
  const customerPhone = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewCustomerPhone(state, activeThread?.customerId || ""),
  );
  const customerRole = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewCustomerRole(state, activeThread?.customerId || ""),
  );

  const threadStage = activeThread?.stage || "";
  const assigneeId = activeThread?.assigneeId || "unassigned";
  const priority = activeThread?.priority || "normal";

  const [, copyEmail] = useCopyToClipboard();
  const [, copyExternalId] = useCopyToClipboard();
  const [, copyPhone] = useCopyToClipboard();

  return (
    <Sheet onOpenChange={setOpen} open={open}>
      <SheetTrigger asChild>
        <Button className="flex h-7 w-7 md:hidden" size="icon" variant="ghost">
          <PanelRight />
        </Button>
      </SheetTrigger>
      <SheetContent className="px-3 w-full">
        <SheetHeader className="hidden">
          <SheetTitle>Open Thread Details Menu</SheetTitle>
          <SheetDescription>Thread Details Menu</SheetDescription>
        </SheetHeader>
        <ScrollArea
          className={cn(
            "flex w-full flex-col",
            "h-[calc(100dvh-2rem)]"
          )}
        >
          {/* properties */}
          <div className="flex items-center justify-between border-b py-2">
            <span className="font-serif text-sm font-medium">Properties</span>
          </div>
          {/* forms */}
          <div className="flex flex-col gap-1 border-b py-2">
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <span className="font-serif text-sm font-medium">Status</span>
              </div>
              <div className="flex-1">
                <SetThreadStatusForm
                  stage={threadStage}
                  threadId={activeThread.threadId}
                  token={token}
                  workspaceId={workspaceId}
                />
              </div>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <span className="font-serif text-sm font-medium">Priority</span>
              </div>
              <div className="flex-1">
                <SetThreadPriorityForm
                  priority={priority}
                  threadId={activeThread.threadId}
                  token={token}
                  workspaceId={workspaceId}
                />
              </div>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex-1">
                <span className="font-serif text-sm font-medium">Assignee</span>
              </div>
              <div className="flex-1">
                <SetThreadAssigneeForm
                  assigneeId={assigneeId}
                  threadId={activeThread.threadId}
                  token={token}
                  workspaceId={workspaceId}
                />
              </div>
            </div>
            <ThreadLabels
              threadId={activeThread.threadId}
              token={token}
              workspaceId={workspaceId}
              workspaceLabels={workspaceLabels}
            />
          </div>
          {/* customer */}
          <div className="flex flex-col py-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <Avatar className="h-7 w-7">
                  <AvatarImage
                    src={`https://avatar.vercel.sh/${activeThread?.customerId || ""}`}
                  />
                  <AvatarFallback>{getInitials(customerName)}</AvatarFallback>
                </Avatar>
                <div className="flex flex-col">
                  <div className="text-sm font-semibold">{customerName}</div>
                  <div className="text-xs text-muted-foreground">
                    {customerRoleVerboseName(customerRole)}
                  </div>
                </div>
              </div>
              <Button size="icon" variant="ghost">
                <DotsHorizontalIcon className="h-4 w-4" />
              </Button>
            </div>
            <div className="flex items-center justify-between">
              <div className="text-xs font-medium">Email</div>
              <div className="flex items-center space-x-2">
                <div className="text-xs">{customerEmail || "n/a"}</div>
                <Button
                  className="text-muted-foreground"
                  onClick={() => copyEmail(customerEmail || "n/a")}
                  size="icon"
                  type="button"
                  variant="ghost"
                >
                  <CopyIcon className="h-3 w-3" />
                </Button>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <div className="text-xs font-medium">External ID</div>
              <div className="flex items-center space-x-2">
                <div className="text-xs">{customerExternalId || "n/a"}</div>
                <Button
                  className="text-muted-foreground"
                  onClick={() => copyExternalId(customerExternalId || "n/a")}
                  size="icon"
                  type="button"
                  variant="ghost"
                >
                  <CopyIcon className="h-3 w-3" />
                </Button>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <div className="text-xs font-medium">Phone</div>
              <div className="flex items-center space-x-2">
                <div className="text-xs">{customerPhone || "n/a"}</div>
                <Button
                  className="text-muted-foreground"
                  onClick={() => copyPhone(customerPhone || "n/a")}
                  size="icon"
                  type="button"
                  variant="ghost"
                >
                  <CopyIcon className="h-3 w-3" />
                </Button>
              </div>
            </div>
          </div>
          {/*  recent customer events */}
          <div className="flex flex-col">
            <RecentCustomerEvents
              customerId={activeThread.customerId}
              token={token}
              workspaceId={workspaceId}
            />
          </div>
        </ScrollArea>
      </SheetContent>
    </Sheet>
  );
}
