CREATE TABLE workspace
(
    "workspaceId" VARCHAR(255) NOT NULL, -- primary key
    "name"        VARCHAR(255) NOT NULL,
    "publicName"  VARCHAR(255) NOT NULL,
    "createdAt"   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updatedAt"   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    "versionId"   UUID         NOT NULL,
    "syncedAt"    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT workspace_workspace_id_pkey PRIMARY KEY ("workspaceId")
);

CREATE TABLE member
(
    "memberId"    VARCHAR(255) NOT NULL, -- Unique identifier for the member
    "workspaceId" VARCHAR(255) NOT NULL, -- Reference to workspace this member belongs to
    name          VARCHAR(255) NOT NULL, -- Display name of the member
    "publicName"  VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    role          VARCHAR(255) NOT NULL, -- Member's role in the workspace
    permissions   JSONB        NOT NULL,
    "avatarUrl"   VARCHAR(511) NOT NULL,
    "createdAt"   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updatedAt"   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    "versionId"   UUID         NOT NULL,
    "syncedAt"    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT member_member_id_pkey PRIMARY KEY ("memberId"),
    CONSTRAINT member_workspace_id_fkey FOREIGN KEY ("workspaceId") REFERENCES workspace ("workspaceId")
);


CREATE TABLE customer
(
    "customerId"      VARCHAR(255) NOT NULL,               -- primary key
    "workspaceId"     VARCHAR(255) NOT NULL,               -- fk to workspace
    "externalId"      VARCHAR(255) NULL,                   -- external id of the customer (optional identifier)
    email             VARCHAR(255) NULL,                   -- email address of the customer (optional identifier)
    phone             VARCHAR(255) NULL,                   -- phone number of the customer (optional identifier)
    name              VARCHAR(255) NOT NULL,               -- display name of the customer
    role              VARCHAR(255) NOT NULL,               -- role/type of the customer
    "avatarUrl"       VARCHAR(511) NOT NULL,
    "isEmailVerified" BOOLEAN      NOT NULL DEFAULT FALSE, -- whether email has been verified
    "createdAt"       TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updatedAt"       TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    "versionId"       UUID         NOT NULL,
    "syncedAt"        TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT customer_customer_id_pkey PRIMARY KEY ("customerId"),
    CONSTRAINT customer_workspace_id_fkey FOREIGN KEY ("workspaceId") REFERENCES workspace ("workspaceId"),
    CONSTRAINT customer_workspace_id_external_id_key UNIQUE ("workspaceId", "externalId"),
    CONSTRAINT customer_workspace_id_email_key UNIQUE ("workspaceId", email),
    CONSTRAINT customer_workspace_id_phone_key UNIQUE ("workspaceId", phone)
);

CREATE TABLE thread
(
    "threadId"          VARCHAR(255) NOT NULL,                           -- Unique identifier for the thread
    "workspaceId"       VARCHAR(255) NOT NULL,                           -- Workspace this thread belongs to
    "customerId"        VARCHAR(255) NOT NULL,                           -- Customer associated with this thread
    "assigneeId"        VARCHAR(255) NULL,                               -- Member assigned to handle this thread
    "assignedAt"        TIMESTAMP    NULL,                               -- When the thread was assigned
    title               TEXT         NOT NULL,                           -- Thread title
    description         TEXT         NOT NULL,                           -- Thread description
    "previewText"       TEXT         NOT NULL,
    status              VARCHAR(127) NOT NULL,                           -- Current status of the thread
    "statusChangedAt"   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP, -- When status was last changed
    "statusChangedById" VARCHAR(255) NOT NULL,                           -- Member who last changed the status
    stage               VARCHAR(127) NOT NULL,                           -- Current stage in the workflow
    replied             BOOLEAN      NOT NULL DEFAULT FALSE,             -- Whether a member has replied
    priority            VARCHAR(255) NOT NULL,                           -- Thread priority level
    channel             VARCHAR(127) NOT NULL,                           -- Communication channel used
    "createdById"       VARCHAR(255) NOT NULL,                           -- Member who created the thread
    "updatedById"       VARCHAR(255) NOT NULL,                           -- Member who last updated the thread
    labels              JSONB        NULL,
    "createdAt"         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    "updatedAt"         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    "versionId"         UUID         NOT NULL,
    "syncedAt"          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT thread_threadId_pkey PRIMARY KEY ("threadId"),
    CONSTRAINT thread_workspace_id_fkey FOREIGN KEY ("workspaceId") REFERENCES workspace ("workspaceId"),
    CONSTRAINT thread_customer_id_fkey FOREIGN KEY ("customerId") REFERENCES customer ("customerId"),
    CONSTRAINT thread_assignee_id_fkey FOREIGN KEY ("assigneeId") REFERENCES member ("memberId"),
    CONSTRAINT thread_status_changed_by_id_fkey FOREIGN KEY ("statusChangedById") REFERENCES member ("memberId"),
    CONSTRAINT thread_created_by_id_fkey FOREIGN KEY ("createdById") REFERENCES member ("memberId"),
    CONSTRAINT thread_updated_by_id_fkey FOREIGN KEY ("updatedById") REFERENCES member ("memberId")
);
