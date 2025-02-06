import { CustomerShapeMap, MemberShapeMap, ThreadShapeMap } from "@/db/store";
import {
  CustomerRow,
  CustomerRowUpdates,
  MemberRow,
  MemberRowUpdates,
  ThreadRow,
  ThreadRowUpdates,
} from "@/db/sync";

// Represents the Member shape from sync engine
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

// Represents the Member shape updates from sync engine
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

// Maps sync engine `member` table row to Member shape
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

// Represents the Customer shape from sync engine
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

// Represents the Customer shape updates from sync engine
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

// Maps sync engine `customer` table row to Customer shape
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

// Represents the label assigned to Thread shape
export type ThreadLabelShape = {
  createdAt: string;
  labelId: string;
  name: string;
  updatedAt: string;
};
export type ThreadLabelShapeMap = Record<string, ThreadLabelShape>;

// Represents the Thread shape from sync engine
export type ThreadShape = {
  assignedAt: null | string;
  assigneeId: null | string;
  channel: string;
  createdAt: string;
  createdById: string;
  customerId: string;
  description: string;
  lastInboundAt: string | null;
  lastOutboundAt: string | null;
  labels: null | ThreadLabelShapeMap;
  previewText: string;
  priority: string;
  replied: boolean;
  stage: string;
  status: string;
  statusChangedAt: string;
  statusChangedById: string;
  syncedAt: string;
  threadId: string;
  title: string;
  updatedAt: string;
  updatedById: string;
  versionId: string;
  workspaceId: string;
};

export type ThreadShapeUpdates = {
  assignedAt?: null | string;
  assigneeId?: null | string;
  channel?: string;
  createdAt?: string;
  createdById?: string;
  customerId?: string;
  description?: string;
  lastInboundAt?: string | null;
  lastOutboundAt?: string | null;
  labels?: null | ThreadLabelShapeMap;
  previewText?: string;
  priority?: string;
  replied?: boolean;
  stage?: string;
  status?: string;
  statusChangedAt?: string;
  statusChangedById?: string;
  syncedAt?: string;
  threadId?: string;
  title?: string;
  updatedAt?: string;
  updatedById?: string;
  versionId?: string;
  workspaceId?: string;
};

// Maps sync engine `thread` table row to Thread shape
export function threadRowToShape(row: ThreadRow): ThreadShape {
  return {
    assignedAt: row.assigned_at,
    assigneeId: row.assignee_id,
    channel: row.channel,
    createdAt: row.created_at,
    createdById: row.created_by_id,
    customerId: row.customer_id,
    description: row.description,
    lastInboundAt: row.last_inbound_at,
    lastOutboundAt: row.last_outbound_at,
    labels: row.labels,
    previewText: row.preview_text,
    priority: row.priority,
    replied: row.replied,
    stage: row.stage,
    status: row.status,
    statusChangedAt: row.status_changed_at,
    statusChangedById: row.status_changed_by_id,
    syncedAt: row.synced_at,
    threadId: row.thread_id,
    title: row.title,
    updatedAt: row.updated_at,
    updatedById: row.updated_by_id,
    versionId: row.version_id,
    workspaceId: row.workspace_id,
  };
}

export function threadsToMap(threads: ThreadShape[]): ThreadShapeMap {
  return new Map(threads.map((thread) => [thread.threadId, thread]));
}

export function takeThreadUpdates(thread: ThreadRowUpdates): ThreadRowUpdates {
  const result: ThreadShapeUpdates = {};
  if (thread.assigned_at !== undefined) result.assignedAt = thread.assigned_at;
  if (thread.assignee_id !== undefined) result.assigneeId = thread.assignee_id;
  if (thread.channel !== undefined) result.channel = thread.channel;
  if (thread.created_at !== undefined) result.createdAt = thread.created_at;
  if (thread.created_by_id !== undefined)
    result.createdById = thread.created_by_id;
  if (thread.customer_id !== undefined) result.customerId = thread.customer_id;
  if (thread.description !== undefined) result.description = thread.description;
  if (thread.last_inbound_at !== undefined)
    result.lastInboundAt = thread.last_inbound_at;
  if (thread.labels !== undefined) result.labels = thread.labels;
  if (thread.last_outbound_at !== undefined)
    result.lastOutboundAt = thread.last_outbound_at;
  if (thread.preview_text !== undefined)
    result.previewText = thread.preview_text;
  if (thread.priority !== undefined) result.priority = thread.priority;
  if (thread.replied !== undefined) result.replied = thread.replied;
  if (thread.stage !== undefined) result.stage = thread.stage;
  if (thread.status !== undefined) result.status = thread.status;
  if (thread.status_changed_at !== undefined)
    result.statusChangedAt = thread.status_changed_at;
  if (thread.status_changed_by_id !== undefined)
    result.statusChangedById = thread.status_changed_by_id;
  if (thread.synced_at !== undefined) result.syncedAt = thread.synced_at;
  if (thread.thread_id !== undefined) result.threadId = thread.thread_id;
  if (thread.title !== undefined) result.title = thread.title;
  if (thread.updated_at !== undefined) result.updatedAt = thread.updated_at;
  if (thread.updated_by_id !== undefined)
    result.updatedById = thread.updated_by_id;
  if (thread.version_id !== undefined) result.versionId = thread.version_id;
  if (thread.workspace_id !== undefined)
    result.workspaceId = thread.workspace_id;

  return result;
}
