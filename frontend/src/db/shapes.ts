import { MemberShapeMap} from "@/db/store";
import { MemberRow, WorkspaceRow } from "@/db/sync";

export type WorkspaceShape = {
  createdAt: string;
  name: string;
  publicName: string;
  syncedAt: string;
  updatedAt: string;
  versionId: string;
  workspaceId: string;
};

export type MemberShape = {
  avatarUrl: string;
  createdAt: string;
  memberId: string;
  name: string;
  permissions: Record<string, unknown>;
  publicName: string;
  role: string;
  syncedAt: string;
  updatedAt: string;
  versionId: string;
  workspaceId: string;
};

export function workspaceRowToShape(row: WorkspaceRow): WorkspaceShape {
  return {
    createdAt: row.created_at,
    name: row.name,
    publicName: row.public_name,
    syncedAt: row.synced_at,
    updatedAt: row.updated_at,
    versionId: row.version_id,
    workspaceId: row.workspace_id,
  };
}

export function memberRowToShape(row: MemberRow): MemberShape {
  return {
    avatarUrl: row.avatar_url,
    createdAt: row.created_at,
    memberId: row.member_id,
    name: row.name,
    permissions: row.permissions,
    publicName: row.public_name,
    role: row.role,
    syncedAt: row.synced_at,
    updatedAt: row.updated_at,
    versionId: row.version_id,
    workspaceId: row.workspace_id,
  };
}

export function membersToMap(members: MemberShape[]): MemberShapeMap {
  return members.reduce((acc, member) => {
    acc[member.memberId] = member;
    return acc;
  }, {} as MemberShapeMap);
}