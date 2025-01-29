import {
  accountSchema,
  authMemberSchema,
  customerEventSchema,
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
import { z } from "zod";

export type Account = z.infer<typeof accountSchema>;

export type AuthMember = z.infer<typeof authMemberSchema>;

export type Workspace = z.infer<typeof workspaceSchema>;

export type Pat = z.infer<typeof patSchema>;
export type PatResponse = z.infer<typeof patSchema>;

export type Label = z.infer<typeof labelSchema>;
export type LabelResponse = z.infer<typeof labelSchema>;

export type WorkspaceMetrics = z.infer<typeof threadCountMetricsSchema>;
export type WorkspaceMetricsResponse = z.infer<typeof workspaceMetricsSchema>;

export type ThreadResponse = z.infer<typeof threadSchema>;

export type ThreadMessageResponse = z.infer<typeof threadMessageSchema>;

export type ThreadLabelResponse = z.infer<typeof threadLabelSchema>;

export type CustomerEventResponse = z.infer<typeof customerEventSchema>;

export type MessageAttachmentResponse = z.infer<typeof messageAttachmentSchema>;

export type PostmarkMailServerSetting = z.infer<
  typeof postmarkMailServerSettingSchema
>;

export type Pk = string;

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
