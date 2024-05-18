import { Button } from "@/components/ui/button";

import { Icons } from "@/components/icons";

import { ExitIcon, QuestionMarkIcon } from "@radix-ui/react-icons";

export const metadata = {
  title: "Select Or Create Workspace - Zyg AI",
};

export default function CreateOrSelectWorkspaceLayout({ children }) {
  return (
    <div className="relative flex min-h-screen flex-col bg-background">
      <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 max-w-screen-2xl items-center">
          <div className="flex flex-1 items-end">
            <Icons.logo className="my-auto mr-2 h-5 w-5" />
            <span className="font-semibold">Zyg.</span>
          </div>
          <div className="flex justify-between space-x-2 md:justify-end">
            <Button variant="outline" size="default">
              <QuestionMarkIcon />
              Help
            </Button>
            <Button variant="outline" size="icon">
              <ExitIcon />
            </Button>
          </div>
        </div>
      </header>
      <main className="relative py-6">{children}</main>
    </div>
  );
}
