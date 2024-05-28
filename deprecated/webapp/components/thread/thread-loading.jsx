import { Skeleton } from "@/components/ui/skeleton";

export default function ThreadLoading() {
  return (
    <div className="space-y-2 p-1">
      <div className="flex flex-col gap-2 border p-4">
        <Skeleton className="h-5 w-28" />
        <div className="flex">
          <Skeleton className="h-6 w-6 rounded-full" />
          <div className="ml-2 space-y-2">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-3 w-40" />
            <Skeleton className="h-3 w-60" />
          </div>
        </div>
      </div>

      <div className="flex flex-col gap-2 border p-4">
        <Skeleton className="h-5 w-28" />
        <div className="flex">
          <Skeleton className="h-6 w-6 rounded-full" />
          <div className="ml-2 space-y-2">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-3 w-40" />
            <Skeleton className="h-3 w-60" />
          </div>
        </div>
      </div>

      <div className="flex flex-col gap-2 border p-4">
        <Skeleton className="h-5 w-28" />
        <div className="flex">
          <Skeleton className="h-6 w-6 rounded-full" />
          <div className="ml-2 space-y-2">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-3 w-40" />
            <Skeleton className="h-3 w-60" />
          </div>
        </div>
      </div>

      <div className="flex flex-col gap-2 border p-4">
        <Skeleton className="h-5 w-28" />
        <div className="flex">
          <Skeleton className="h-6 w-6 rounded-full" />
          <div className="ml-2 space-y-2">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-3 w-40" />
            <Skeleton className="h-3 w-60" />
          </div>
        </div>
      </div>
    </div>
  );
}
