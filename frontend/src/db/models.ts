import {
  accountSchema,
  authMemberSchema,
  customerEventSchema,
  customerSchema,
  labelSchema,
  messageAttachmentSchema,
  patSchema,
  postmarkMailServerSettingSchema,
  threadCountMetricsSchema,
  threadLabelSchema,
  threadMessageSchema,
  threadSchema,
  workspaceMetricsSchema,
  workspaceSchema,
} from "@/db/schema";
import { CustomerMap } from "@/db/store";
import { z } from "zod";

export type Account = z.infer<typeof accountSchema>;

export type AuthMember = z.infer<typeof authMemberSchema>;

export type Workspace = z.infer<typeof workspaceSchema>;

export type Pat = z.infer<typeof patSchema>;
export type PatResponse = z.infer<typeof patSchema>;

export type Customer = z.infer<typeof customerSchema>;
export type CustomerResponse = z.infer<typeof customerSchema>;

export type Label = z.infer<typeof labelSchema>;
export type LabelResponse = z.infer<typeof labelSchema>;

export type WorkspaceMetrics = z.infer<typeof threadCountMetricsSchema>;
export type WorkspaceMetricsResponse = z.infer<typeof workspaceMetricsSchema>;

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
export type ThreadResponse = z.infer<typeof threadSchema>;

export type ThreadMessageResponse = z.infer<typeof threadMessageSchema>;

export type ThreadLabelResponse = z.infer<typeof threadLabelSchema>;

export type CustomerEventResponse = z.infer<typeof customerEventSchema>;

export type MessageAttachmentResponse = z.infer<typeof messageAttachmentSchema>;

export type PostmarkMailServerSetting = z.infer<
  typeof postmarkMailServerSettingSchema
>;

export type Pk = string;

// Deprecated
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

export function threadsToMap(
  threads: ThreadResponse[],
): Record<string, Thread> {
  return threads.reduce(
    (acc, thread) => {
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

      acc[threadId] = {
        threadId,
        ...rest,
        assigneeId,
        customerId,
        inboundFirstSeqId: inboundFirstSeqId || null,
        inboundLastSeqId: inboundLastSeqId || null,
        outboundFirstSeqId: outboundFirstSeqId || null,
        outboundLastSeqId: outboundLastSeqId || null,
      };

      return acc;
    },
    {} as Record<string, Thread>,
  );
}







export function labelsToMap(labels: LabelResponse[]): Record<string, Label> {
  return labels.reduce(
    (acc, label) => {
      const { createdAt, icon, labelId, name, updatedAt } = label;

      acc[labelId] = {
        createdAt,
        icon,
        labelId,
        name,
        updatedAt,
      };

      return acc;
    },
    {} as Record<string, Label>,
  );
}

export function customersToMap(
  customers: CustomerResponse[],
): CustomerMap {
  return customers.reduce(
    (acc, customer) => {
      const { customerId, ...rest } = customer;

      acc[customerId] = {
        customerId,
        ...rest,
      };

      return acc;
    },
    {} as CustomerMap,
  );
}

export function patsToMap(pats: PatResponse[]): Record<string, Pat> {
  return pats.reduce(
    (acc, pat) => {
      const { patId, ...rest } = pat;

      acc[patId] = {
        patId,
        ...rest,
      };

      return acc;
    },
    {} as Record<string, Pat>,
  );
}
