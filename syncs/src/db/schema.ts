import {
  pgTable,
  varchar,
  timestamp,
  uuid,
  foreignKey,
  jsonb,
  unique,
  boolean,
  text,
} from "drizzle-orm/pg-core";
import { sql } from "drizzle-orm";

export const workspace = pgTable("workspace", {
  workspaceId: varchar("workspace_id", { length: 255 }).primaryKey().notNull(),
  name: varchar({ length: 255 }).notNull(),
  publicName: varchar("public_name", { length: 255 }).notNull(),
  createdAt: timestamp("created_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  updatedAt: timestamp("updated_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
  versionId: uuid("version_id").notNull(),
  syncedAt: timestamp("synced_at", { mode: "string" }).default(
    sql`CURRENT_TIMESTAMP`,
  ),
});

export const member = pgTable(
  "member",
  {
    memberId: varchar("member_id", { length: 255 }).primaryKey().notNull(),
    workspaceId: varchar("workspace_id", { length: 255 }).notNull(),
    name: varchar({ length: 255 }).notNull(),
    publicName: varchar("public_name", { length: 255 }).notNull(),
    role: varchar({ length: 255 }).notNull(),
    permissions: jsonb().notNull(),
    avatarUrl: varchar("avatar_url", { length: 511 }).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    versionId: uuid("version_id").notNull(),
    syncedAt: timestamp("synced_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
  },
  (table) => [
    foreignKey({
      columns: [table.workspaceId],
      foreignColumns: [workspace.workspaceId],
      name: "member_workspace_id_fkey",
    }),
  ],
);

export const customer = pgTable(
  "customer",
  {
    customerId: varchar("customer_id", { length: 255 }).primaryKey().notNull(),
    workspaceId: varchar("workspace_id", { length: 255 }).notNull(),
    externalId: varchar("external_id", { length: 255 }),
    email: varchar({ length: 255 }),
    phone: varchar({ length: 255 }),
    name: varchar({ length: 255 }).notNull(),
    role: varchar({ length: 255 }).notNull(),
    avatarUrl: varchar("avatar_url", { length: 511 }).notNull(),
    isEmailVerified: boolean("is_email_verified").default(false).notNull(),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    versionId: uuid("version_id").notNull(),
    syncedAt: timestamp("synced_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
  },
  (table) => [
    foreignKey({
      columns: [table.workspaceId],
      foreignColumns: [workspace.workspaceId],
      name: "customer_workspace_id_fkey",
    }),
    unique("customer_workspace_id_external_id_key").on(
      table.workspaceId,
      table.externalId,
    ),
    unique("customer_workspace_id_email_key").on(
      table.workspaceId,
      table.email,
    ),
    unique("customer_workspace_id_phone_key").on(
      table.workspaceId,
      table.phone,
    ),
  ],
);

export const thread = pgTable(
  "thread",
  {
    threadId: varchar("thread_id", { length: 255 }).primaryKey().notNull(),
    workspaceId: varchar("workspace_id", { length: 255 }).notNull(),
    customerId: varchar("customer_id", { length: 255 }).notNull(),
    assigneeId: varchar("assignee_id", { length: 255 }),
    assignedAt: timestamp("assigned_at", { mode: "string" }),
    title: text().notNull(),
    description: text().notNull(),
    previewText: text("preview_text").notNull(),
    status: varchar({ length: 127 }).notNull(),
    statusChangedAt: timestamp("status_changed_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    statusChangedById: varchar("status_changed_by_id", {
      length: 255,
    }).notNull(),
    stage: varchar({ length: 127 }).notNull(),
    replied: boolean().default(false).notNull(),
    priority: varchar({ length: 255 }).notNull(),
    channel: varchar({ length: 127 }).notNull(),
    createdById: varchar("created_by_id", { length: 255 }).notNull(),
    updatedById: varchar("updated_by_id", { length: 255 }).notNull(),
    labels: jsonb(),
    inboundSeqId: varchar("inbound_seq_id", { length: 255 }),
    outboundSeqId: varchar("outbound_seq_id", { length: 255 }),
    createdAt: timestamp("created_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    updatedAt: timestamp("updated_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
    versionId: uuid("version_id").notNull(),
    syncedAt: timestamp("synced_at", { mode: "string" }).default(
      sql`CURRENT_TIMESTAMP`,
    ),
  },
  (table) => [
    foreignKey({
      columns: [table.workspaceId],
      foreignColumns: [workspace.workspaceId],
      name: "thread_workspace_id_fkey",
    }),
    foreignKey({
      columns: [table.customerId],
      foreignColumns: [customer.customerId],
      name: "thread_customer_id_fkey",
    }),
    foreignKey({
      columns: [table.assigneeId],
      foreignColumns: [member.memberId],
      name: "thread_assignee_id_fkey",
    }),
    foreignKey({
      columns: [table.statusChangedById],
      foreignColumns: [member.memberId],
      name: "thread_status_changed_by_id_fkey",
    }),
    foreignKey({
      columns: [table.createdById],
      foreignColumns: [member.memberId],
      name: "thread_created_by_id_fkey",
    }),
    foreignKey({
      columns: [table.updatedById],
      foreignColumns: [member.memberId],
      name: "thread_updated_by_id_fkey",
    }),
  ],
);
