import { getSession, isAuthenticated } from "@/utils/supabase/helpers";
import { createClient } from "@/utils/supabase/server";
import * as React from "react";

import { redirect } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import ThreadList from "@/components/dashboard/thread-list";

import { DoubleArrowUpIcon, MixerHorizontalIcon } from "@radix-ui/react-icons";

import { CheckCircle, CircleIcon, EclipseIcon } from "lucide-react";

export const metadata = {
  title: "Unassigned Threads - Zyg AI",
};

/**
 * Fetches the list of thread chats for a given workspace and assigned member.
 *
 * @param {string} workspaceId - The ID of the workspace.
 * @param {string} [authToken=""] - The authentication token (optional).
 * @returns {Promise<{ data: any, error: Error | null }>} - The response object containing the data and error (if any).
 */
async function getMyThreadChatListAPI(workspaceId, authToken = "") {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}/threads/chat/with/unassigned/`,
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

export default async function AssignedToMePage({ params }) {
  const { workspaceId } = params;

  const supabase = createClient();
  if (!(await isAuthenticated(supabase))) {
    return redirect("/login/");
  }

  const { token, error: tokenErr } = await getSession(supabase);
  if (tokenErr) {
    return redirect("/login/");
  }

  const threads = [];
  const { error, data } = await getMyThreadChatListAPI(workspaceId, token);

  if (error) {
    return (
      <main className="col-span-3 lg:col-span-4">
        <div className="container mt-12">
          <h1 className="mb-1 text-3xl font-bold">Error</h1>
          <p className="mb-4 text-red-500">
            There was an error fetching your threads. Please try again later.
          </p>
        </div>
      </main>
    );
  } else {
    threads.push(...data);
  }

  return (
    <React.Fragment>
      <main className="col-span-3 lg:col-span-4">
        <div className="container">
          <div className="mb-4 mt-4 text-xl">Unassigned Threads</div>
          <Tabs defaultValue="todo">
            <div className="mb-4 sm:flex sm:justify-between">
              <TabsList className="grid grid-cols-3">
                <TabsTrigger value="todo">
                  <div className="flex items-center">
                    <CircleIcon className="mr-1 h-4 w-4 text-indigo-500" />
                    Todo
                  </div>
                </TabsTrigger>
                <TabsTrigger value="snoozed">
                  <div className="flex items-center">
                    <EclipseIcon className="mr-1 h-4 w-4 text-fuchsia-500" />
                    Snoozed
                  </div>
                </TabsTrigger>
                <TabsTrigger value="done">
                  <div className="flex items-center">
                    <CheckCircle className="mr-1 h-4 w-4 text-green-500" />
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
                endpoint="/threads/chat/with/unassigned/"
              />
            </TabsContent>
            <TabsContent value="snoozed" className="m-0">
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
    </React.Fragment>
  );
}
