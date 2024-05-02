import { getAuthToken, isAuthenticated } from "@/utils/supabase/helpers";
import { createClient } from "@/utils/supabase/server";

import { redirect } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { Sidebar } from "@/components/dashboard/sidebar";
import ThreadList from "@/components/dashboard/thread-list";

import { DoubleArrowUpIcon, MixerHorizontalIcon } from "@radix-ui/react-icons";

import { CheckCircle, CircleIcon, EclipseIcon } from "lucide-react";

/**
 * Fetches the list of thread chats for a given workspace.
 *
 * @param {string} workspaceId - The ID of the workspace.
 * @param {string} [authToken=""] - The authentication token (optional).
 * @returns {Promise<{ data: any, error: Error | null }>} - The response object containing the data and error (if any).
 */
async function getThreadChatListAPI(workspaceId, authToken = "") {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}/threads/chat/`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${authToken}`,
        },
      }
    );

    if (!response.ok) {
      const { status, statusText } = response;
      return {
        error: new Error(
          `Error fetching thread chats: ${status} ${statusText}`
        ),
      };
    }

    const data = await response.json();
    return { data, error: null };
  } catch (err) {
    console.error("error fetching workspace thread chats", err);
    return { data: null, error: err };
  }
}

/**
 *
 * TODO:
 * fetch last 10 threads based on status: todo, snoozed, closed
 * load them in parallel into the tabs as initial threads
 * based on query key status set the default view.
 *
 *
 */
export default async function DashboardPage({ params }) {
  const { workspaceId } = params;

  const supabase = createClient();

  if (!(await isAuthenticated(supabase))) {
    return redirect("/login/");
  }

  const authToken = await getAuthToken(supabase);

  const threads = [];
  const { error, data } = await getThreadChatListAPI(workspaceId, authToken);

  if (error) {
    return (
      <div className="container mt-12">
        <h1 className="mb-1 text-3xl font-bold">Error</h1>
        <p className="mb-4 text-red-500">
          There was an error fetching your threads. Please try again later.
        </p>
      </div>
    );
  } else {
    threads.push(...data);
  }

  return (
    <div className="grid lg:grid-cols-5">
      <Sidebar className="hidden lg:block lg:border-r" />
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
                threads={threads}
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
