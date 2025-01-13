import { z } from "zod";

import {
  accountResponseSchema,
  authMemberResponseSchema,
  CustomerResponse,
  customerResponseSchema,
  LabelResponse,
  labelResponseSchema,
  PatResponse,
  patResponseSchema,
  threadCountMetricsSchema,
  threadLabelMetricsSchema,
  ThreadResponse,
} from "./schema";

export type Account = z.infer<typeof accountResponseSchema>;

export type AuthMember = z.infer<typeof authMemberResponseSchema>;

export type Pat = z.infer<typeof patResponseSchema>;

export type Customer = z.infer<typeof customerResponseSchema>;

export type Label = z.infer<typeof labelResponseSchema>;

export type WorkspaceMetrics = z.infer<typeof threadCountMetricsSchema>;

export type LabelMetrics = z.infer<typeof threadLabelMetricsSchema>;

// Represents the thread model stored in local state, diff from the zod schema
export type Thread = {
  assigneeId: null | string;
  channel: string;
  createdAt: string;
  customerId: string;
  description: string;
  inboundFirstSeqId: null | string;
  inboundLastSeqId: null | string;
  outboundFirstSeqId: null | string;
  outboundLastSeqId: null | string;
  previewText: string;
  priority: string;
  replied: boolean;
  stage: string;
  status: string;
  threadId: string;
  title: string;
  updatedAt: string;
};

export type Pk = string;

export function threadTransformer() {
  return {
    normalize(thread: ThreadResponse): [Pk, Thread] {
      const {
        assignee,
        customer,
        inboundFirstSeqId,
        inboundLastSeqId,
        outboundFirstSeqId,
        outboundLastSeqId,
        threadId,
        ...rest
      } = thread;
      const customerId = customer.customerId;
      const assigneeId = assignee?.memberId || null;
      return [
        threadId,
        {
          threadId,
          ...rest,
          assigneeId: assigneeId,
          customerId: customerId,
          inboundFirstSeqId: inboundFirstSeqId || null,
          inboundLastSeqId: inboundLastSeqId || null,
          outboundFirstSeqId: outboundFirstSeqId || null,
          outboundLastSeqId: outboundLastSeqId || null,
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

// export function memberTransformer() {
//   return {
//     normalize(member: MemberResponse): [Pk, Member] {
//       const { memberId, ...rest } = member;
//       return [
//         memberId,
//         {
//           memberId,
//           ...rest,
//         },
//       ];
//     },
//   };
// }

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
