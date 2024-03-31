import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { PlusIcon, DotsHorizontalIcon } from "@radix-ui/react-icons";

export default function SettingsMembersPage() {
  return (
    <div className="container md:mx-auto">
      <div className="max-w-2xl">
        <div className="pt-8 lg:pt-12">
          <div className="flex items-center justify-between pb-8">
            <header className="text-xl font-semibold">Members</header>
            <Button size="sm" className="h-7">
              <PlusIcon className="mr-1 h-4 w-4" />
              Invite
            </Button>
          </div>
          <Separator />
        </div>
        <div className="pt-8">
          <Tabs defaultValue="members">
            <div className="flex w-1/2">
              <TabsList className="flex">
                <TabsTrigger value="members">Active Members</TabsTrigger>
                <TabsTrigger value="invites">Pending Invites</TabsTrigger>
              </TabsList>
            </div>
            <TabsContent value="members">
              <div className="mt-8 flex flex-col gap-2">
                <div className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left">
                  <div className="flex w-full flex-col gap-1">
                    <div className="flex items-center">
                      <div className="flex items-center gap-2">
                        <Avatar className="h-8 w-8">
                          <AvatarImage src="https://github.com/shadcn.png" />
                          <AvatarFallback>CN</AvatarFallback>
                        </Avatar>
                        <div className="">
                          <div className="font-normal"> Sanchit Rk</div>
                          <div className="text-xs text-muted-foreground">
                            sanchitrrk@gmail.com
                          </div>
                        </div>
                      </div>
                      <div className="ml-auto mr-2 text-sm text-muted-foreground">
                        Primary Owner
                      </div>
                      <Button variant="ghost" size="sm">
                        <DotsHorizontalIcon className="h-4 w-4" />
                      </Button>
                    </div>
                    <div>
                      <div className="text-xs text-muted-foreground">
                        Member since 12th Jan 2024
                      </div>
                    </div>
                  </div>
                </div>
                <div className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left">
                  <div className="flex w-full flex-col gap-1">
                    <div className="flex items-center">
                      <div className="flex items-center gap-2">
                        <Avatar className="h-8 w-8">
                          <AvatarImage src="https://github.com/shadcn.png" />
                          <AvatarFallback>CN</AvatarFallback>
                        </Avatar>
                        <div className="">
                          <div className="font-normal"> Sanchit Rk</div>
                          <div className="text-xs text-muted-foreground">
                            sanchitrrk@gmail.com
                          </div>
                        </div>
                      </div>
                      <div className="ml-auto mr-2 text-sm text-muted-foreground">
                        Primary Owner
                      </div>
                      <Button variant="ghost" size="sm">
                        <DotsHorizontalIcon className="h-4 w-4" />
                      </Button>
                    </div>
                    <div>
                      <div className="text-xs text-muted-foreground">
                        Member since 12th Jan 2024
                      </div>
                    </div>
                  </div>
                </div>
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
    </div>
  );
}
