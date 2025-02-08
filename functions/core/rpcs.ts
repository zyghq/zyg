/**
 * This file contains RPC (Remote Procedure Call) functions
 * used to communicate with other services.
 */

import { ZYG_SRV_BASE_URL } from "./config.ts";

export async function upsertThreadSummaryRPC() {
  const payload = {
    jsonrpc: "2.0",
    method: "upsertThreadSummary",
    params: {},
    id: Date.now(),
  };

  const response = await fetch(`${ZYG_SRV_BASE_URL}/v1/rpc/threads/`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  const result = await response.json();
  if (result.error) {
    throw new Error(result.error.message);
  }

  return result.result;
}
