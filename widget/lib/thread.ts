import { z } from "zod";

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

export type ThreadResponse = z.infer<typeof threadResponseSchema>;

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

export type ThreadChatResponse = z.infer<typeof threadChatResponseSchema>;
