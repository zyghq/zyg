import { eventSeverityIcon } from "@/components/icons";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { getCustomerEvents } from "@/db/api";
import { CustomerEventResponse } from "@/db/schema";
import { ClockIcon, InfoCircledIcon } from "@radix-ui/react-icons";
import { useQuery } from "@tanstack/react-query";

function EventCard({ event }: { event: CustomerEventResponse }) {
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
      <CardContent className="space-y-1 p-2">
        <div className="flex items-center justify-between">
          <Badge className="text-xs" variant="outline">
            {eventSeverityIcon(event.severity, { className: "mr-1 h-3 w-3" })}
            {event.severity}
          </Badge>
          <span className="flex items-center text-xs text-muted-foreground">
            <ClockIcon className="mr-1 h-3 w-3" />
            {formatDate(event.timestamp)}
          </span>
        </div>
        <div className="text-sm font-medium">{event.event}</div>
        <p className="font-mono text-sm text-muted-foreground">{event.body}</p>
        <div className="font-mono text-xs text-muted-foreground">
          Event ID {event.eventId}
        </div>
      </CardContent>
    </Card>
  );
}

export function CustomerEvents({
  customerId,
  jwt,
  workspaceId,
}: {
  customerId: string;
  jwt: string;
  workspaceId: string;
}) {
  const {
    data: events,
    error,
    isPending,
  } = useQuery({
    enabled: !!customerId,
    initialData: [],
    queryFn: async () => {
      const { data, error } = await getCustomerEvents(
        jwt,
        workspaceId,
        customerId,
      );
      if (error) throw new Error("failed to fetch customer events");
      return data;
    },
    queryKey: ["customerEvents", workspaceId, customerId, jwt],
    staleTime: 1000 * 60,
  });

  if (error) {
    return <div>Error</div>;
  }

  if (isPending) {
    return <div>Loading...</div>;
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

  return (
    <>
      <div className="flex items-center gap-2">
        <InfoCircledIcon className="h-4 w-4" />
        <div className="font-mono text-xs">No events yet.</div>
      </div>
      <div className="pb-2 text-xs">
        Learn{" "}
        <a
          className="underline"
          href="https://zyg.ai/docs/events?utm_source=app&utm_medium=docs&utm_campaign=onboarding"
          rel="noopener noreferrer"
          target="_blank"
        >
          how to send customer events to Zyg.
        </a>
      </div>
    </>
  );
}
