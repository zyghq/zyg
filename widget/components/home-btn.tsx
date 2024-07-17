import { ArrowLeftIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function HomeButton() {
  return (
    <Button variant="outline" size="sm" className="mr-1" asChild>
      <Link href="/">
        <ArrowLeftIcon className="h-4 w-4" />
      </Link>
    </Button>
  );
}
