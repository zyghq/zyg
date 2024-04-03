import * as React from "react";
import SearchBar from "@/components/searchbar";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function SearchPage() {
  return (
    <React.Fragment>
      <SearchBar />
      <div className="pt-4 px-2 mt-auto border-t">
        <Link href="/threads/threadId/">
          <Button
            variant="secondary"
            className="w-full bg-blue-700 hover:bg-blue-800 text-white"
          >
            Send us a message
          </Button>
        </Link>
        <footer className="flex flex-col justify-center items-center border-t w-full h-8 mt-4">
          <a
            href="https://www.zyg.ai/"
            className="text-xs font-semibold text-gray-500"
            target="_blank"
          >
            Powered by Zyg.
          </a>
        </footer>
      </div>
    </React.Fragment>
  );
}
