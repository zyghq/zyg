"use client";

import * as React from "react";
import { Customer, CustomerContext, SdkCustomerResponse } from "@/lib/customer";

export function CustomerProvider({ children }: { children: React.ReactNode }) {
  const [customer, setCustomer] = React.useState<Customer | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);
  const [hasError, setHasError] = React.useState(false);

  React.useEffect(() => {
    window.parent.postMessage("ifc:ready", "*");
  }, []);

  React.useEffect(() => {
    const onMessageHandler = (e: MessageEvent) => {
      try {
        console.log("*********** ifc:onMessageHandler ***********");
        const data = JSON.parse(e.data);

        if (data.type === "customer") {
          const response = JSON.parse(data.data) as SdkCustomerResponse;
          const { widgetId, ...body } = response;
          fetch(
            `${process.env.NEXT_PUBLIC_XAPI_URL}/widgets/${widgetId}/init/`,
            {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify(body),
            }
          )
            .then((response) => {
              if (!response.ok) {
                throw new Error("Not Found");
              }
              return response.json();
            })
            .then((data) => {
              console.log("response data", data);
              const { jwt, name } = data;
              setCustomer({
                jwt,
                name,
              });
              // send post message to the parent
              window.parent.postMessage("ifc:ack", "*");
            })
            .catch((err) => {
              console.error("Error processing request:", err);
              setHasError(true);
              // send post message to the parent
              window.parent.postMessage("ifc:error", "*");
            });
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
