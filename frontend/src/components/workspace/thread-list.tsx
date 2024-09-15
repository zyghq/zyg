import { ThreadLinkItem } from "@/components/workspace/thread-list-item";
import { Virtuoso } from "react-virtuoso";
import { Thread } from "@/db/models";

export function ThreadListV3({
  workspaceId,
  threads,
}: {
  workspaceId: string;
  threads: Thread[];
}) {
  return (
    <Virtuoso
      useWindowScroll
      totalCount={threads.length}
      itemContent={(index) => (
        <ThreadLinkItem workspaceId={workspaceId} thread={threads[index]} />
      )}
    />
  );
}
