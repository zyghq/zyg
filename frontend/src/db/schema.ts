import { z } from "zod";

export const accountResponseSchema = z.object({
  accountId: z.string(),
  email: z.string(),
  provider: z.string(),
  name: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const workspaceResponseSchema = z.object({
  workspaceId: z.string(),
  name: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const authMemberResponseSchema = z.object({
  memberId: z.string(),
  name: z.string(),
  role: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const patResponseSchema = z.object({
  accountId: z.string(),
  patId: z.string(),
  token: z.string(),
  name: z.string(),
  description: z.string().nullable().default(null),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const memberResponseSchema = z.object({
  memberId: z.string(),
  name: z.string(),
  role: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const workspaceMetricsResponseSchema = z.object({
  count: z.object({
    active: z.number().default(0),
    done: z.number().default(0),
    snoozed: z.number().default(0),
    assignedToMe: z.number().default(0),
    unassigned: z.number().default(0),
    otherAssigned: z.number().default(0),
    labels: z
      .array(
        z.object({
          labelId: z.string(),
          name: z.string().default(""),
          icon: z.string().default(""),
          count: z.number().default(0),
        })
      )
      .default([]),
  }),
});

export const customerResponseSchema = z.object({
  customerId: z.string(),
  externalId: z.string().nullable().default(null),
  email: z.string().nullable().default(null),
  phone: z.string().nullable().default(null),
  name: z.string(),
  isVerified: z.boolean(),
  role: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const threadResponseSchema = z.object({
  threadId: z.string(),
  customer: z.object({
    customerId: z.string(),
    name: z.string(),
  }),
  title: z.string(),
  description: z.string(),
  sequence: z.number(),
  status: z.string(),
  read: z.boolean(),
  replied: z.boolean(),
  priority: z.string(),
  spam: z.boolean(),
  channel: z.string(),
  previewText: z.string(),
  assignee: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  ingressFirstSeq: z.number().nullable().default(null),
  ingressLastSeq: z.number().nullable().default(null),
  ingressCustomer: z
    .object({
      customerId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  egressFirstSeq: z.number().nullable().default(null),
  egressLastSeq: z.number().nullable().default(null),
  egressMember: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const labelResponseSchema = z.object({
  labelId: z.string(),
  name: z.string(),
  icon: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const threadLabelResponseSchema = z.object({
  threadLabelId: z.string(),
  threadId: z.string(),
  labelId: z.string(),
  name: z.string(),
  icon: z.string(),
  addedBy: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const threadChatResponseSchema = z.object({
  threadId: z.string(),
  chatId: z.string(),
  body: z.string(),
  sequence: z.number(),
  customer: z
    .object({
      customerId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  member: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  isHead: z.boolean(),
  createdAt: z.string(),
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
