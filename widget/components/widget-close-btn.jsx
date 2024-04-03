"use client";
import { Button } from "@/components/ui/button";
import { Cross1Icon } from "@radix-ui/react-icons";

export default function WidgetCloseButton() {
  const handleClose = () => {
    console.log("close called from iframe ....");
    window.parent.postMessage("close", "*");
  };
  return (
    <Button size="icon" variant="ghost" onClick={handleClose}>
      <Cross1Icon className="h-4 w-4" />
    </Button>
  );
}
