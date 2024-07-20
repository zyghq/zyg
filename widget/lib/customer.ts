import * as React from "react";

type KV = { [key: string]: string };

export interface SdkCustomerResponse {
  widgetId: string;
  anonId?: string;
  customerExternalId?: string;
  customerEmail?: string;
  customerPhone?: string;
  customerHash?: string;
  traits?: KV;
}

export interface Customer {
  name: string;
  jwt: string;
}

export interface CustomerContext {
  customer: Customer | null;
  isLoading: boolean;
  hasError: boolean;
}

export const CustomerContext = React.createContext<CustomerContext | null>(
  null
);

export const useCustomer = () => {
  const context = React.useContext(CustomerContext);
  if (!context) {
    throw new Error("useCustomer must be used within a CustomerProvider");
  }
  return context;
};
