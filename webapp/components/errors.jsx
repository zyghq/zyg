import { Icons } from "@/components/icons";

export function OopsDefault() {
  return (
    <div className="mt-12 flex flex-col items-center space-y-1">
      <Icons.oops className="h-12 w-12" />
      <div className="text-xs">something went wrong.</div>
    </div>
  );
}
