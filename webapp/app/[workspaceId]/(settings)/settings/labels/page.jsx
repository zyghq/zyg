import { Separator } from "@/components/ui/separator";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { WorkspaceLabelAddOrEditForm } from "@/components/settings/forms";
import { LabelItem } from "@/components/settings/label";

export default function SettingsLabelsPage() {
  return (
    <div className="container md:mx-auto">
      <div className="max-w-2xl">
        <div className="pt-8 lg:pt-12">
          <div className="pb-8">
            <header className="text-xl font-semibold">Labels</header>
          </div>
          <Separator />
        </div>
        <div className="pt-8">
          <Tabs defaultValue="active">
            <div className="flex">
              <TabsList className="flex">
                <TabsTrigger value="active">
                  Active
                  <span className="ml-1 font-mono text-muted-foreground">
                    12
                  </span>
                </TabsTrigger>
                <TabsTrigger value="archived">
                  Archived
                  <span className="ml-1 font-mono text-muted-foreground">
                    0
                  </span>
                </TabsTrigger>
              </TabsList>
            </div>
            <TabsContent value="active">
              <div className="mt-8 flex flex-col gap-4">
                <WorkspaceLabelAddOrEditForm />
                <div className="flex flex-col gap-1">
                  <LabelItem label={"urgent"} />
                  <LabelItem label={"bug"} />
                </div>
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
    </div>
  );
}
