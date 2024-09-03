import { z } from "zod";

// 1:1 mapping with `WidgetInitPayload`
export const sdkInitResponseSchema = z.object({
  widgetId: z.string(),
  sessionId: z.string().optional(),
  externalId: z.string().optional(),
  email: z.string().optional(),
  phone: z.string().optional(),
  customerHash: z.string().optional(),
  isVerified: z.boolean().optional(),
  traits: z.record(z.string()).optional(),
});

export type SdkInitResponse = z.infer<typeof sdkInitResponseSchema>;

export interface HomeLink {
  id: string;
  title: string;
  href: string;
  previewText?: string;
}

export interface WidgetLayout {
  title: string;
  ctaSearchButtonText: string;
  ctaMessageButtonText: string;
  tabs: string[];
  defaultTab: string;
  homeLinks: HomeLink[];
}
