import * as React from "react";
import { cookies } from "next/headers";
import { Header } from "@/components/headers";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { MagnifyingGlassIcon } from "@radix-ui/react-icons";
import Link from "next/link";
import { Icons } from "@/components/icons";
import StartThreadLink from "@/components/start-thread-link";
import { isAuthenticated } from "@/utils/helpers";
import { redirect } from "next/navigation";
import ThreadList from "@/components/thread-list";

async function GetThreadsAPI(cookies) {
  if (!cookies) {
    return { error: new Error("no cookie store provided"), data: null };
  }
  const token = cookies.get("__zygtoken");
  if (!token) {
    return { error: new Error("no auth token found or not set"), data: null };
  }
  const { value = "" } = token;

  const url = `${process.env.ZYG_API_URL}/-/threads/chat/`;

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${value}`,
      },
    });
    if (!response.ok) {
      return { error: new Error("Failed to fetch threads"), data: null };
    }
    const data = await response.json();
    return { error: null, data };
  } catch (error) {
    console.error("Error fetching threads", error);
    return { error, data: null };
  }
}

export default async function WelcomePage() {
  const cookieStore = cookies();

  if (!(await isAuthenticated(cookieStore))) {
    return redirect("/authenticate/");
  }

  const threads = [];
  const { error, data } = await GetThreadsAPI(cookieStore);
  if (error) {
    console.error("Error fetching threads", error);
  } else {
    threads.push(...data);
  }

  return (
    <React.Fragment>
      <div className="flex flex-col">
        <Header label="Hey! How can we help?" />
        <div className="px-2">
          <Button variant="outline" className="w-full" asChild>
            <Link href="/search/" className="flex">
              <MagnifyingGlassIcon className="h-4 w-4 mr-1" />
              Ask anything...
            </Link>
          </Button>
        </div>
        <div className="p-2">
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
            <TabsContent value="threads">
              <ThreadList threads={threads} />
            </TabsContent>
          </Tabs>
        </div>
      </div>
      <div className="pt-4 px-2 mt-auto border-t">
        <StartThreadLink />
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
