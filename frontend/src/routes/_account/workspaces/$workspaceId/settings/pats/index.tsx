import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { deletePat } from "@/db/api";
import { useWorkspaceStore } from "@/providers";
import { PlusIcon } from "@radix-ui/react-icons";
import { createFileRoute, Link } from "@tanstack/react-router";
import { format } from "date-fns";
import { KeyRoundIcon } from "lucide-react";
import * as React from "react";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/pats/"
)({
  component: PatSettings,
});

const formatDate = (date: string) => {
  const dateObj = new Date(date);
  return format(dateObj, "MMMM d, yyyy");
};

function DeleteConfirmationDialog({
  deleteFunc,
  isDeleted,
  isDeleting,
  isError,
  resetFunc,
}: {
  deleteFunc: () => void;
  isDeleted: boolean;
  isDeleting: boolean;
  isError: boolean;
  resetFunc: () => void;
}) {
  const [open, setOpen] = React.useState(false);

  React.useEffect(() => {
    if (isDeleted) {
      setOpen(false);
    }
  }, [isDeleted]);

  function onDelete(
    e: React.BaseSyntheticEvent<
      MouseEvent,
      EventTarget & HTMLButtonElement,
      EventTarget
    >
  ) {
    e.preventDefault();
    deleteFunc();
  }

  function onCancel() {
    setOpen(false);
    setTimeout(() => {
      resetFunc();
    }, 100);
  }

  return (
    <AlertDialog onOpenChange={setOpen} open={open}>
      <AlertDialogTrigger asChild>
        <Button size="sm" variant="destructive">
          Delete
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>
            Are you sure you want to delete this token?
          </AlertDialogTitle>
          <AlertDialogDescription>
            This action <span className="font-bold">CANNOT</span> be undone.
            Proceeding will permanently delete your token, affecting any usage
            associated with it.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={() => onCancel()}>
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction disabled={isDeleting} onClick={onDelete}>
            I understand, delete it
          </AlertDialogAction>
        </AlertDialogFooter>
        {isError && (
          <div className="text-xs text-red-500">
            Something went wrong. Please try again later.
          </div>
        )}
      </AlertDialogContent>
    </AlertDialog>
  );
}

function PatSettings() {
  const { token } = Route.useRouteContext();
  const { workspaceId } = Route.useParams();
  const workspaceStore = useWorkspaceStore();

  const pats = useStore(workspaceStore, (state) => state.viewPats(state));

  const [isDeleting, setIsDeleting] = React.useState(false);
  const [isDeleted, setIsDeleted] = React.useState(false);
  const [isError, setIsError] = React.useState(false);

  function onDelete(token: string, patId: string) {
    return async () => {
      setIsDeleting(true);
      const { error } = await deletePat(token, patId);
      if (error) {
        setIsError(true);
      } else {
        workspaceStore.getState().deletePat(patId);
        setIsDeleted(true);
      }
      setIsDeleting(false);
    };
  }

  function reset() {
    setIsDeleting(false);
    setIsDeleted(false);
    setIsError(false);
  }

  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="flex items-center justify-between my-12">
            <header className="text-xl font-semibold">
              Personal Access Tokens
            </header>
            <Button asChild size="sm">
              <Link
                params={{ workspaceId }}
                to="/workspaces/$workspaceId/settings/pats/add"
              >
                <PlusIcon className="mr-1 h-4 w-4" />
                Add Token
              </Link>
            </Button>
          </div>
          <Separator />
        </div>
        <div className="mt-8 flex flex-col gap-2">
          {pats && pats.length > 0 ? (
            <React.Fragment>
              {pats.map((pat) => (
                <div
                  className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left"
                  key={pat.patId}
                >
                  <div className="flex w-full flex-col gap-1">
                    <div className="flex items-center">
                      <div className="flex items-center gap-2">
                        <KeyRoundIcon className="h-5 w-5" />
                        <div className="flex flex-col">
                          <div className="font-normal">{pat.name}</div>
                          <div className="text-xs">{pat.token}</div>
                        </div>
                      </div>
                      <div className="ml-auto">
                        <DeleteConfirmationDialog
                          deleteFunc={onDelete(token, pat.patId)}
                          isDeleted={isDeleted}
                          isDeleting={isDeleting}
                          isError={isError}
                          resetFunc={reset}
                        />
                      </div>
                    </div>
                    <div className="flex flex-col">
                      <div className="text-sm text-muted-foreground">
                        {pat.description}
                      </div>
                      <div className="text-xs">
                        {`Added on ${formatDate(pat.createdAt)}`}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </React.Fragment>
          ) : (
            <div className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left">
              <div className="flex w-full flex-col gap-1">
                <div className="text-md">
                  No personal access tokens created.
                </div>
                <div className="text-sm text-muted-foreground">
                  {`When you create a personal access token, it will appear here.`}
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
