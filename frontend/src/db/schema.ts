import { z } from "zod";

export const accountResponseSchema = z.object({
  accountId: z.string(),
  createdAt: z.string(),
  email: z.string(),
  name: z.string(),
  provider: z.string(),
  updatedAt: z.string(),
});

export const workspaceResponseSchema = z.object({
  createdAt: z.string(),
  name: z.string(),
  updatedAt: z.string(),
  workspaceId: z.string(),
});

export const authMemberResponseSchema = z.object({
  createdAt: z.string(),
  memberId: z.string(),
  name: z.string(),
  role: z.string(),
  updatedAt: z.string(),
});

export const patResponseSchema = z.object({
  accountId: z.string(),
  createdAt: z.string(),
  description: z.string().nullable().default(null),
  name: z.string(),
  patId: z.string(),
  token: z.string(),
  updatedAt: z.string(),
});

export const memberResponseSchema = z.object({
  createdAt: z.string(),
  memberId: z.string(),
  name: z.string(),
  role: z.string(),
  updatedAt: z.string(),
});

export const workspaceMetricsResponseSchema = z.object({
  count: z.object({
    active: z.number().default(0),
    assignedToMe: z.number().default(0),
    done: z.number().default(0),
    labels: z
      .array(
        z.object({
          count: z.number().default(0),
          icon: z.string().default(""),
          labelId: z.string(),
          name: z.string().default(""),
        })
      )
      .default([]),
    otherAssigned: z.number().default(0),
    snoozed: z.number().default(0),
    unassigned: z.number().default(0),
  }),
});

export const customerResponseSchema = z.object({
  createdAt: z.string(),
  customerId: z.string(),
  email: z.string().nullable().default(null),
  externalId: z.string().nullable().default(null),
  isVerified: z.boolean(),
  name: z.string(),
  phone: z.string().nullable().default(null),
  role: z.string(),
  updatedAt: z.string(),
});

export const threadResponseSchema = z.object({
  assignee: z
    .object({
      memberId: z.string(),
      name: z.string(), // TODO: add support for avatarUrl
    })
    .nullable()
    .default(null),
  channel: z.string(),
  createdAt: z.string(),
  customer: z.object({
    customerId: z.string(),
    name: z.string(), // TODO: add support for avatarUrl
  }),
  description: z.string(),
  inboundCustomer: z
    .object({
      customerId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  inboundFirstSeqId: z.string().nullable().default(null),
  inboundLastSeqId: z.string().nullable().default(null),
  outboundFirstSeqId: z.string().nullable().default(null),
  outboundLastSeqId: z.string().nullable().default(null),
  outboundMember: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  previewText: z.string(),
  priority: z.string(),
  replied: z.boolean(),
  stage: z.string(),
  status: z.string(),
  statusChangedAt: z.string(),
  threadId: z.string(),
  title: z.string(),
  updatedAt: z.string(),
});

export const labelResponseSchema = z.object({
  createdAt: z.string(),
  icon: z.string(),
  labelId: z.string(),
  name: z.string(),
  updatedAt: z.string(),
});

export const threadLabelResponseSchema = z.object({
  addedBy: z.string(),
  createdAt: z.string(),
  icon: z.string(),
  labelId: z.string(),
  name: z.string(),
  threadId: z.string(),
  threadLabelId: z.string(),
  updatedAt: z.string(),
});

export const threadChatResponseSchema = z.object({
  body: z.string(),
  chatId: z.string(),
  createdAt: z.string(),
  customer: z
    .object({
      customerId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  isHead: z.boolean(),
  member: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  sequence: z.number(),
  threadId: z.string(),
  updatedAt: z.string(),
});

export type ThreadResponse = z.infer<typeof threadResponseSchema>;

export type ThreadChatResponse = z.infer<typeof threadChatResponseSchema>;

export type LabelResponse = z.infer<typeof labelResponseSchema>;
export type ThreadLabelResponse = z.infer<typeof threadLabelResponseSchema>;

export type MemberResponse = z.infer<typeof memberResponseSchema>;

export type WorkspaceMetricsResponse = z.infer<
  typeof workspaceMetricsResponseSchema
>;

export type CustomerResponse = z.infer<typeof customerResponseSchema>;

export type PatResponse = z.infer<typeof patResponseSchema>;
