import { createFileRoute } from "@tanstack/react-router";
import { LabelItem } from "@/components/workspace/settings/label";
import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { z } from "zod";
import { useForm, SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useStore } from "zustand";
import { useWorkspaceStore } from "@/providers";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export const Route = createFileRoute(
  "/_auth/workspaces/$workspaceId/settings/labels"
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
  const workspaceStore = useWorkspaceStore();
  const labels = useStore(workspaceStore, (state) => state.viewLabels(state));

  const form = useForm({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
    },
  });

  const onSubmit: SubmitHandler<FormInputs> = async (inputs) => {
    // Do something with the form values.
    // ✅ This will be type-safe and validated.
    console.log(inputs);
  };

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
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="flex items-center space-x-2"
                >
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem className="w-full">
                        <FormControl>
                          <Input placeholder="Label" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <Button type="submit">Add</Button>
                </form>
              </Form>
              {labels && labels.length > 0 ? (
                <div className="flex flex-col gap-1">
                  {labels.map((label) => (
                    <LabelItem key={label.labelId} label={label.name} />
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
