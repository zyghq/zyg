import Image from "next/image";
import { Button } from "@/components/ui/button";
import { QuestionMarkIcon, ExitIcon } from "@radix-ui/react-icons";

export const metadata = {
  title: "All Workspaces - Zyg AI",
};

export default function CreateOrSelectWorkspaceLayout({ children }) {
  return (
    <div className="relative flex min-h-screen flex-col bg-background">
      <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 max-w-screen-2xl items-center">
          <div className="flex flex-1 items-end">
            <Image alt="Zyg Logo" width={32} height={32} src="/logo.png" />
            <div className="ml-2 font-bold">Zyg.</div>
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
