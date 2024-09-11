// import { cn } from "@/lib/utils";
// import { Link } from "@tanstack/react-router";
// import { Badge } from "@/components/ui/badge";
// import { formatDistanceToNow } from "date-fns";
// import { Thread } from "@/db/models";
// import { ChatBubbleIcon, ResetIcon } from "@radix-ui/react-icons";
// import Avatar from "boring-avatars";
// import { useStore } from "zustand";
// import { useWorkspaceStore } from "@/providers";

// function ThreadItem({
//   workspaceId,
//   item,
//   variant = "default",
// }: {
//   workspaceId: string;
//   item: Thread;
//   variant?: string;
// }) {
//   // const WorkspaceStore = useRouteContext({
//   //   from: "/_auth/workspaces/$workspaceId/_workspace",
//   //   select: (context) => context.WorkspaceStore,
//   // });
//   const workspaceStore = useWorkspaceStore();
//   const customerName = useStore(workspaceStore, (state) =>
//     state.viewCustomerName(state, item.customerId)
//   );

//   //   const message = item.messages[0];
//   //   const name = item?.customer?.name || "Customer";
//   //   const { assignee } = item;

//   //   const renderLabels = () => {
//   //     if (result.isSuccess && result.data && result.data.length) {
//   //       return (
//   //         <div className="flex gap-1">
//   //           {result.data.map((label) => (
//   //             <Badge
//   //               key={label.labelId}
//   //               variant="outline"
//   //               className="font-normal"
//   //             >
//   //               {label.name}
//   //             </Badge>
//   //           ))}
//   //         </div>
//   //       );
//   //     }
//   //     return <div className="min-h-5"></div>;
//   //   };

//   return (
//     <Link
//       to={"/workspaces/$workspaceId/threads/$threadId"}
//       params={{ workspaceId, threadId: item.threadId }}
//       className={cn(
//         "flex flex-col items-start gap-2 rounded-lg border px-3 py-3 text-left text-sm transition-all hover:bg-accent",
//         variant === "compress" && "gap-0 rounded-none py-5"
//       )}
//     >
//       <div className="flex w-full flex-col gap-1">
//         <div className="flex items-center">
//           <div className="flex items-center gap-2">
//             <ChatBubbleIcon />
//             <div className="font-semibold">{customerName}</div>
//             {!item.read && (
//               <span className="flex h-2 w-2 rounded-full bg-blue-600" />
//             )}
//           </div>
//           <div
//             className={cn(
//               "ml-auto mr-2 text-xs",
//               !item.replied ? "text-foreground" : "text-muted-foreground"
//             )}
//           >
//             {formatDistanceToNow(new Date(item.updatedAt), {
//               addSuffix: true,
//             })}
//           </div>
//           {item.assigneeId && (
//             <Avatar size={24} name={item.assigneeId} variant="marble" />
//           )}
//         </div>
//         {item.replied ? (
//           <div className="flex">
//             <Badge variant="outline" className="font-normal">
//               <div className="flex items-center gap-1">
//                 <ResetIcon className="h-3 w-3" />
//                 replied to
//               </div>
//             </Badge>
//           </div>
//         ) : (
//           <div className="flex">
//             <Badge
//               variant="outline"
//               className="bg-indigo-100 font-normal dark:bg-indigo-500"
//             >
//               <div className="flex items-center gap-1">
//                 <ResetIcon className="h-3 w-3" />
//                 awaiting reply
//               </div>
//             </Badge>
//           </div>
//         )}
//         {/* {item?.title ? <div className="font-medium">{item?.title}</div> : null} */}
//       </div>
//       <div className="line-clamp-2 text-muted-foreground">
//         {item.previewText}
//       </div>
//     </Link>
//   );
// }

// export function ThreadListV2({
//   workspaceId,
//   threads,
//   variant = "default",
// }: {
//   workspaceId: string;
//   threads: Thread[];
//   variant?: string;
// }) {
//   return (
//     <div
//       className={cn("flex flex-col gap-2", variant === "compress" && "gap-0")}
//     >
//       {threads.map((item: Thread) => (
//         <ThreadItem
//           key={item.threadId}
//           workspaceId={workspaceId}
//           item={item}
//           variant={variant}
//         />
//       ))}
//     </div>
//   );
// }

import { Mail } from "lucide-react";

interface SupportRequest {
  id: string;
  sender: string;
  subject: string;
  preview: string;
  tags: string[];
  time: string;
}

const supportRequests: SupportRequest[] = [
  {
    id: "1",
    sender: "Sanchit Rk",
    subject: "Help Support for Zyg",
    preview: "Hello, I'm having trouble logging into th...",
    tags: ["Needs first response"],
    time: "5d",
  },
  {
    id: "2",
    sender: "Manmohini Sharma",
    subject: "Trouble with DB queries",
    preview: "There seems to be some issues with...",
    tags: ["Needs first response", "Bug", "Technical"],
    time: "5d",
  },
  {
    id: "3",
    sender: "Sanchit Rk",
    subject: "Support request",
    preview: "My second mail, when should I check the lo...",
    tags: ["Needs first response"],
    time: "5d",
  },
  {
    id: "4",
    sender: "Sanchit Rk",
    subject: "Support request",
    preview: "Hey Thanks the issue is now resolve, Login ...",
    tags: ["Needs first response"],
    time: "5d",
  },
  {
    id: "5",
    sender: "Sanchit Rk",
    subject: "Support request",
    preview: "No preview",
    tags: ["Needs first response"],
    time: "5h",
  },
];

export function ThreadListV2() {
  return (
    <div>
      {supportRequests.map((request) => (
        <div className="block unicode-bidi-isolate border-b" key={request.id}>
          <div className="relative top-10 left-4">
            <Mail className="w-4 h-4 text-muted-foreground" />
          </div>
          <a
            href="#"
            className="grid gap-x-4 items-start py-4 px-8 hover:bg-zinc-50 no-underline"
          >
            ...
          </a>
        </div>
        // <a
        //   key={request.id}
        //   href={`#${request.id}`}
        //   className="grid items-center border-b border-gray-200 hover:bg-gray-50"
        //   style={{
        //     padding: "10px 32px",
        //     gridTemplateColumns: "18px 200px auto 240px",
        //     gridTemplateRows: "1fr",
        //   }}
        // >
        //   <Mail className="w-4 h-4 text-muted-foreground" />
        //   <div className="font-medium text-sm truncate">{request.sender}</div>
        //   <div className="flex items-center space-x-2 overflow-hidden">
        //     <span className="font-medium text-sm whitespace-nowrap">
        //       {request.subject}
        //     </span>
        //     <span className="text-sm text-gray-500 truncate">
        //       {request.preview}
        //     </span>
        //   </div>
        //   <div className="flex items-center justify-end space-x-2">
        //     <div className="flex flex-wrap justify-end gap-1">
        //       {request.tags.map((tag, index) => (
        //         <span
        //           key={index}
        //           className="px-2 py-1 text-xs font-medium text-blue-600 bg-blue-100 rounded-full"
        //         >
        //           {tag}
        //         </span>
        //       ))}
        //     </div>
        //     <span className="text-sm text-gray-500 ml-2">{request.time}</span>
        //   </div>
        // </a>
      ))}
    </div>
  );
}
