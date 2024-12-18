import { Button } from "@/components/ui/button.tsx";
import { Power } from "lucide-react";

export function EnableEmail() {
  return (
    <Button>
      <Power className="mr-2 h-5 w-5" /> Enable Email
    </Button>
  );
}
