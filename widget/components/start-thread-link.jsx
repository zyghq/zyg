import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function StartThreadLink() {
  return (
    <Link href="/threads/">
      <Button
        variant="secondary"
        className="w-full bg-blue-700 hover:bg-blue-800 text-white"
      >
        Send us a message
      </Button>
    </Link>
  );
}
