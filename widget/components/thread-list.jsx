"use client";
import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Link from "next/link";
import Cookies from "js-cookie";

export default function ThreadList({ threads }) {
  const result = useQuery({
    queryKey: ["thchats"],
    queryFn: async () => {
      const token = Cookies.get("__zygtoken") || "";
      const response = await fetch(`http://localhost:8080/-/threads/chat/`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });
      if (!response.ok) {
        throw new Error("Failed to fetch threads");
      }
      return response.json();
    },
    initialData: threads,
  });

  console.log("result", result.data);

  if (result.isPending) {
    return <div>Loading...</div>;
  }

  if (result.isError) {
    return <div>Error</div>;
  }

  return (
    <div className="flex space-y-1">
      <Card className="w-full">
        <CardHeader>
          <CardTitle>Thread 1</CardTitle>
          <CardDescription>
            Deploy your new project in one-click.
          </CardDescription>
        </CardHeader>
      </Card>
    </div>
  );
}
