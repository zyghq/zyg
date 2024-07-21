"use client";

import * as React from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Customer, CustomerContext, SdkCustomerResponse } from "@/lib/customer";

export const ReactQueryClientProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [queryClient] = React.useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            staleTime: 1000 * 60 * 5,
          },
        },
      })
  );
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
};

export function CustomerProvider({ children }: { children: React.ReactNode }) {
  const [customer, setCustomer] = React.useState<Customer | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);
  const [hasError, setHasError] = React.useState(false);

  React.useEffect(() => {
    window.parent.postMessage("ifc:ready", "*");
  }, []);

  React.useEffect(() => {
    const onMessageHandler = async (e: MessageEvent) => {
      try {
        console.log("*********** ifc:onMessageHandler ***********");
        const data = JSON.parse(e.data);
        if (data.type === "customer") {
          try {
            const response = JSON.parse(data.data) as SdkCustomerResponse;
            const { widgetId, ...body } = response;
            const fetchResponse = await fetch(
              `/api/widgets/${widgetId}/init/`,
              {
                method: "POST",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify(body),
              }
            );

            if (!fetchResponse.ok) {
              throw new Error("Not Found");
            }

            const responseData = await fetchResponse.json();
            console.log("response data", responseData);

            const { jwt, name } = responseData;
            setCustomer({
              widgetId,
              jwt,
              name,
            });

            // send post message to the parent
            window.parent.postMessage("ifc:ack", "*");
          } catch (err) {
            console.error("Error processing request:", err);
            setHasError(true);
            // send post message to the parent
            window.parent.postMessage("ifc:error", "*");
          }
        }
        if (data.type === "start") {
          setIsLoading(false);
        }
      } catch (err) {
        console.error("evt on message something went wrong:", err);
      } finally {
        console.log("*********** ifc:onMessageHandler ***********");
      }
    };
    window.addEventListener("message", onMessageHandler);
    return () => {
      window.removeEventListener("message", onMessageHandler);
    };
  }, []);

  const value = {
    isLoading,
    hasError,
    customer,
  };

  return (
    <CustomerContext.Provider value={value}>
      {children}
    </CustomerContext.Provider>
  );
}
