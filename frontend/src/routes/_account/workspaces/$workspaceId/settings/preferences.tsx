import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Separator } from "@/components/ui/separator";
import { useTheme } from "@/hooks/theme";
import { LaptopIcon, MoonIcon, SunIcon } from "@radix-ui/react-icons";
import { createFileRoute } from "@tanstack/react-router";
import { CheckIcon } from "lucide-react";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/settings/preferences",
)({
  component: Preferences,
});

function ThemeSelector() {
  const { setTheme, theme } = useTheme();
  const systemTheme = window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "dark"
    : "light";

  const renderSystem = (theme: { label: string; value: string }) => {
    return (
      <>
        <RadioGroupItem
          className="peer sr-only"
          id={theme.value}
          onClick={() => setTheme("system")}
          value={theme.value}
        />
        <Label
          className={`flex flex-grow cursor-pointer flex-col items-center justify-between rounded-md border p-4 ${systemTheme === "dark" ? "bg-primary text-white dark:bg-neutral-900" : "dark:bg-white dark:text-primary-foreground"}`}
          htmlFor={theme.value}
        >
          <LaptopIcon className="mb-3 h-5 w-5" />
          <span className="text-lg font-bold">Aa</span>
          <p className="text-sm">{theme.label}</p>
        </Label>
        <div className="absolute right-2 top-2 z-10 flex h-5 w-5 items-center justify-center rounded-full bg-primary opacity-0 transition-opacity peer-data-[state=checked]:opacity-100">
          <CheckIcon className="h-4 w-4 text-primary-foreground" />
        </div>
      </>
    );
  };

  const renderLight = (theme: { label: string; value: string }) => {
    return (
      <>
        <RadioGroupItem
          className="peer sr-only"
          id={theme.value}
          onClick={() => setTheme("light")}
          value={theme.value}
        />
        <Label
          className={`flex flex-grow cursor-pointer flex-col items-center justify-between rounded-md border p-4 dark:bg-white dark:text-primary-foreground`}
          htmlFor={theme.value}
        >
          <SunIcon className="mb-3 h-5 w-5" />
          <span className="text-lg font-bold">Aa</span>
          <p className="text-sm">{theme.label}</p>
        </Label>
        <div className="absolute right-2 top-2 z-10 flex h-5 w-5 items-center justify-center rounded-full bg-primary opacity-0 transition-opacity peer-data-[state=checked]:opacity-100">
          <CheckIcon className="h-4 w-4 text-primary-foreground" />
        </div>
      </>
    );
  };

  const renderDark = (theme: { label: string; value: string }) => {
    return (
      <>
        <RadioGroupItem
          className="peer sr-only"
          id={theme.value}
          onClick={() => setTheme("dark")}
          value={theme.value}
        />
        <Label
          className={`flex flex-grow cursor-pointer flex-col items-center justify-between rounded-md border bg-primary p-4 text-white dark:bg-neutral-900`}
          htmlFor={theme.value}
        >
          <MoonIcon className="mb-3 h-5 w-5" />
          <span className="text-lg font-bold">Aa</span>
          <p className="text-sm">{theme.label}</p>
        </Label>
        <div className="absolute right-2 top-2 z-10 flex h-5 w-5 items-center justify-center rounded-full bg-primary opacity-0 transition-opacity peer-data-[state=checked]:opacity-100">
          <CheckIcon className="h-4 w-4 text-primary-foreground" />
        </div>
      </>
    );
  };

  return (
    <RadioGroup
      className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"
      defaultValue={theme}
    >
      {[
        {
          label: "System",
          value: "system",
        },
        { label: "Light", value: "light" },
        { label: "Dark", value: "dark" },
      ].map((theme) => (
        <div className="relative flex" key={theme.value}>
          {theme.value === "system"
            ? renderSystem(theme)
            : theme.value === "light"
              ? renderLight(theme)
              : renderDark(theme)}
        </div>
      ))}
    </RadioGroup>
  );
}

function Preferences() {
  return (
    <div className="container">
      <div className="max-w-2xl">
        <div className="my-12">
          <div className="my-12">
            <header className="text-xl font-semibold">Preferences</header>
          </div>
          <Separator />
        </div>
        <div className="flex flex-col gap-2">
          <div className="text-xl">Theme</div>
          <ThemeSelector />
        </div>
      </div>
    </div>
  );
}
