"use client";

import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";

import { ArrowLeftIcon } from "@radix-ui/react-icons";

export function GoBack() {
  const router = useRouter();
  return (
    <Button onClick={() => router.back()} variant="outline" size="icon">
      <ArrowLeftIcon className="h-4 w-4" />
    </Button>
  );
}
