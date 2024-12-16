import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Check, Copy } from "lucide-react";

interface DNSRecord {
  hostname: string;
  status: "Pending" | "Verified";
  type: string;
  value: string;
}

export function DNSRecords({ records }: { records: DNSRecord[]}) {
  return (
    <Card className="rounded-md shadow-none">
      <CardContent className="p-4">
        <div className="space-y-4">
          {records.map((record, index) => (
            <div className="space-y-4" key={index}>
              <div className="grid gap-4">
                <div className="grid gap-2">
                  <div className="text-sm font-medium text-muted-foreground">
                    Type
                  </div>
                  <div className="font-mono">{record.type}</div>
                </div>
                <div className="grid gap-2">
                  <div className="text-sm font-medium text-muted-foreground">
                    Hostname
                  </div>
                  <div className="flex items-start gap-2">
                    <Button
                      className="mt-0.5 h-6 w-6 shrink-0"
                      onClick={() =>
                        navigator.clipboard.writeText(record.hostname)
                      }
                      size="icon"
                      variant="ghost"
                    >
                      <Copy className="h-4 w-4" />
                      <span className="sr-only">Copy hostname</span>
                    </Button>
                    <div className="break-all font-mono">{record.hostname}</div>
                  </div>
                </div>
                <div className="grid gap-2">
                  <div className="text-sm font-medium text-muted-foreground">
                    Value
                  </div>
                  <div className="flex items-start gap-2">
                    <Button
                      className="mt-0.5 h-6 w-6 shrink-0"
                      onClick={() =>
                        navigator.clipboard.writeText(record.value)
                      }
                      size="icon"
                      variant="ghost"
                    >
                      <Copy className="h-4 w-4" />
                      <span className="sr-only">Copy value</span>
                    </Button>
                    <div className="break-all font-mono">{record.value}</div>
                  </div>
                </div>
                <div className="grid gap-2">
                  <div className="text-sm font-medium text-muted-foreground">
                    Status
                  </div>
                  <div className="flex items-center gap-2 text-green-600">
                    <Check className="h-4 w-4" />
                    {record.status}
                  </div>
                </div>
              </div>
              {index < records.length - 1 && <div className="border-t" />}
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
