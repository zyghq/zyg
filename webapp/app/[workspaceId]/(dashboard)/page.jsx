import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { createClient } from "@/utils/supabase/server";
import { isAuthenticated, getAuthToken } from "@/utils/supabase/helpers";
import { Sidebar } from "@/components/sidebar";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { DoubleArrowUpIcon, MixerHorizontalIcon } from "@radix-ui/react-icons";
import { EclipseIcon, CircleIcon, CheckCircle } from "lucide-react";
import ThreadList from "@/components/thread-list";

async function getWorkspaceThreadChats(workspaceId, authToken) {
  try {
    const response = await fetch(
      `${process.env.ZYG_API_URL}/workspaces/${workspaceId}/threads/chat/`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
      },
    );

    if (!response.ok) {
      const { status, statusText } = response;
      throw new Error(`Error fetching thread chats: ${status} ${statusText}`);
    }

    const data = await response.json();
    return [null, data];
  } catch (err) {
    console.error(err);
    return [err, null];
  }
}

export default async function DashboardPage({ params }) {
  const { workspaceId } = params;
  const cookieStore = cookies();
  const supabase = createClient(cookieStore);

  if (!(await isAuthenticated(supabase))) {
    return redirect("/login/");
  }

  const authToken = await getAuthToken(supabase);
  const [err, threadChats] = await getWorkspaceThreadChats(
    workspaceId,
    authToken,
  );

  if (err) {
    console.log("render error component");
  }

  return (
    <div className="grid lg:grid-cols-5">
      <Sidebar className="hidden border-r lg:block" />
      <main className="col-span-3 lg:col-span-4">
        <div className="container">
          <div className="mb-4 mt-5 text-xl font-semibold">Threads</div>
          <Tabs defaultValue="todo">
            <div className="mb-4 sm:flex sm:justify-between">
              <TabsList className="grid grid-cols-3">
                <TabsTrigger value="todo">
                  <div className="flex items-center">
                    <CircleIcon className="mr-1 h-3 w-3 text-indigo-500" />
                    Todo
                  </div>
                </TabsTrigger>
                <TabsTrigger value="inprogress">
                  <div className="flex items-center">
                    <EclipseIcon className="mr-1 h-3 w-3 text-fuchsia-500" />
                    In Progress
                  </div>
                </TabsTrigger>
                <TabsTrigger value="done">
                  <div className="flex items-center">
                    <CheckCircle className="mr-1 h-3 w-3 text-green-500" />
                    Done
                  </div>
                </TabsTrigger>
              </TabsList>
              <div className="mt-4 sm:my-auto">
                <Button variant="ghost" size="sm">
                  <MixerHorizontalIcon className="mr-1 h-3 w-3" />
                  Filters
                </Button>
                <Button variant="ghost" size="sm">
                  <DoubleArrowUpIcon className="mr-1 h-3 w-3" />
                  Sort
                </Button>
              </div>
            </div>
            <TabsContent value="todo" className="m-0">
              <ThreadList
                workspaceId={workspaceId}
                items={threadChats}
                className="h-[calc(100dvh-14rem)]"
              />
            </TabsContent>
            <TabsContent value="inprogress" className="m-0">
              {/* <ThreadList
          items={threads.filter((item) => !item.read)}
          className="h-[calc(100dvh-14rem)]"
        /> */}
            </TabsContent>
            <TabsContent value="done" className="m-0">
              {/* <Threads items={threads} /> */}
              {/* <ScrollArea className="h-[calc(100vh-14rem)] pr-1">
          <div className="flex flex-col gap-2">...</div>
        </ScrollArea> */}
            </TabsContent>
          </Tabs>
        </div>
      </main>
    </div>
  );
}
