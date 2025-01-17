import { z } from "zod";

export const accountSchema = z.object({
  accountId: z.string(),
  createdAt: z.string(),
  email: z.string(),
  name: z.string(),
  provider: z.string(),
  updatedAt: z.string(),
});

export const workspaceSchema = z.object({
  createdAt: z.string(),
  name: z.string(),
  updatedAt: z.string(),
  workspaceId: z.string(),
});

export const authMemberSchema = z.object({
  createdAt: z.string(),
  memberId: z.string(),
  name: z.string(),
  role: z.string(),
  updatedAt: z.string(),
});

export const patSchema = z.object({
  accountId: z.string(),
  createdAt: z.string(),
  description: z.string().nullable().default(null),
  name: z.string(),
  patId: z.string(),
  token: z.string(),
  updatedAt: z.string(),
});

export const threadLabelMetricsSchema = z.object({
  count: z.number().default(0),
  icon: z.string().default(""),
  labelId: z.string(),
  name: z.string().default(""),
});

export const threadCountMetricsSchema = z.object({
  active: z.number().default(0),
  assignedToMe: z.number().default(0),
  hold: z.number().default(0),
  labels: z.array(threadLabelMetricsSchema).default([]),
  needsFirstResponse: z.number().default(0),
  needsNextResponse: z.number().default(0),
  otherAssigned: z.number().default(0),
  snoozed: z.number().default(0),
  unassigned: z.number().default(0),
  waitingOnCustomer: z.number().default(0),
});

export const workspaceMetricsSchema = z.object({
  count: threadCountMetricsSchema,
});

export const threadSchema = z.object({
  assignee: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  channel: z.string(),
  createdAt: z.string(),
  customer: z.object({
    customerId: z.string(),
    name: z.string(),
  }),
  description: z.string(),
  inboundFirstSeqId: z.string().nullable().default(null),
  inboundLastSeqId: z.string().nullable().default(null),
  outboundFirstSeqId: z.string().nullable().default(null),
  outboundLastSeqId: z.string().nullable().default(null),
  previewText: z.string(),
  priority: z.string(),
  replied: z.boolean(),
  stage: z.string(),
  status: z.string(),
  statusChangedAt: z.string(),
  threadId: z.string(),
  title: z.string(),
  updatedAt: z.string(),
});

export const labelSchema = z.object({
  createdAt: z.string(),
  icon: z.string(),
  labelId: z.string(),
  name: z.string(),
  updatedAt: z.string(),
});

export const threadLabelSchema = z.object({
  addedBy: z.string(),
  createdAt: z.string(),
  icon: z.string(),
  labelId: z.string(),
  name: z.string(),
  threadId: z.string(),
  threadLabelId: z.string(),
  updatedAt: z.string(),
});

export const messageAttachmentSchema = z.object({
  attachmentId: z.string(),
  contentKey: z.string(),
  contentType: z.string(),
  contentUrl: z.string().default(""),
  createdAt: z.string(),
  error: z.string().default(""),
  hasError: z.boolean(),
  md5Hash: z.string(),
  messageId: z.string(),
  name: z.string(),
  spam: z.boolean(),
  updatedAt: z.string(),
});

export const threadMessageSchema = z.object({
  attachments: z.array(messageAttachmentSchema).default([]),
  channel: z.string(),
  createdAt: z.string(),
  customer: z
    .object({
      customerId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  htmlBody: z.string(),
  markdownBody: z.string(),
  member: z
    .object({
      memberId: z.string(),
      name: z.string(),
    })
    .nullable()
    .default(null),
  messageId: z.string(),
  textBody: z.string(),
  threadId: z.string(),
  updatedAt: z.string(),
});

// Component schemas
const ComponentText = z.object({
  componentText: z.object({
    text: z.string(),
    textColor: z.string(),
    textSize: z.string(),
  }),
});

const ComponentSpacer = z.object({
  componentSpacer: z.object({
    spacerSize: z.string(),
  }),
});

const ComponentLinkButton = z.object({
  componentLinkButton: z.object({
    linkButtonLabel: z.string(),
    linkButtonUrl: z.string().url(),
  }),
});

const ComponentDivider = z.object({
  componentDivider: z.object({
    dividerSize: z.string(),
  }),
});

const ComponentCopyButton = z.object({
  componentCopyButton: z.object({
    copyButtonToolTipLabel: z.string(),
    copyButtonValue: z.string(),
  }),
});

const ComponentBadge = z.object({
  componentBadge: z.object({
    badgeColor: z.string(),
    badgeLabel: z.string(),
  }),
});

// Row component needs to reference other component types
const ComponentRow = z.object({
  componentRow: z.object({
    rowAsideContent: z.array(
      z.union([
        ComponentBadge,
        ComponentText,
        ComponentSpacer,
        ComponentLinkButton,
        ComponentDivider,
        ComponentCopyButton,
      ]),
    ),
    rowMainContent: z.array(
      z.union([
        ComponentBadge,
        ComponentText,
        ComponentSpacer,
        ComponentLinkButton,
        ComponentDivider,
        ComponentCopyButton,
      ]),
    ),
  }),
});

// Union type for all possible components
const Component = z.union([
  ComponentText,
  ComponentSpacer,
  ComponentLinkButton,
  ComponentDivider,
  ComponentCopyButton,
  ComponentBadge,
  ComponentRow,
]);

// Main customer event schema
export const customerEventSchema = z.object({
  components: z.array(Component),
  createdAt: z.string(),
  eventId: z.string(),
  severity: z.string(),
  timestamp: z.string(),
  title: z.string(),
  updatedAt: z.string(),
});

export const postmarkMailServerSettingSchema = z.object({
  createdAt: z.string(),
  dkimHost: z.string().optional().nullable(),
  dkimTextValue: z.string().optional().nullable(),
  dkimUpdateStatus: z.string().optional().nullable(),
  dnsDomainId: z.number().int().optional().nullable(),
  dnsVerifiedAt: z.string().optional().nullable(),
  domain: z.string(),
  email: z.string().email(),
  hasDNS: z.boolean(),
  hasError: z.boolean(),
  hasForwardingEnabled: z.boolean(),
  inboundEmail: z.string().optional().nullable(),
  isDNSVerified: z.boolean(),
  isEnabled: z.boolean(),
  returnPathDomain: z.string().optional().nullable(),
  returnPathDomainCNAME: z.string().optional().nullable(), // the CNAME value
  returnPathDomainVerified: z.boolean(),
  serverId: z.number().int().positive(),
  serverToken: z.string(),
  updatedAt: z.string(),
  workspaceId: z.string(),
});
