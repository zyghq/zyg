import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";

export function QuickThreadSummary() {
  return (
    <Accordion className="flex flex-col" collapsible type="single">
      <AccordionItem className="border-none" value="item-1">
        <AccordionTrigger className="px-1.5 py-0 text-xs text-muted-foreground hover:no-underline">
          Summary
        </AccordionTrigger>
        <AccordionContent className="px-1.5 py-0 text-xs text-muted-foreground">
          ...
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  );
}
