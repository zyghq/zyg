import * as React from "react";
import { z } from "zod";
import { createFileRoute, Outlet } from "@tanstack/react-router";

import { useStore } from "zustand";
import { WorkspaceStoreStateType } from "@/db/store";
import { Header } from "@/components/workspace/header";
import SideNavLinks from "@/components/workspace/sidenav-links";
import { useAccountStore, useWorkspaceStore } from "@/providers";

//
// for more: https://tanstack.com/router/latest/docs/framework/react/guide/search-params
// usage of `.catch` or `default` matters.
const reasonsSchema = (validValues: string[]) => {
  const sanitizeArray = (arr: string[]) => {
    // remove duplicates
    const uniqueValues = [...new Set(arr)];
    // filter only valid values
    const uniqueValidValues: string[] = uniqueValues.filter((val) =>
      validValues.includes(val)
    );

    // no valid values
    if (uniqueValidValues.length === 0) {
      throw new Error("invalid reason(s) passed");
    }

    if (uniqueValidValues.length === 1) {
      return uniqueValidValues[0];
    }

    return uniqueValidValues;
  };
  return z.union([
    z.string().refine((value) => validValues.includes(value)),
    z.array(z.string()).transform(sanitizeArray),
    // .refine(
    //   (arr) =>
    //     arr.length === validValues.length &&
    //     validValues.every((val) => arr.includes(val))
    // ),
    z.undefined(),
  ]);
};

const prioritiesSchema = (validValues: string[]) => {
  const sanitizeArray = (arr: string[]) => {
    // remove duplicates
    const uniqueValues = [...new Set(arr)];
    // filter only valid values
    const uniqueValidValues: string[] = uniqueValues.filter((val) =>
      validValues.includes(val)
    );

    // no valid values
    if (uniqueValidValues.length === 0) {
      throw new Error("invalid prioritie(s) passed");
    }

    if (uniqueValidValues.length === 1) {
      return uniqueValidValues[0];
    }

    return uniqueValidValues;
  };
  return z.union([
    z.string().refine((value) => validValues.includes(value)),
    z.array(z.string()).transform(sanitizeArray),
    z.undefined(),
  ]);
};

const assigneesScheme = z.union([
  z.string(),
  z.array(z.string()),
  z.undefined(),
]);

const threadSearchSchema = z.object({
  status: z.enum(["todo", "snoozed", "done"]).catch("todo"),
  reasons: reasonsSchema(["replied", "unreplied"]).catch(""),
  sort: z
    .enum(["last-message-dsc", "created-asc", "created-dsc"])
    .catch("last-message-dsc"),
  priorities: prioritiesSchema(["urgent", "high", "normal", "low"]).catch(""),
  assignees: assigneesScheme.catch(""),
});

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/_workspace"
)({
  validateSearch: (search) => threadSearchSchema.parse(search),
  component: WorkspaceLayout,
});

function WorkspaceLayout() {
  const accountStore = useAccountStore();
  const workspaceStore = useWorkspaceStore();

  const email = useStore(accountStore, (state) => state.getEmail(state));

  const workspaceId = useStore(
    workspaceStore,
    (state: WorkspaceStoreStateType) => state.getWorkspaceId(state)
  );
  const workspaceName = useStore(
    workspaceStore,
    (state: WorkspaceStoreStateType) => state.getWorkspaceName(state)
  );

  const memberId = useStore(workspaceStore, (state: WorkspaceStoreStateType) =>
    state.getMemberId(state)
  );

  const metrics = useStore(workspaceStore, (state: WorkspaceStoreStateType) =>
    state.getMetrics(state)
  );

  return (
    <React.Fragment>
      <Header
        email={email}
        workspaceId={workspaceId}
        workspaceName={workspaceName}
        metrics={metrics}
        memberId={memberId}
      />
      <div className="flex min-h-screen">
        <aside className="hidden sticky top-14 h-[calc(100vh-theme(spacing.14))] w-80 overflow-y-auto md:block md:border-r">
          <SideNavLinks
            maxHeight="h-[calc(100dvh-8rem)]"
            email={email}
            workspaceId={workspaceId}
            workspaceName={workspaceName}
            metrics={metrics}
            memberId={memberId}
          />
        </aside>
        <main className="flex-1 mt-14 pb-4">
          <Outlet />
        </main>
      </div>
    </React.Fragment>
  );
}
