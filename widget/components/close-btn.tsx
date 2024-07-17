"use client";

import { Button } from "@/components/ui/button";
import { Cross1Icon } from "@radix-ui/react-icons";

export default function CloseButton() {
  const handleClose = () => {
    window.parent.postMessage("close", "*");
  };

  return (
    <Button variant="outline" size="icon" onClick={handleClose}>
      <Cross1Icon className="h-4 w-4" />
    </Button>
  );
}
