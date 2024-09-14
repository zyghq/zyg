import { cn } from "@/lib/utils";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Link, getRouteApi } from "@tanstack/react-router";
import { buttonVariants } from "@/components/ui/button";
import { PersonIcon, CircleIcon } from "@radix-ui/react-icons";
import { TagIcon } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { LabelMetrics } from "@/db/models";

const routeApi = getRouteApi("/_account/workspaces/$workspaceId/_workspace");

export function MyThreadsLink({
  workspaceId,
  memberId,
  assignedCount,
  openClose = () => {},
}: {
  workspaceId: string;
  memberId: string;
  assignedCount: number;
  openClose?: (isOpen: boolean) => void | undefined;
}) {
  const routeSearch = routeApi.useSearch();
  const { sort } = routeSearch;
  return (
    <Link
      to="/workspaces/$workspaceId/threads/me"
      onClick={() => openClose(false)}
      params={{ workspaceId }}
      search={{ sort }}
      className={cn(
        buttonVariants({ variant: "ghost" }),
        "flex w-full justify-between px-3 dark:text-accent-foreground"
      )}
      activeOptions={{ exact: true, includeSearch: false }}
      activeProps={{
        className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
      }}
    >
      {({ isActive }: { isActive: boolean }) => (
        <>
          {isActive ? (
            <>
              <div className="flex gap-x-2">
                <Avatar className="h-5 w-5">
                  <AvatarImage
                    src={`https://avatar.vercel.sh/${memberId}`}
                    alt={memberId}
                  />
                  <AvatarFallback>CN</AvatarFallback>
                </Avatar>
                <div className="font-semibold">Your Threads</div>
              </div>
              <span className="font-mono text-muted-foreground pr-2">
                {assignedCount}
              </span>
            </>
          ) : (
            <>
              <div className="flex gap-x-2">
                <Avatar className="h-5 w-5">
                  <AvatarImage
                    src={`https://avatar.vercel.sh/${memberId}`}
                    alt={memberId}
                  />
                  <AvatarFallback>M</AvatarFallback>
                </Avatar>
                <div className="font-normal">Your Threads</div>
              </div>
              <span className="font-mono text-muted-foreground pr-2">
                {assignedCount}
              </span>
            </>
          )}
        </>
      )}
    </Link>
  );
}

export function UnassignedThreadsLink({
  workspaceId,
  unassignedCount,
  openClose = () => {},
}: {
  workspaceId: string;
  unassignedCount: number;
  openClose?: (isOpen: boolean) => void | undefined;
}) {
  const routeSearch = routeApi.useSearch();
  const { sort } = routeSearch;
  return (
    <Link
      to="/workspaces/$workspaceId/threads/unassigned"
      onClick={() => openClose(false)}
      params={{ workspaceId }}
      search={{ sort }}
      className={cn(
        buttonVariants({ variant: "ghost" }),
        "flex w-full justify-between px-3 dark:text-accent-foreground"
      )}
      activeOptions={{ exact: true, includeSearch: false }}
      activeProps={{
        className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
      }}
    >
      {({ isActive }: { isActive: boolean }) => (
        <>
          {isActive ? (
            <>
              <div className="flex">
                <PersonIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                <div className="font-semibold">Unassigned Threads</div>
              </div>
              <span className="font-mono text-muted-foreground pr-2">
                {unassignedCount}
              </span>
            </>
          ) : (
            <>
              <div className="flex">
                <PersonIcon className="my-auto mr-2 h-4 w-4 text-muted-foreground" />
                <div className="font-normal">Unassigned Threads</div>
              </div>
              <span className="font-mono text-muted-foreground pr-2">
                {unassignedCount}
              </span>
            </>
          )}
        </>
      )}
    </Link>
  );
}

export function TodoThreadsLink({
  workspaceId,
  activeCount,
  openClose = () => {},
}: {
  workspaceId: string;
  activeCount: number;
  openClose?: (isOpen: boolean) => void | undefined;
}) {
  const routeSearch = routeApi.useSearch();
  const { sort } = routeSearch;
  return (
    <Link
      onClick={() => openClose(false)}
      to="/workspaces/$workspaceId/threads/todo"
      params={{ workspaceId }}
      search={{ sort }}
      className={cn(
        buttonVariants({ variant: "ghost" }),
        "flex w-full justify-between px-3 dark:text-accent-foreground"
      )}
      activeOptions={{ exact: true, includeSearch: false }}
      activeProps={{
        className: "bg-indigo-50 hover:bg-indigo-100 dark:bg-accent",
      }}
    >
      {({ isActive }: { isActive: boolean }) => (
        <>
          {isActive ? (
            <>
              <div className="flex gap-x-2">
                <CircleIcon className="my-auto h-4 w-4 text-indigo-500" />
                <div className="font-semibold">Todo</div>
              </div>
              <Badge className="bg-indigo-500 font-mono text-white hover:bg-indigo-600">
                {activeCount}
              </Badge>
            </>
          ) : (
            <>
              <div className="flex gap-x-2">
                <CircleIcon className="my-auto h-4 w-4 text-indigo-500" />
                <div className="font-normal">Todo</div>
              </div>
              <Badge variant="outline" className="font-mono font-light">
                {activeCount}
              </Badge>
            </>
          )}
        </>
      )}
    </Link>
  );
}

export function LabelThreadsLink({
  workspaceId,
  label,
}: {
  workspaceId: string;
  label: LabelMetrics;
}) {
  const routeSearch = routeApi.useSearch();
  const { sort } = routeSearch;
  return (
    <Link
      key={label.labelId}
      to="/workspaces/$workspaceId/threads/labels/$labelId"
      params={{ workspaceId, labelId: label.labelId }}
      search={{ sort: sort }}
      className={cn(
        buttonVariants({ variant: "ghost" }),
        "flex w-full justify-between px-3 dark:text-accent-foreground"
      )}
      activeOptions={{ exact: true }}
      activeProps={{
        className: "bg-indigo-100 hover:bg-indigo-200 dark:bg-accent",
      }}
    >
      {({ isActive }) => (
        <>
          {isActive ? (
            <>
              <div className="flex">
                <TagIcon className="my-auto mr-1 h-3 w-3 text-muted-foreground" />
                <div className="font-normal capitalize text-foreground">
                  {label.name}
                </div>
              </div>
              <div className="font-mono font-light text-muted-foreground">
                {label.count}
              </div>
            </>
          ) : (
            <>
              <div className="flex">
                <TagIcon className="my-auto mr-1 h-3 w-3 text-muted-foreground" />
                <div className="font-normal capitalize text-foreground">
                  {label.name}
                </div>
              </div>
              <div className="font-mono font-light text-muted-foreground">
                {label.count}
              </div>
            </>
          )}
        </>
      )}
    </Link>
  );
}
