"use client";

import { createClient } from "@/utils/supabase/client";
import { getSession } from "@/utils/supabase/helpers";
import { useQuery } from "@tanstack/react-query";
import * as React from "react";

export default function DashboardContainer({ workspaceId }) {
  const supabase = createClient();
  const threadsData = useQuery({
    queryKey: ["threads", workspaceId, supabase],
    queryFn: async () => {
      const { token, error: sessErr } = await getSession(supabase);
      if (sessErr) throw new Error("session expired or not found");
      const url = `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}/threads/chat/`;
      const response = await fetch(url, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        const { status, statusText } = response;
        throw new Error(`error fetching threads: ${status} ${statusText}`);
      }
      return await response.json();
    },
  });

  if (threadsData.isPending) console.log("fetching threads");
  if (threadsData.error)
    console.log("error fetching threads", threadsData.error);

  //   if (result.isError) console.log("error fetching threads", result.error);
  //   if (result.isPending) console.log("fetching threads");

  console.log("result", threadsData.data);

  return (
    <React.Fragment>
      <div>header</div>
      <div className="flex flex-col">
        <div className="grid lg:grid-cols-5">
          <div>sidebar</div>
          <div>content</div>
        </div>
      </div>
    </React.Fragment>
  );
}
