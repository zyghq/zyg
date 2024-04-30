import Link from "next/link";
import { cn } from "@/lib/utils";
import { CommandMenu } from "@/components/commander";
import { MobileNav } from "@/components/mobile-nav";
import { ModeToggle } from "@/components/theme";
import { buttonVariants } from "@/components/ui/button";
import { ArrowLeftRightIcon } from "lucide-react";
import { Icons } from "@/components/icons";

// TODO: rename to DashboardHeader or MainHeader
export function Header() {
  return (
    <header className="sticky top-0 z-50 flex h-14 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-4 flex w-full items-center">
        <div className="hidden md:flex">
          <Link href="/" className="flex items-center space-x-2">
            <Icons.logo className="h-5 w-5" />
            <span className="hidden font-semibold sm:inline-block">Zyg.</span>
          </Link>
        </div>
        <MobileNav />
        <div className="flex flex-1 items-center justify-between space-x-2 md:justify-end">
          <div className="w-full flex-1 md:w-auto md:flex-none">
            <CommandMenu />
          </div>
          <nav className="flex items-center">
            <Link href="/workspaces/">
              <div
                className={cn(
                  buttonVariants({
                    variant: "ghost",
                  }),
                  "w-9 px-0",
                )}
              >
                <ArrowLeftRightIcon className="h-4 w-4" />
                <span className="sr-only">Switch Workspace</span>
              </div>
            </Link>
            <ModeToggle />
          </nav>
        </div>
      </div>
    </header>
  );
}
