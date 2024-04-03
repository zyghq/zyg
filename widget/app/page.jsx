import * as React from "react";
import { Header } from "@/components/headers";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { MagnifyingGlassIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { Icons } from "@/components/icons";

export default function WelcomePage() {
  return (
    <React.Fragment>
      <div className="flex flex-col">
        <Header label="Hey! How can we help?" />
        <div className="px-4">
          <Button variant="outline" className="w-full" asChild>
            <Link href="/search/" className="flex">
              <MagnifyingGlassIcon className="h-4 w-4 mr-1" />
              Ask anything...
            </Link>
          </Button>
        </div>
        <div className="p-4">
          <Tabs defaultValue="home">
            <TabsList>
              <TabsTrigger value="home">Home</TabsTrigger>
              <TabsTrigger value="threads">Threads</TabsTrigger>
            </TabsList>
            <TabsContent value="home">
              <div className="flex flex-col items-center justify-center space-y-4">
                <Icons.nothing className="w-40" />
                <p className="text-center text-gray-600 dark:text-gray-300">
                  Nothing to see here yet.
                </p>
              </div>
            </TabsContent>
            <TabsContent value="threads">....</TabsContent>
          </Tabs>
        </div>
      </div>
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