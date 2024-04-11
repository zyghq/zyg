"use client";

import * as React from "react";
import { Icons } from "@/components/icons";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

export default function AuthenticatePage() {
  const router = useRouter();
  const result = useQuery({
    queryKey: ["me"],
    queryFn: async () => {
      const response = await fetch(`/api/auth/`, {
        method: "GET",
      });
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      return response.json();
    },
  });

  React.useEffect(() => {
    if (result.isSuccess) {
      return router.push("/");
    }
  }, [result.isSuccess, router]);

  if (result.isError) {
    return (
      <div className="flex flex-col mx-auto my-auto">
        <div className="flex flex-col items-center">
          <Icons.oops className="w-8 h-8" />
          <div className="text-xs text-red-500">Authentication Error</div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col mx-auto my-auto">
      <Icons.spinner className="w-5 h-5" />
    </div>
  );
}
