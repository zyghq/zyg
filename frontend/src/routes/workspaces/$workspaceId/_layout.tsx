import { z } from "zod";
import { createFileRoute, Outlet } from "@tanstack/react-router";

import { useStore } from "zustand";
import { WorkspaceStoreStateType } from "@/db/store";
import { Header } from "@/components/workspace/header";
import { SideNav } from "@/components/workspace/sidenav";

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

const threadSearchSchema = z.object({
  status: z.enum(["todo", "snoozed", "done"]).catch("todo"),
  reasons: reasonsSchema(["replied", "unreplied"]).catch(""),
  sort: z
    .enum(["last-message-dsc", "created-asc", "created-dsc"])
    .catch("last-message-dsc"),
  priorities: prioritiesSchema(["urgent", "high", "normal", "low"]).catch(""),
});

export const Route = createFileRoute("/workspaces/$workspaceId/_layout")({
  validateSearch: (search) => threadSearchSchema.parse(search),
  component: () => <WorkspaceLayout />,
});

function WorkspaceLayout() {
  const { WorkspaceStore, AccountStore } = Route.useRouteContext();

  const email = useStore(AccountStore.useContext(), (state) =>
    state.getEmail(state)
  );

  const workspaceId = useStore(
    WorkspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getWorkspaceId(state)
  );
  const workspaceName = useStore(
    WorkspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getWorkspaceName(state)
  );

  const memberId = useStore(
    WorkspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getMemberId(state)
  );

  const metrics = useStore(
    WorkspaceStore.useContext(),
    (state: WorkspaceStoreStateType) => state.getMetrics(state)
  );

  return (
    <div vaul-drawer-wrapper="">
      <div className="flex flex-col">
        <Header
          workspaceId={workspaceId}
          workspaceName={workspaceName}
          metrics={metrics}
          memberId={memberId}
        />
        <div className="flex flex-col">
          <div className="grid lg:grid-cols-5">
            <SideNav
              email={email}
              workspaceId={workspaceId}
              workspaceName={workspaceName}
              metrics={metrics}
              memberId={memberId}
            />
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  );
}
