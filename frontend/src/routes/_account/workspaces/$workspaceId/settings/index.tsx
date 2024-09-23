import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { updateWorkspace } from "@/db/api";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { zodResolver } from "@hookform/resolvers/zod";
import { CheckIcon, ExclamationTriangleIcon } from "@radix-ui/react-icons";
import { CheckCircledIcon, ClipboardCopyIcon } from "@radix-ui/react-icons";
import { createFileRoute, useRouterState } from "@tanstack/react-router";
import { useCopyToClipboard } from "@uidotdev/usehooks";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";
import { useStore } from "zustand";

const formSchema = z.object({
  name: z.string().min(3),
});

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/"
)({
  component: GeneralSettings,
});

type FormInputs = {
  name: string;
};

function GeneralSettings() {
  const { workspaceId } = Route.useParams();
  const { token } = Route.useRouteContext();
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const workspaceStore = useWorkspaceStore();

  const workspaceName = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.getWorkspaceName(state)
  );

  const [copiedText, copyToClipboard] = useCopyToClipboard();
  const hasCopiedText = Boolean(copiedText);

  const form = useForm({
    defaultValues: {
      name: workspaceName,
    },
    resolver: zodResolver(formSchema),
  });

  const { formState } = form;

  const { errors, isSubmitSuccessful, isSubmitting } = formState;

  const isUpdating = isLoading || isSubmitting;

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    try {
      const { name } = inputs;
      const { data, error } = await updateWorkspace(token, workspaceId, {
        name,
      });
      if (error) {
        console.error(error);
        form.setError("root", {
          message: "Something went wrong. Please try again later.",
          type: "serverError",
        });
        return;
      }
      if (data) {
        const { workspaceName } = data;
        workspaceStore.getState().updateWorkspaceName(workspaceName);
      }
    } catch (err) {
      console.error(err);
      form.setError("root", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    }
  };

  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="my-12">
            <header className="text-xl font-semibold">Settings</header>
          </div>
          <Separator />
        </div>
        <Form {...form}>
          <form className="space-y-8" onSubmit={form.handleSubmit(onSubmit)}>
            {errors?.root?.type === "serverError" && (
              <Alert variant="destructive">
                <ExclamationTriangleIcon className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>
                  {`${errors?.root?.message || "Please try again later."}`}
                </AlertDescription>
              </Alert>
            )}
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Zyg" {...field} />
                  </FormControl>
                  <FormDescription>
                    Typically your company or team name.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <div>
              <div className="text-sm font-normal">Workspace ID</div>
              <div className="text-xs text-muted-foreground">
                We might ask you for support inquiries.
              </div>
              <div className="mt-4 flex w-full items-center">
                <code className="mr-2 rounded-lg border bg-muted p-2">
                  {workspaceId}
                </code>
                <Button
                  aria-disabled={hasCopiedText}
                  disabled={hasCopiedText}
                  onClick={() => copyToClipboard(workspaceId)}
                  size="sm"
                  type="button"
                  variant="ghost"
                >
                  {hasCopiedText ? (
                    <CheckCircledIcon className="h-5 w-5 text-green-500" />
                  ) : (
                    <ClipboardCopyIcon className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
            <div className="flex space-x-4">
              <Button
                aria-disabled={isUpdating}
                aria-label="Create Workspace"
                disabled={isUpdating}
                type="submit"
              >
                Save
              </Button>
              {isSubmitSuccessful && (
                <div className="flex my-auto text-green-500">
                  <CheckIcon className="my-auto h-4 w-4" />
                  <div className="text-sm my-auto">Workspace Updated!</div>
                </div>
              )}
            </div>
          </form>
        </Form>
      </div>
    </div>
  );
}
