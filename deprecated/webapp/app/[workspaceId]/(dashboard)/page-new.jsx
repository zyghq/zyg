import * as React from "react";

import DashboardContainer from "@/components/dashboard/container";

export const metadata = {
  title: "All Threads - Zyg AI",
};

export default function DashboardPage({ params }) {
  const { workspaceId } = params;
  return (
    <div>
      <DashboardContainer workspaceId={workspaceId} />
    </div>
  );
}
