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
  widgetId: string;
  name: string;
  jwt: string;
  isVerified?: boolean;
  anonId?: string;
  customerExternalId?: string;
  customerEmail?: string;
  customerPhone?: string;
  customerHash?: string;
}

export interface Identities {
  name: string;
  customerEmail: string;
  customerPhone: string;
  customerExternalId: string;
  isVerified: boolean;
}

export interface CustomerContext {
  customer: Customer | null;
  isLoading: boolean;
  hasError: boolean;
  setIdentities: (identities: Identities) => void;
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
