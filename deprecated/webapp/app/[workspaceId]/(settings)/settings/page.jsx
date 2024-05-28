import Link from "next/link";
import { Separator } from "@/components/ui/separator";
import { WorkspaceEditForm } from "@/components/settings/forms";

export default function SettingsGeneralPage({ params }) {
  const { workspaceId } = params;
  return (
    <div className="container md:mx-auto">
      <div className="max-w-2xl">
        <div className="pt-8 lg:pt-12">
          <div className="pb-8">
            <header className="text-xl font-semibold">Settings</header>
          </div>
          <Separator />
        </div>
        <div className="pt-8">
          <WorkspaceEditForm workspaceId={workspaceId} />
        </div>
      </div>
    </div>
  );
}
