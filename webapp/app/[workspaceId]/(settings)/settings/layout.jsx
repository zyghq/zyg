import { Button } from "@/components/ui/button";
import { ArrowLeftIcon } from "@radix-ui/react-icons";
import { SettingsHeader } from "@/components/settings-header";
import SettingsMenuScroll from "@/components/settings-menu-scroll";

export const metadata = {
  title: "Settings - Zyg AI",
};

export default function SettingsLayout({ children }) {
  return (
    <div vaul-drawer-wrapper="">
      <div className="flex flex-col">
        <SettingsHeader />
        <div className="flex flex-col">
          <div className="flex">
            <div className="hidden min-w-80 flex-col border-r lg:flex">
              <SettingsMenuScroll className="h-[calc(100dvh-8rem)]" />
              <div className="sticky bottom-0 flex h-14 border-t">
                <div className="flex w-full items-center">
                  <div className="ml-4">
                    <Button size="sm" variant="outline">
                      <ArrowLeftIcon className="mr-2 h-4 w-4" />
                      Support
                    </Button>
                  </div>
                </div>
              </div>
            </div>
            <main className="flex-1">{children}</main>
          </div>
        </div>
      </div>
    </div>
  );
}
