import { getSession } from "@/utils/supabase/helpers";
import { createClient } from "@/utils/supabase/server";

import { Header } from "@/components/dashboard/header";
import { Sidebar } from "@/components/dashboard/sidebar";

export const metadata = {
  title: "All Threads - Zyg AI",
};

async function getThreadMetricsAPI(workspaceId, authToken = "") {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}/threads/chat/metrics/`,
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
          `Error fetching thread metrics: ${status} ${statusText}`
        ),
      };
    }

    const data = await response.json();
    return { data, error: null };
  } catch (err) {
    console.error("error fetching workspace thread metrics", err);
    return { data: null, error: err };
  }
}

async function getWorkspaceAPI(workspaceId, authToken = "") {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_ZYG_URL}/workspaces/${workspaceId}/`,
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
          `Error fetching workspace details: ${status} ${statusText}`
        ),
      };
    }

    const data = await response.json();
    return { data, error: null };
  } catch (err) {
    console.error("error fetching workspace details", err);
    return { data: null, error: err };
  }
}

export default async function DashboardLayout({ params, children }) {
  const { workspaceId } = params;

  const supabase = createClient();
  const { token, error: tokenErr } = await getSession(supabase);

  if (tokenErr) {
    return redirect("/login/");
  }

  const workspaceResp = getWorkspaceAPI(workspaceId, token);
  const metricsResp = getThreadMetricsAPI(workspaceId, token);

  const [workspaceData, metricsData] = await Promise.all([
    workspaceResp,
    metricsResp,
  ]);

  const { error: workspaceError, data: workspace } = workspaceData;
  const { error: metricsError, data: metrics } = metricsData;

  if (workspaceError || metricsError) {
    return (
      <div className="flex h-screen flex-col">
        <div className="mx-auto my-auto">
          <h1 className="mb-1 text-3xl font-bold">Error</h1>
          <p className="mb-4 text-red-500">
            There was an error. Please try again later.
          </p>
        </div>
      </div>
    );
  }

  const { name: workspaceName } = workspace;

  return (
    <div vaul-drawer-wrapper="">
      <div className="flex flex-col">
        <Header
          workspaceId={workspaceId}
          workspaceName={workspaceName}
          metrics={metrics}
        />
        <div className="flex flex-col">
          <div className="grid lg:grid-cols-5">
            <Sidebar
              workspaceId={workspaceId}
              workspaceName={workspaceName}
              metrics={metrics}
            />
            {children}
          </div>
        </div>
      </div>
    </div>
  );
}
