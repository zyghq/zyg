import { relations } from "drizzle-orm/relations";
import { workspace, member, customer, thread } from "./schema.ts";

export const memberRelations = relations(member, ({one, many}) => ({
	workspace: one(workspace, {
		fields: [member.workspaceId],
		references: [workspace.workspaceId]
	}),
	threads_assigneeId: many(thread, {
		relationName: "thread_assigneeId_member_memberId"
	}),
	threads_statusChangedById: many(thread, {
		relationName: "thread_statusChangedById_member_memberId"
	}),
	threads_createdById: many(thread, {
		relationName: "thread_createdById_member_memberId"
	}),
	threads_updatedById: many(thread, {
		relationName: "thread_updatedById_member_memberId"
	}),
}));

export const workspaceRelations = relations(workspace, ({many}) => ({
	members: many(member),
	customers: many(customer),
	threads: many(thread),
}));

export const customerRelations = relations(customer, ({one, many}) => ({
	workspace: one(workspace, {
		fields: [customer.workspaceId],
		references: [workspace.workspaceId]
	}),
	threads: many(thread),
}));

export const threadRelations = relations(thread, ({one}) => ({
	workspace: one(workspace, {
		fields: [thread.workspaceId],
		references: [workspace.workspaceId]
	}),
	customer: one(customer, {
		fields: [thread.customerId],
		references: [customer.customerId]
	}),
	member_assigneeId: one(member, {
		fields: [thread.assigneeId],
		references: [member.memberId],
		relationName: "thread_assigneeId_member_memberId"
	}),
	member_statusChangedById: one(member, {
		fields: [thread.statusChangedById],
		references: [member.memberId],
		relationName: "thread_statusChangedById_member_memberId"
	}),
	member_createdById: one(member, {
		fields: [thread.createdById],
		references: [member.memberId],
		relationName: "thread_createdById_member_memberId"
	}),
	member_updatedById: one(member, {
		fields: [thread.updatedById],
		references: [member.memberId],
		relationName: "thread_updatedById_member_memberId"
	}),
}));