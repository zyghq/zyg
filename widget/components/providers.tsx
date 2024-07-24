"use client";

import * as React from "react";
import {
  QueryClient,
  QueryClientProvider,
  useMutation,
} from "@tanstack/react-query";
import { Customer, CustomerContext, SdkCustomerResponse } from "@/lib/customer";
import { useQuery } from "@tanstack/react-query";

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

async function initWidgetRequest(payload: SdkCustomerResponse) {
  const { widgetId, ...body } = payload;
  const response = await fetch(`/api/widgets/${widgetId}/init`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    throw new Error("Not Found");
  }

  const responseData = await response.json();
  const { jwt, name } = responseData;
  return {
    widgetId,
    jwt,
    name,
  };
}

export function CustomerProvider({ children }: { children: React.ReactNode }) {
  const [customer, setCustomer] = React.useState<Customer | null>(null);
  const [isLoading, setIsLoading] = React.useState(true);
  const [hasError, setHasError] = React.useState(false);

  // React.useEffect(() => {
  //   window.parent.postMessage("ifc:ready", "*");
  // }, []);

  const _ = useQuery({
    queryKey: ["ifc:ready"],
    queryFn: async () => {
      window.parent.postMessage("ifc:ready", "*");
      return true;
    },
  });

  const initMutate = useMutation({
    mutationFn: async (payload: SdkCustomerResponse) => {
      const response = await initWidgetRequest(payload);
      const customer = {
        ...response,
        anonId: payload.anonId,
        customeExternalId: payload.customerExternalId,
        customerEmail: payload.customerEmail,
        customerPhone: payload.customerPhone,
      };
      return customer;
    },
    onSuccess: (result) => {
      setCustomer(result);
      window.parent.postMessage("ifc:ack", "*");
    },
    onError: (error, variables, context) => {
      console.log("onError", error, variables, context);
      setHasError(true);
    },
  });

  React.useEffect(() => {
    const onMessageHandler = async (e: MessageEvent) => {
      try {
        console.log("*********** ifc:onMessageHandler ***********");
        const data = JSON.parse(e.data);
        if (data.type === "customer") {
          const response = JSON.parse(data.data) as SdkCustomerResponse;
          initMutate.mutate({ ...response });
        }
        if (data.type === "start") {
          setIsLoading(false);
        }
      } catch (err) {
        console.error("error procesing evt message:", err);
      } finally {
        console.log("*********** ifc:onMessageHandler ***********");
      }
    };
    window.addEventListener("message", onMessageHandler);
    return () => {
      window.removeEventListener("message", onMessageHandler);
    };
  }, [initMutate]);

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
