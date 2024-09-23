import { z } from "zod";

import {
  accountResponseSchema,
  authMemberResponseSchema,
  CustomerResponse,
  customerResponseSchema,
  LabelResponse,
  labelResponseSchema,
  MemberResponse,
  memberResponseSchema,
  PatResponse,
  patResponseSchema,
  ThreadResponse,
  workspaceResponseSchema,
} from "./schema";

export type Account = z.infer<typeof accountResponseSchema>;

export type Workspace = z.infer<typeof workspaceResponseSchema>;

export type AuthMember = z.infer<typeof authMemberResponseSchema>;

export type Pat = z.infer<typeof patResponseSchema>;

export type Customer = z.infer<typeof customerResponseSchema>;

export type Label = z.infer<typeof labelResponseSchema>;

export type Member = z.infer<typeof memberResponseSchema>;

export type LabelMetrics = {
  count: number;
  icon: string;
  labelId: string;
  name: string;
};

export type WorkspaceMetrics = {
  active: number;
  assignedToMe: number;
  done: number;
  labels: [] | LabelMetrics[];
  otherAssigned: number;
  snoozed: number;
  unassigned: number;
};

// Represents the thread model stored in local state, diff from the zod schema
export type Thread = {
  assigneeId: null | string;
  channel: string;
  createdAt: string;
  customerId: string;
  description: string;
  inboundCustomerId: null | string;
  inboundFirstSeqId: null | string;
  inboundLastSeqId: null | string;
  outboundFirstSeqId: null | string;
  outboundLastSeqId: null | string;
  outboundMemberId: null | string;
  previewText: string;
  priority: string;
  replied: boolean;
  stage: string;
  status: string;
  threadId: string;
  title: string;
  updatedAt: string;
};

export type ThreadChat = {
  body: string;
  chatId: string;
  createdAt: string;
  customerId: null | string;
  isHead: boolean;
  memberId: null | string;
  sequence: number;
  threadId: string;
  updatedAt: string;
};

export type Pk = string;

export function threadTransformer() {
  return {
    normalize(thread: ThreadResponse): [Pk, Thread] {
      const {
        assignee,
        customer,
        inboundCustomer,
        inboundFirstSeqId,
        inboundLastSeqId,
        outboundFirstSeqId,
        outboundLastSeqId,
        outboundMember,
        threadId,
        ...rest
      } = thread;
      const customerId = customer.customerId;
      const assigneeId = assignee?.memberId || null;
      const inboundCustomerId = inboundCustomer?.customerId || null;
      const outboundMemberId = outboundMember?.memberId || null;
      return [
        threadId,
        {
          threadId,
          ...rest,
          assigneeId: assigneeId,
          customerId: customerId,
          inboundCustomerId: inboundCustomerId || null,
          inboundFirstSeqId: inboundFirstSeqId || null,
          inboundLastSeqId: inboundLastSeqId || null,
          outboundFirstSeqId: outboundFirstSeqId || null,
          outboundLastSeqId: outboundLastSeqId || null,
          outboundMemberId: outboundMemberId || null,
        },
      ];
    },
  };
}

export function labelTransformer() {
  return {
    normalize(label: LabelResponse): [Pk, Label] {
      const { createdAt, icon, labelId, name, updatedAt } = label;
      return [
        labelId,
        {
          createdAt,
          icon,
          labelId,
          name,
          updatedAt,
        },
      ];
    },
  };
}

export function customerTransformer() {
  return {
    normalize(customer: CustomerResponse): [Pk, Customer] {
      const { customerId, ...rest } = customer;
      return [
        customerId,
        {
          customerId,
          ...rest,
        },
      ];
    },
  };
}

export function memberTransformer() {
  return {
    normalize(member: MemberResponse): [Pk, Member] {
      const { memberId, ...rest } = member;
      return [
        memberId,
        {
          memberId,
          ...rest,
        },
      ];
    },
  };
}

export function patTransformer() {
  return {
    normalize(pat: PatResponse): [Pk, Pat] {
      const { patId, ...rest } = pat;
      return [
        patId,
        {
          patId,
          ...rest,
        },
      ];
    },
  };
}
