import { QueryFilter } from "@/lib/filters";
import { getSession } from "@/utils/supabase/helpers";
import { createClient } from "@/utils/supabase/server";
import * as React from "react";

import { redirect } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import ThreadList from "@/components/dashboard/thread-list";
import ThreadListContainer from "@/components/dashboard/thread-list-container";

import { DoubleArrowUpIcon, MixerHorizontalIcon } from "@radix-ui/react-icons";

import { CheckCircle, CircleIcon, EclipseIcon } from "lucide-react";

export const metadata = {
  title: "My Threads - Zyg AI",
};

async function getMyThreadChatListAPI(url, authToken = "") {
  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${authToken}`,
      },
    });

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

export default async function AssignedToMePage({ params, searchParams }) {
  const { workspaceId } = params;
  const filters = new QueryFilter(searchParams);
  const queryParams = filters.buildQuery();
  const cleanedQueryParams = filters.buildCleanedQuery();

  if (filters.redirect) {
    return redirect(`/${workspaceId}/threads/me/?${queryParams.toString()}`);
  }

  const supabase = createClient();
  const { token, error: tokenErr } = await getSession(supabase);
  if (tokenErr) {
    return redirect("/login/");
  }

  const threads = [];

  const url = `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}/threads/chat/with/me/?${cleanedQueryParams.toString()}`;

  const { error, data } = await getMyThreadChatListAPI(url, token);

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
    <ThreadListContainer
      name="My Threads"
      workspaceId={workspaceId}
      threads={threads}
      status={filters.getStatus()}
      url={url}
    />
  );
}
