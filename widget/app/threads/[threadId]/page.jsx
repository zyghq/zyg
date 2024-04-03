import { ThreadHeader } from "@/components/headers";
import Room from "@/components/room";

export default function ThreadPage({ params }) {
  const { threadId } = params;
  console.log("rendering thread page", threadId);
  return (
    <div className="flex flex-col max-h-[calc(100dvh-2rem)]">
      <ThreadHeader />
      <Room roomId={threadId} />
    </div>
  );
}
