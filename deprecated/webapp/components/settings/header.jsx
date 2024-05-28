import { Button } from "@/components/ui/button";
import { SidebarMobile } from "@/components/settings/sidebar-mobile";
import { ArrowLeftIcon } from "@radix-ui/react-icons";

export function Header() {
  return (
    <header className="sticky top-0 z-50 flex h-14 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-4 flex w-full items-center">
        <div className="hidden md:flex">
          <Button size="sm" variant="outline">
            <ArrowLeftIcon className="mr-2 h-4 w-4" />
            Settings
          </Button>
        </div>
        <SidebarMobile />
      </div>
    </header>
  );
}
