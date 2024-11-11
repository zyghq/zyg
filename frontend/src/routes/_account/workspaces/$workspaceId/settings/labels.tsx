import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { LabelItem } from "@/components/workspace/settings/label";
import { createWorkspaceLabel } from "@/db/api";
import { Label } from "@/db/models";
import { WorkspaceStoreState } from "@/db/store";
import { useWorkspaceStore } from "@/providers";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { SubmitHandler, useForm } from "react-hook-form";
import { z } from "zod";
import { useStore } from "zustand";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/labels",
)({
  component: LabelSettings,
});

type FormInputs = {
  name: string;
};

const formSchema = z.object({
  name: z.string().min(2),
});

function LabelSettings() {
  const { workspaceId } = Route.useParams();
  const { token } = Route.useRouteContext();
  const workspaceStore = useWorkspaceStore();
  const labels = useStore(workspaceStore, (state: WorkspaceStoreState) =>
    state.viewLabels(state),
  ) as Label[];

  const form = useForm({
    defaultValues: {
      name: "",
    },
    resolver: zodResolver(formSchema),
  });

  const mutation = useMutation({
    mutationFn: async (inputs: { name: string }) => {
      const { name } = inputs;
      const { data, error } = await createWorkspaceLabel(token, workspaceId, {
        name,
      });
      if (error) {
        throw new Error(error.message);
      }
      if (!data) {
        throw new Error("no data returned");
      }
      return data;
    },
    onError: (error) => {
      console.error(error);
      form.setError("name", {
        message: "Something went wrong. Please try again later.",
        type: "serverError",
      });
    },
    onSuccess: (data) => {
      workspaceStore.getState().addLabel(data);
      form.reset();
    },
  });

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    await mutation.mutateAsync(inputs);
  };

  const { formState } = form;
  const { isSubmitting } = formState;

  const isCreating = isSubmitting || mutation.isPending;

  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="my-12">
            <header className="text-xl font-semibold">Labels</header>
          </div>
          <Separator />
        </div>
        <Tabs defaultValue="active">
          <div className="flex">
            <TabsList className="flex">
              <TabsTrigger value="active">
                Active
                <span className="ml-1 text-muted-foreground">
                  {labels.length}
                </span>
              </TabsTrigger>
              <TabsTrigger value="archived">
                Archived
                <span className="ml-1 text-muted-foreground">0</span>
              </TabsTrigger>
            </TabsList>
          </div>
          <TabsContent value="active">
            <div className="mt-8 flex flex-col gap-4">
              <Form {...form}>
                <form
                  className="flex space-x-2"
                  onSubmit={form.handleSubmit(onSubmit)}
                >
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem className="w-full">
                        <FormControl>
                          <Input
                            autoComplete="off"
                            placeholder="Label"
                            required
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <Button disabled={isCreating} type="submit">
                    Add
                  </Button>
                </form>
              </Form>
              {labels && labels.length > 0 ? (
                <div className="flex flex-col gap-1">
                  {labels.map((label) => (
                    <LabelItem
                      key={label.labelId}
                      label={label.name}
                      labelId={label.labelId}
                      token={token}
                      workspaceId={workspaceId}
                    />
                  ))}
                </div>
              ) : (
                <div className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left">
                  <div className="flex w-full flex-col gap-1">
                    <div className="text-md">There are no labels.</div>
                    <div className="text-sm text-muted-foreground">
                      {`You can add labels to your workspace to categorize your threads.`}
                    </div>
                  </div>
                </div>
              )}
            </div>
          </TabsContent>
          <TabsContent value="archived">
            <div className="mt-8 flex flex-col">
              <div className="flex flex-col items-start gap-2 rounded-lg border p-3 text-left">
                <div className="flex w-full flex-col gap-1">
                  <div className="text-md">No archived labels.</div>
                  <div className="text-sm text-muted-foreground">
                    {`Instead of deleting labels from your workspace, you can archive them. Archived labels will still be applied to current threads but cannot be added to new threads.`}
                  </div>
                </div>
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
