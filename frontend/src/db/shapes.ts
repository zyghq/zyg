import { CustomerShapeMap, MemberShapeMap } from "@/db/store";
import {
  CustomerRow,
  CustomerRowUpdates,
  MemberRow,
  MemberRowUpdates,
} from "@/db/sync";

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

export type CustomerShape = {
  avatarUrl: string;
  createdAt: string;
  customerId: string;
  email: string;
  externalId: string;
  isEmailVerified: boolean;
  name: string;
  phone: string;
  role: string;
  syncedAt: string;
  updatedAt: string;
  versionId: string;
  workspaceId: string;
};

export type CustomerShapeUpdates = {
  avatarUrl?: string;
  createdAt?: string;
  customerId?: string;
  email?: string;
  externalId?: string;
  isEmailVerified?: boolean;
  name?: string;
  phone?: string;
  role?: string;
  syncedAt?: string;
  updatedAt?: string;
};

export function customerRowToShape(row: CustomerRow): CustomerShape {
  return {
    avatarUrl: row.avatar_url,
    createdAt: row.created_at,
    customerId: row.customer_id,
    email: row.email,
    externalId: row.external_id,
    isEmailVerified: row.is_email_verified,
    name: row.name,
    phone: row.phone,
    role: row.role,
    syncedAt: row.synced_at,
    updatedAt: row.updated_at,
    versionId: row.version_id,
    workspaceId: row.workspace_id,
  };
}

export function customersToMap(customers: CustomerShape[]): CustomerShapeMap {
  return customers.reduce((acc, customer) => {
    const { customerId, ...rest } = customer;

    acc[customerId] = {
      customerId,
      ...rest,
    };

    return acc;
  }, {} as CustomerShapeMap);
}

export function takeCustomerUpdates(
  customer: CustomerRowUpdates,
): CustomerShapeUpdates {
  const result: CustomerShapeUpdates = {};
  if (customer.avatar_url !== undefined) result.avatarUrl = customer.avatar_url;
  if (customer.created_at !== undefined) result.createdAt = customer.created_at;
  if (customer.customer_id !== undefined)
    result.customerId = customer.customer_id;
  if (customer.email !== undefined) result.email = customer.email;
  if (customer.external_id !== undefined)
    result.externalId = customer.external_id;
  if (customer.is_email_verified !== undefined)
    result.isEmailVerified = customer.is_email_verified;
  if (customer.name !== undefined) result.name = customer.name;
  if (customer.phone !== undefined) result.phone = customer.phone;
  if (customer.role !== undefined) result.role = customer.role;
  if (customer.synced_at !== undefined) result.syncedAt = customer.synced_at;
  if (customer.updated_at !== undefined) result.updatedAt = customer.updated_at;

  return result;
}
