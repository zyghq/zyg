"use client";
import * as React from "react";
import { useMutation } from "@tanstack/react-query";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Icons } from "@/components/icons";
import {
  ArrowLeftIcon,
  ActivityLogIcon,
  PaperPlaneIcon,
} from "@radix-ui/react-icons";
import { ThumbsUpIcon, ThumbsDownIcon } from "lucide-react";
import Link from "next/link";
import { ScrollArea } from "@/components/ui/scroll-area";
import axios from "axios";

export default function SearchBar() {
  const [text, setText] = React.useState("");
  const [query, setQuery] = React.useState("");
  const [requestId, setRequestId] = React.useState("");

  const queryMutation = useMutation({
    mutationFn: (text) => {
      return axios.post(
        "http://127.0.0.1:8080/-/threads/qa/",
        { query: text },
        {
          headers: {
            "Content-Type": "application/json",
            Authorization:
              "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ3b3Jrc3BhY2VJZCI6Indya2NvNjBlcGt0aWR1N3NvZDk2bDkwIiwiZXh0ZXJuYWxJZCI6Inh4eHgtMTExLXp6enoiLCJlbWFpbCI6InNhbmNoaXRycmtAZ21haWwuY29tIiwicGhvbmUiOiIrOTE3NzYwNjg2MDY4IiwiaXNzIjoiYXV0aC56eWcuYWkiLCJzdWIiOiJjX2NvNjFhYmt0aWR1MXQzaTNkbjYwIiwiYXVkIjpbImN1c3RvbWVyIl0sImV4cCI6MTc0Mzc1Nzg3MSwibmJmIjoxNzEyMjIxODcxLCJpYXQiOjE3MTIyMjE4NzEsImp0aSI6Indya2NvNjBlcGt0aWR1N3NvZDk2bDkwOmNfY282MWFia3RpZHUxdDNpM2RuNjAifQ.epCQ4aXvYPXIhVrX6TtfYrq0XxYXT18kIWsOae8HvUQ",
          },
        }
      );
    },
    onSuccess: (response) => {
      const { data } = response;
      const { threadId, query } = data;
      setQuery(query);
      setRequestId(threadId);
      setText("");
      // setQuery(text);
      // const { data } = response;
      // const { text: textResponse, requestId } = data;
      // console.log(textResponse);
      // console.log(requestId);
      // setText("");
      // setRequestId(requestId);
    },
  });

  const evalMutation = useMutation({
    mutationFn: ({ requestId, score }) => {
      return axios.post(
        `http://127.0.0.1:8080/workspaces/3a690e9f85544f6f82e6bdc432418b11/-/queries/${requestId}/`,
        { eval: score }
      );
    },
    onSuccess: () => {},
  });

  const onSubmit = (event) => {
    event.preventDefault();
    queryMutation.mutate(text);
  };

  const evaluate = (score) => {
    evalMutation.mutate({ requestId, score });
  };

  const renderResult = () => {
    if (queryMutation.isPending) {
      return (
        <div className="flex justify-start mt-4 px-2">
          <div className="flex items-center space-x-1">
            <div className="w-1 h-1 bg-gray-600 rounded-full animate-pulse" />
            <div className="w-1 h-1 bg-gray-600 rounded-full animate-pulse" />
            <div className="w-1 h-1 bg-gray-600 rounded-full animate-pulse" />
          </div>
        </div>
      );
    }

    if (queryMutation.isError) {
      return (
        <div className="flex flex-col justify-center items-center mt-4 px-2 space-y-1">
          <Icons.oops className="h-10 w-10" />
          <div className="text-xs">something went wrong.</div>
        </div>
      );
    }
    if (queryMutation.isSuccess) {
      const { data: result } = queryMutation;
      const { data } = result;
      const { answers = [] } = data;
      const answer = answers[0];
      return (
        <ScrollArea className="h-[calc(100dvh-14rem)]">
          <div className="flex flex-col px-4">
            <div className="mt-2">
              <div className="text-xl">{query}</div>
            </div>
            <div className="flex mt-4 items-center">
              <ActivityLogIcon className="h-4 w-4 mr-2" />
              <div className="font-semibold">Answers</div>
            </div>
            <div className="my-2">
              <p className="whitespace-pre-line">{answer.answer}</p>
            </div>
            <div className="flex justify-end items-center">
              <div className="font-semibold">Helpful?</div>
              <Button
                variant="ghost"
                className="mx-2"
                onClick={() => evaluate(1)}
              >
                <ThumbsUpIcon className="h-4 w-4" />
              </Button>
              <Button variant="ghost" onClick={() => evaluate(0)}>
                <ThumbsDownIcon className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </ScrollArea>
      );
    }

    return (
      <div className="flex flex-col items-center mt-auto">
        <Icons.nothing className="h-40 w-40" />
      </div>
    );
  };

  return (
    <React.Fragment>
      <form onSubmit={onSubmit}>
        <div className="flex items-center justify-between py-4 px-1 border-b">
          <Button variant="outline" size="sm" className="mr-1" asChild>
            <Link href="/">
              <ArrowLeftIcon className="h-4 w-4" />
            </Link>
          </Button>
          <Textarea
            name="text"
            type="text"
            placeholder="Ask anything..."
            rows={1}
            cols={1}
            value={text}
            onChange={(e) => setText(e.target.value)}
            required
          />
          <Button
            variant="outline"
            size="sm"
            className="ml-1"
            type="submit"
            disabled={queryMutation.isPending}
          >
            <PaperPlaneIcon className="h-4 w-4" />
          </Button>
        </div>
      </form>
      {renderResult()}
    </React.Fragment>
  );
}
