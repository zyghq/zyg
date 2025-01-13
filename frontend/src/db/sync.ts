import { ShapeStreamOptions } from "@electric-sql/client";

export type WorkspaceRow = {
  created_at: string;
  name: string;
  public_name: string;
  synced_at: string;
  updated_at: string;
  version_id: string;
  workspace_id: string;
};

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

export function syncWorkspaceShape({
  token,
  workspaceId,
}: {
  token: string;
  workspaceId: string;
}): ShapeStreamOptions {
  const url = `${import.meta.env.VITE_ZYG_URL}/v1/sync/workspaces/${workspaceId}/shapes/parts/workspace/`;
  return {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    params: {
      table: "workspace",
    },
    url,
  };
}

export function syncWorkspaceMemberShape({
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
