CREATE TABLE workspace
(
    workspace_id VARCHAR(255) NOT NULL, -- primary key
    name         VARCHAR(255) NOT NULL,
    public_name  VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    version_id   UUID         NOT NULL,
    synced_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT workspace_workspace_id_pkey PRIMARY KEY (workspace_id)
);

CREATE TABLE member
(
    member_id    VARCHAR(255) NOT NULL, -- Unique identifier for the member
    workspace_id VARCHAR(255) NOT NULL, -- Reference to workspace this member belongs to
    name         VARCHAR(255) NOT NULL, -- Display name of the member
    public_name  VARCHAR(255) NOT NULL,
    role         VARCHAR(255) NOT NULL, -- Member's role in the workspace
    permissions  JSONB        NOT NULL,
    avatar_url   VARCHAR(511) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    version_id   UUID         NOT NULL,
    synced_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT member_member_id_pkey PRIMARY KEY (member_id),
    CONSTRAINT member_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id)
);

CREATE TABLE customer
(
    customer_id       VARCHAR(255) NOT NULL,               -- primary key
    workspace_id      VARCHAR(255) NOT NULL,               -- fk to workspace
    external_id       VARCHAR(255) NULL,                   -- external id of the customer (optional identifier)
    email             VARCHAR(255) NULL,                   -- email address of the customer (optional identifier)
    phone             VARCHAR(255) NULL,                   -- phone number of the customer (optional identifier)
    name              VARCHAR(255) NOT NULL,               -- display name of the customer
    role              VARCHAR(255) NOT NULL,               -- role/type of the customer
    avatar_url        VARCHAR(511) NOT NULL,
    is_email_verified BOOLEAN      NOT NULL DEFAULT FALSE, -- whether email has been verified
    created_at        TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    version_id        UUID         NOT NULL,
    synced_at         TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT customer_customer_id_pkey PRIMARY KEY (customer_id),
    CONSTRAINT customer_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT customer_workspace_id_external_id_key UNIQUE (workspace_id, external_id),
    CONSTRAINT customer_workspace_id_email_key UNIQUE (workspace_id, email),
    CONSTRAINT customer_workspace_id_phone_key UNIQUE (workspace_id, phone)
);

CREATE TABLE thread
(
    thread_id            VARCHAR(255) NOT NULL,                           -- Unique identifier for the thread
    workspace_id         VARCHAR(255) NOT NULL,                           -- Workspace this thread belongs to
    customer_id          VARCHAR(255) NOT NULL,                           -- Customer associated with this thread
    assignee_id          VARCHAR(255) NULL,                               -- Member assigned to handle this thread
    assigned_at          TIMESTAMP    NULL,                               -- When the thread was assigned
    title                TEXT         NOT NULL,                           -- Thread title
    description          TEXT         NOT NULL,                           -- Thread description
    preview_text         TEXT         NOT NULL,
    status               VARCHAR(127) NOT NULL,                           -- Current status of the thread
    status_changed_at    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP, -- When status was last changed
    status_changed_by_id VARCHAR(255) NOT NULL,                           -- Member who last changed the status
    stage                VARCHAR(127) NOT NULL,                           -- Current stage in the workflow
    replied              BOOLEAN      NOT NULL DEFAULT FALSE,             -- Whether a member has replied
    priority             VARCHAR(255) NOT NULL,                           -- Thread priority level
    channel              VARCHAR(127) NOT NULL,                           -- Communication channel used
    created_by_id        VARCHAR(255) NOT NULL,                           -- Member who created the thread
    updated_by_id        VARCHAR(255) NOT NULL,                           -- Member who last updated the thread
    labels               JSONB        NULL,
    last_inbound_at      TIMESTAMP    NULL,
    last_outbound_at     TIMESTAMP    NULL,
    created_at           TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    version_id           UUID         NOT NULL,
    synced_at            TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT thread_thread_id_pkey PRIMARY KEY (thread_id),
    CONSTRAINT thread_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT thread_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
    CONSTRAINT thread_assignee_id_fkey FOREIGN KEY (assignee_id) REFERENCES member (member_id),
    CONSTRAINT thread_status_changed_by_id_fkey FOREIGN KEY (status_changed_by_id) REFERENCES member (member_id),
    CONSTRAINT thread_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES member (member_id),
    CONSTRAINT thread_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES member (member_id)
);
