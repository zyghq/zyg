import { ThreadLabelShapeMap } from "@/db/shapes";
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

export type ThreadRow = {
  assigned_at: null | string;
  assignee_id: null | string;
  channel: string;
  created_at: string;
  created_by_id: string;
  customer_id: string;
  description: string;
  last_inbound_at: string | null;
  last_outbound_at: string | null;
  labels: null | ThreadLabelShapeMap;
  preview_text: string;
  priority: string;
  replied: boolean;
  stage: string;
  status: string;
  status_changed_at: string;
  status_changed_by_id: string;
  synced_at: string;
  thread_id: string;
  title: string;
  updated_at: string;
  updated_by_id: string;
  version_id: string;
  workspace_id: string;
};

export type ThreadRowUpdates = {
  assigned_at?: null | string;
  assignee_id?: null | string;
  channel?: string;
  created_at?: string;
  created_by_id?: string;
  customer_id?: string;
  description?: string;
  last_inbound_at?: null | string;
  labels?: null | ThreadLabelShapeMap;
  last_outbound_at?: null | string;
  preview_text?: string;
  priority?: string;
  replied?: boolean;
  stage?: string;
  status?: string;
  status_changed_at?: string;
  status_changed_by_id?: string;
  synced_at?: string;
  thread_id?: string;
  title?: string;
  updated_at?: string;
  updated_by_id?: string;
  version_id?: string;
  workspace_id?: string;
};

export function syncThreadsShape({
  token,
  workspaceId,
}: {
  token: string;
  workspaceId: string;
}): ShapeStreamOptions {
  const url = `${import.meta.env.VITE_ZYG_URL}/v1/sync/workspaces/${workspaceId}/shapes/parts/threads/`;
  return {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    params: {
      table: "thread",
    },
    url,
  };
}
