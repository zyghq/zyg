import { RenderComponents } from "@/components/event/components";
import { eventSeverityIcon } from "@/components/icons";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { CustomerEventResponse } from "@/db/schema";
import { ClockIcon } from "@radix-ui/react-icons";
import { DefaultError } from "@tanstack/react-query";
import { AlertCircle } from "lucide-react";

function EventError() {
  return (
    <Alert className="bg-card" variant="destructive">
      <AlertCircle className="h-4 w-4" />
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>Something went wrong.</AlertDescription>
    </Alert>
  );
}

export function CustomerEvents({
  error,
  events,
}: {
  error: DefaultError | null;
  events: CustomerEventResponse[];
}) {
  if (error) {
    return <EventError />;
  }

  if (events && events.length > 0) {
    return (
      <div className="flex flex-col gap-1">
        {events.map((event) => (
          <EventCard event={event} key={event.eventId} />
        ))}
      </div>
    );
  }
  return null;
}

export function EventCard({ event }: { event: CustomerEventResponse }) {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString("en-US", {
      day: "numeric",
      hour: "2-digit",
      hour12: true,
      minute: "2-digit",
      month: "short",
    });
  };
  return (
    <Card className="shadow-xs">
      <CardHeader className="border-b p-3">
        <div className="flex items-center">
          {eventSeverityIcon(event.severity, {
            className: "h-5 w-5",
          })}
          <div className="flex w-full gap-2">
            <span className="text-sm font-medium">{event.title}</span>
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
              <ClockIcon className="h-4 w-4" />
              {formatDate(event.timestamp)}
            </span>
          </div>
        </div>
      </CardHeader>
      <CardContent className="p-3">
        {RenderComponents({ components: event.components })}
      </CardContent>
    </Card>
  );
}
