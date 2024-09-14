import { createFileRoute } from "@tanstack/react-router";

import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PlusIcon, DotsHorizontalIcon } from "@radix-ui/react-icons";
import { useWorkspaceStore } from "@/providers";
import { useStore } from "zustand";
import React from "react";
import Avatar from "boring-avatars";
import { format } from "date-fns";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/members"
)({
  component: MemberSettings,
});

// date-fns format date string to 12 Jan 2024
const formatDate = (date: string) => {
  const dateObj = new Date(date);
  return format(dateObj, "MMMM d, yyyy");
};

function MemberSettings() {
  const workspaceStore = useWorkspaceStore();
  const members = useStore(workspaceStore, (state) => state.viewMembers(state));
  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="flex items-center justify-between my-12">
            <header className="text-xl font-semibold">Members</header>
            <Button size="sm" className="h-7">
              <PlusIcon className="mr-1 h-4 w-4" />
              Invite
            </Button>
          </div>
          <Separator />
        </div>
        <Tabs defaultValue="members">
          <div className="flex w-1/2">
            <TabsList className="flex">
              <TabsTrigger value="members">Active Members</TabsTrigger>
              <TabsTrigger value="invites">Pending Invites</TabsTrigger>
            </TabsList>
          </div>
          <TabsContent value="members">
            <div className="mt-8 flex flex-col gap-2">
              {members && members.length > 0 ? (
                <React.Fragment>
                  {members.map((member) => (
                    <div
                      key={member.memberId}
                      className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left"
                    >
                      <div className="flex w-full flex-col gap-1">
                        <div className="flex items-center">
                          <div className="flex items-center gap-2">
                            <Avatar
                              name={member.memberId}
                              size={32}
                              variant="marble"
                            />
                            <div className="flex flex-col">
                              <div className="font-normal">{member.name}</div>
                              <div className="text-xs text-muted-foreground"></div>
                            </div>
                          </div>
                          <div className="ml-auto mr-2 text-sm text-muted-foreground capitalize">
                            {member.role}
                          </div>
                          <Button variant="ghost" size="sm">
                            <DotsHorizontalIcon className="h-4 w-4" />
                          </Button>
                        </div>
                        <div>
                          <div className="text-xs text-muted-foreground">
                            {`Member since ${formatDate(member.createdAt)}`}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </React.Fragment>
              ) : (
                <div className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left">
                  <div className="flex w-full flex-col gap-1">
                    <div className="text-md">No members.</div>
                    <div className="text-sm text-muted-foreground">
                      {`When someone joins this workspace, they'll appear here.`}
                    </div>
                  </div>
                </div>
              )}
            </div>
          </TabsContent>
          <TabsContent value="invites">
            <div className="mt-8 flex flex-col">
              <div className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left">
                <div className="flex w-full flex-col gap-1">
                  <div className="text-md">No pending invites.</div>
                  <div className="text-sm text-muted-foreground">
                    {`When you invite someone to this workspace, they'll appear
                    here until they accept your invite. You can also cancel
                    pending invites.`}
                  </div>
                </div>
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
