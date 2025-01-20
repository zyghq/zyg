import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable";

export function ThreadContent() {
  return (
    <div>....</div>
    // <ResizablePanelGroup className="border" direction="horizontal">
    //   <ResizablePanel defaultSize={50}>
    //     <div className="flex h-[200px] items-center justify-center p-6">
    //       <span className="font-semibold">One</span>
    //     </div>
    //   </ResizablePanel>
    //   <ResizableHandle />
    //   <ResizablePanel defaultSize={50}>
    //     <div className="flex h-[200px] items-center justify-center p-6">
    //       <span className="font-semibold">One</span>
    //     </div>
    //   </ResizablePanel>
    // </ResizablePanelGroup>
  );
}
