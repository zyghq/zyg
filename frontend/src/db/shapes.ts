import { MemberShapeMap } from "@/db/store";
import { MemberRow, MemberRowUpdates } from "@/db/sync";

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

export type MemberShapeUpdates = {
  avatarUrl?: string;
  createdAt?: string;
  memberId?: string;
  name?: string;
  permissions?: Record<string, unknown>;
  publicName?: string;
  role?: string;
  syncedAt?: string;
  updatedAt?: string;
};

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

export function takeMemberUpdates(
  member: MemberRowUpdates,
): MemberShapeUpdates {
  const result: MemberShapeUpdates = {};

  if (member.avatar_url !== undefined) result.avatarUrl = member.avatar_url;
  if (member.created_at !== undefined) result.createdAt = member.created_at;
  if (member.member_id !== undefined) result.memberId = member.member_id;
  if (member.name !== undefined) result.name = member.name;
  if (member.permissions !== undefined) result.permissions = member.permissions;
  if (member.public_name !== undefined) result.publicName = member.public_name;
  if (member.role !== undefined) result.role = member.role;
  if (member.synced_at !== undefined) result.syncedAt = member.synced_at;
  if (member.updated_at !== undefined) result.updatedAt = member.updated_at;

  return result;
}
