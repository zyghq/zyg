import * as React from "react";
import { cookies } from "next/headers";
import { isAuthenticated } from "@/utils/helpers";
import { ThreadHeader } from "@/components/headers";
import StartThreadForm from "@/components/start-thread-form";

export default function ThreadInitPage() {
  const cookieStore = cookies();

  if (!isAuthenticated(cookieStore)) {
    return redirect("/authenticate/");
  }
  return (
    <React.Fragment>
      <ThreadHeader />
      <div className="pt-2 px-2 mt-auto border-t">
        <StartThreadForm />
        <footer className="flex flex-col justify-center items-center border-t w-full h-8 mt-2">
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
