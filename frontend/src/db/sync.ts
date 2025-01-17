import { ShapeStreamOptions } from "@electric-sql/client";

export type MemberRow = {
  avatar_url: string;
  created_at: string;
  member_id: string;
  name: string;
  permissions: Record<string, unknown>;
  public_name: string;
  role: string;
  synced_at: string;
  updated_at: string;
  version_id: string;
  workspace_id: string;
};

export type MemberRowUpdates = {
  avatar_url?: string;
  created_at?: string;
  member_id?: string;
  name?: string;
  permissions?: Record<string, unknown>;
  public_name?: string;
  role?: string;
  synced_at?: string;
  updated_at?: string;
};

export function syncMembersShape({
  token,
  workspaceId,
}: {
  token: string;
  workspaceId: string;
}): ShapeStreamOptions {
  const url = `${import.meta.env.VITE_ZYG_URL}/v1/sync/workspaces/${workspaceId}/shapes/parts/members/`;
  return {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    params: {
      table: "member",
    },
    url,
  };
}

export type CustomerRow = {
  avatar_url: string;
  created_at: string;
  customer_id: string;
  email: string;
  external_id: string;
  is_email_verified: boolean;
  name: string;
  phone: string;
  role: string;
  synced_at: string;
  updated_at: string;
  version_id: string;
  workspace_id: string;
};

export type CustomerRowUpdates = {
  avatar_url?: string;
  created_at?: string;
  customer_id?: string;
  email?: string;
  external_id?: string;
  is_email_verified?: boolean;
  name?: string;
  phone?: string;
  role?: string;
  synced_at?: string;
  updated_at?: string;
};

export function syncCustomersShape({
  token,
  workspaceId,
}: {
  token: string;
  workspaceId: string;
}): ShapeStreamOptions {
  const url = `${import.meta.env.VITE_ZYG_URL}/v1/sync/workspaces/${workspaceId}/shapes/parts/customers/`;
  return {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    params: {
      table: "customer",
    },
    url,
  };
}
