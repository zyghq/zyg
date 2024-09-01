import { z } from "zod";

export const sdkInitResponseSchema = z.object({
  widgetId: z.string(),
  sessionId: z.string().optional(),
  customerExternalId: z.string().optional(),
  customerEmail: z.string().optional(),
  customerPhone: z.string().optional(),
  customerHash: z.string().optional(),
  traits: z.record(z.string()).optional(),
  isVerified: z.boolean().optional(),
});

export type SdkInitResponse = z.infer<typeof sdkInitResponseSchema>;
