import { z } from "zod";
import {
  accountResponseSchema,
  workspaceResponseSchema,
  authMemberResponseSchema,
  memberResponseSchema,
  patResponseSchema,
  customerResponseSchema,
  labelResponseSchema,
  ThreadResponse,
  LabelResponse,
  CustomerResponse,
  MemberResponse,
  PatResponse,
} from "./schema";

export type Account = z.infer<typeof accountResponseSchema>;

export type Workspace = z.infer<typeof workspaceResponseSchema>;

export type AuthMember = z.infer<typeof authMemberResponseSchema>;

export type Pat = z.infer<typeof patResponseSchema>;

export type Customer = z.infer<typeof customerResponseSchema>;

export type Label = z.infer<typeof labelResponseSchema>;

export type Member = z.infer<typeof memberResponseSchema>;

export type LabelMetrics = {
  labelId: string;
  name: string;
  icon: string;
  count: number;
};

export type WorkspaceMetrics = {
  active: number;
  done: number;
  snoozed: number;
  assignedToMe: number;
  unassigned: number;
  otherAssigned: number;
  labels: LabelMetrics[] | [];
};

export type Thread = {
  threadId: string;
  customerId: string;
  title: string;
  description: string;
  sequence: number;
  status: string;
  read: boolean;
  replied: boolean;
  priority: string;
  spam: boolean;
  channel: string;
  previewText: string;
  assigneeId: string | null;
  inboundFirstSeqId: string | null;
  inboundLastSeqId: string | null;
  inboundCustomerId: string | null;
  outboundFirstSeqId: string | null;
  outboundLastSeqId: string | null;
  outboundMemberId: string | null;
  createdAt: string;
  updatedAt: string;
};

export type ThreadChat = {
  threadId: string;
  chatId: string;
  body: string;
  sequence: number;
  customerId: string | null;
  memberId: string | null;
  isHead: boolean;
  createdAt: string;
  updatedAt: string;
};

export type Pk = string;

export function threadTransformer() {
  return {
    normalize(thread: ThreadResponse): [Pk, Thread] {
      const {
        threadId,
        customer,
        assignee,
        inboundFirstSeqId,
        inboundLastSeqId,
        inboundCustomer,
        outboundMember,
        outboundFirstSeqId,
        outboundLastSeqId,
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
          customerId: customerId,
          assigneeId: assigneeId,
          inboundFirstSeqId: inboundFirstSeqId || null,
          inboundLastSeqId: inboundLastSeqId || null,
          inboundCustomerId: inboundCustomerId || null,
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
      const { labelId, name, icon, createdAt, updatedAt } = label;
      return [
        labelId,
        {
          labelId,
          name,
          icon,
          createdAt,
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
