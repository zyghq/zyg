import { createFileRoute } from "@tanstack/react-router";
import { TrendingUp } from "lucide-react";
import { QueueSize, Volume } from "@/components/workspace/insights/overview";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertCircle, MessageCircle } from "lucide-react";

export const Route = createFileRoute(
  "/_account/workspaces/$workspaceId/_workspace/insights"
)({
  component: InsightsComponent,
});

export default function InDevelopment() {
  return (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <AlertCircle className="h-5 w-5 text-yellow-500" />
          <span>Insights Coming Soon</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">
          We're currently developing the insights page. It will be enabled once
          it's ready. Stay tuned!
        </p>
      </CardContent>
      <CardFooter className="flex flex-col sm:flex-row gap-4">
        <Button className="w-full sm:w-auto">Talk to Us</Button>
        <a
          href="https://github.com/zyghq/zyg/discussions"
          target="_blank"
          rel="noopener noreferrer"
          className="w-full sm:w-auto"
        >
          <Button variant="outline" className="w-full">
            <MessageCircle className="mr-2 h-4 w-4" />
            Start Discussion in GitHub
          </Button>
        </a>
      </CardFooter>
    </Card>
  );
}

function InsightsComponent() {
  return (
    <div className="filter blur-md">
      <div className="px-2 sm:px-4 py-4 flex flex-col gap-4">
        <Card className="shadow-sm">
          <CardHeader>
            <CardTitle>Overview</CardTitle>
            <div className="flex justify-between">
              <div>
                <div className="text-xs text-muted-foreground">
                  A snapshot of the number of threads in Todo.
                </div>
                <div className="text-muted-foreground text-xs">
                  <span className="font-mono text-2xl text-foreground">4</span>{" "}
                  in Todo right now.
                </div>
              </div>
              <div>...</div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap sm:flex-nowrap">
              <QueueSize className="w-full sm:w-1/3 max-h-48" />
              <QueueSize className="w-full sm:w-1/3 max-h-48" />
              <QueueSize className="w-full sm:w-1/3 max-h-48" />
            </div>
          </CardContent>
          <CardFooter>
            <div className="flex w-full items-start gap-2 text-sm">
              <div className="grid gap-2">
                <div className="flex items-center gap-2 font-medium leading-none">
                  Trending up by 5.2% this month
                  <TrendingUp className="h-4 w-4" />
                </div>
                <div className="flex items-center gap-2 leading-none text-muted-foreground">
                  January - June 2024
                </div>
              </div>
            </div>
          </CardFooter>
        </Card>
        <Card className="shadow-sm">
          <CardHeader>
            <CardTitle>Support Volume</CardTitle>
            <div className="flex justify-between">
              <div>
                <div className="text-xs text-muted-foreground">
                  The volume of new threads created each day.
                </div>
              </div>
              <div>...</div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap sm:flex-nowrap">
              <Volume className="w-full sm:w-1/3 max-h-48" />
              <Volume className="w-full sm:w-1/3 max-h-48" />
              <Volume className="w-full sm:w-1/3 max-h-48" />
            </div>
          </CardContent>
          <CardFooter>
            <div className="flex w-full items-start gap-2 text-sm">
              <div className="grid gap-2">
                <div className="flex items-center gap-2 font-medium leading-none">
                  Trending up by 5.2% this month
                  <TrendingUp className="h-4 w-4" />
                </div>
                <div className="flex items-center gap-2 leading-none text-muted-foreground">
                  January - June 2024
                </div>
              </div>
            </div>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
