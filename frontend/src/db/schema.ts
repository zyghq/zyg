import { z } from "zod";

// represents the shape of the workspace object
// as per API response.
export const workspaceResponseSchema = z.object({
  workspaceId: z.string(),
  accountId: z.string(),
  name: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

// represents the shape of the workspace metrics object
// as per API response.
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

// represents the shape of the thread chat message object
// as per API response.
/**
 * Schema are subject to change based on the API response
 *
 * a thread chat cannot exists without a customer and a message
 * it can be assigned or unassigned
 * a message is either sent by customer or member, cannot be both as
 * both cannot send the same message
 *
 * `threadChatMessageId` represents the PK of the message
 * `threadChatId` represents the FK of the thread chat
 */
export const threadChatMessagePreviewSchema = z.object({
  threadChatId: z.string(),
  threadChatMessageId: z.string(),
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
  createdAt: z.string(),
  updatedAt: z.string(),
});

// represents the shape of the thread chat object
// as per API response.
export const threadChatResponseSchema = z.object({
  threadChatId: z.string(),
  sequence: z.number(),
  status: z.string(),
  read: z.boolean(),
  replied: z.boolean(),
  customer: z.object({
    customerId: z.string(),
    name: z.string(),
  }),
  assignee: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  createdAt: z.string(),
  updatedAt: z.string(),
  messages: z.array(threadChatMessagePreviewSchema).default([]),
});

export const accountResponseSchema = z.object({
  accountId: z.string(),
  email: z.string(),
  provider: z.string(),
  name: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

// @sanchitrk: membership can have more data like perms, etc.
export const membershipResponseSchema = z.object({
  workspaceId: z.string(),
  accountId: z.string(),
  memberId: z.string(),
  name: z.string().nullable().default(null),
  role: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const workspaceCustomerResponseSchema = z.object({
  workspaceId: z.string(),
  customerId: z.string(),
  externalId: z.string().nullable().default(null),
  email: z.string().nullable().default(null),
  phone: z.string().nullable().default(null),
  name: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const workspaceCustomersResponseSchema = z
  .array(workspaceCustomerResponseSchema)
  .default([]);

export const workspaceLabelResponseSchema = z.object({
  workspaceId: z.string(),
  labelId: z.string(),
  name: z.string(),
  icon: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const workspaceMemberResponseSchema = z.object({
  workspaceId: z.string(),
  accountId: z.string(),
  memberId: z.string(),
  name: z.string().nullable().default(null),
  role: z.string(),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const accountPatSchema = z.object({
  accountId: z.string(),
  patId: z.string(),
  token: z.string(),
  name: z.string(),
  description: z.string().nullable().default(null),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const threadChatMessageSchema = z.object({
  threadChatId: z.string(),
  threadChatMessageId: z.string(),
  body: z.string(),
  sequence: z.number(),
  member: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  customer: z
    .object({
      customerId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  createdAt: z.string(),
  updatedAt: z.string(),
});

export const threadChatMessagesResponseSchema = z.object({
  threadChatId: z.string(),
  sequence: z.number(),
  status: z.string(),
  read: z.boolean(),
  replied: z.boolean(),
  customer: z.object({
    customerId: z.string(),
    name: z.string(),
  }),
  assignee: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  createdAt: z.string(),
  updatedAt: z.string(),
  messages: z.array(threadChatMessageSchema).default([]),
});
