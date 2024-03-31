import Link from "next/link";
import { cn } from "@/lib/utils";
import { CommandMenu } from "@/components/commander";
import { Nav } from "@/components/nav";
import { MobileNav } from "@/components/mobile-nav";
import { ModeToggle } from "@/components/theme";
import { buttonVariants } from "@/components/ui/button";
import { ArrowLeftRightIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ArrowLeftIcon, ChatBubbleIcon, CodeIcon } from "@radix-ui/react-icons";
import { SettingsMobileNav } from "@/components/mobile-nav";

export function SettingsHeader() {
  return (
    <header className="sticky top-0 z-50 flex h-14 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-4 flex w-full items-center">
        <div className="hidden md:flex">
          <Button size="sm" variant="outline">
            <ArrowLeftIcon className="mr-2 h-4 w-4" />
            Settings
          </Button>
        </div>
        <SettingsMobileNav />
      </div>
    </header>
  );
}
