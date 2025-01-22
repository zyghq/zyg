import { Spinner } from "@/components/spinner";
import { CustomerEvents } from "@/components/thread/customer-events";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { getCustomerEvents } from "@/db/api";
import { cn } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";

interface RecentCustomerEventsProps {
  customerId: string;
  token: string;
  triggerClassname?: string;
  workspaceId: string;
}

export function RecentCustomerEvents({
  customerId,
  token,
  triggerClassname,
  workspaceId,
}: RecentCustomerEventsProps) {
  const {
    data: events,
    error: eventsError,
    isPending: eventsIsPending,
  } = useQuery({
    enabled: !!customerId,
    initialData: [],
    queryFn: async () => {
      const { data, error } = await getCustomerEvents(
        token,
        workspaceId,
        customerId,
      );
      if (error) throw new Error("failed to fetch customer events");
      return data;
    },
    queryKey: ["events", workspaceId, customerId, token],
    refetchOnMount: "always",
    staleTime: 0,
  });

  return (
    <Accordion className="w-full" collapsible type="single">
      <AccordionItem className="border-none" value="item-1">
        <AccordionTrigger
          className={cn("text-xs hover:no-underline", triggerClassname || "")}
        >
          Recent Events
        </AccordionTrigger>
        <AccordionContent>
          {eventsIsPending ? (
            <span>
              <Spinner
                className="animate-spin text-muted-foreground"
                size={14}
              />
            </span>
          ) : (
            <CustomerEvents error={eventsError} events={events || []} />
          )}
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
}
