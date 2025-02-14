--
-- Please follow the naming convention for consistency.

-- {table name}_{column name(s)}_{suffix}

-- where the suffix is one of the following:

-- `pkey` for a Primary Key constraint.
-- `key` for a Unique constraint.
-- `excl` for an Exclusion constraint.
-- `idx` for any other of index.
-- `fkey` for a Foreign key.
-- `check` for a Check constraint

-- Standard suffix for sequences is
-- seq for all sequences

-- Thanks.
-- --------------------------------------------------

-- Represents the workos user table
-- Mapped directly based on the WorkOS user object
create table workos_user (
    user_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    profile_picture_url VARCHAR(511) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT workos_user_user_id_pkey PRIMARY KEY (user_id),
    CONSTRAINT workos_user_email_key UNIQUE (email) 
);

-- Represents the auth account table
-- This table is used to store the account information of the user pertaining to auth.
-- Attributes will depend on the auth provider.
CREATE TABLE account
(
    account_id   VARCHAR(255) NOT NULL, -- primary key
    email        VARCHAR(255) NOT NULL,
    provider     VARCHAR(255) NOT NULL,
    auth_user_id VARCHAR(255) NOT NULL, -- key to auth provider
    name         VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT account_account_id_pkey PRIMARY KEY (account_id),
    CONSTRAINT account_email_key UNIQUE (email),
    CONSTRAINT account_auth_user_id_key UNIQUE (auth_user_id)
);

-- Represents the account PAT table - Personal Access Token
-- This table is used to store the account PAT information of the account pertaining to auth.
-- PAT is used to authenticate the account similar to API key.
CREATE TABLE account_pat
(
    account_id  VARCHAR(255) NOT NULL, -- fk to account
    pat_id      VARCHAR(255) NOT NULL, -- primary key
    token       VARCHAR(255) NOT NULL, -- unique token across the system
    name        VARCHAR(255) NOT NULL, -- name of the PAT
    description TEXT         NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT account_pat_pat_id_pkey PRIMARY KEY (pat_id),
    CONSTRAINT account_pat_account_id_fkey FOREIGN KEY (account_id) REFERENCES account (account_id),
    CONSTRAINT account_pat_token_key UNIQUE (token)
);

-- Represents the workspace table
-- This table is used to store the workspace information linked to the account.
-- Account can own multiple workspaces.
CREATE TABLE workspace
(
    workspace_id VARCHAR(255) NOT NULL, -- primary key
    account_id   VARCHAR(255) NOT NULL, -- fk to account
    name         VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT workspace_workspace_id_pkey PRIMARY KEY (workspace_id),
    CONSTRAINT workspace_account_id_fkey FOREIGN KEY (account_id) REFERENCES account (account_id)
);

-- Represents the member table
-- This table is used to store the member information linked to the workspace.
-- Each member is uniquely identified by the combination of `workspace_id` and `account_id`
-- Members can authenticate to the workspace through their linked account
-- System members (bots, integrations) have null account_id
-- The role field determines member permissions within the workspace
CREATE TABLE member
(
    member_id    VARCHAR(255) NOT NULL, -- Unique identifier for the member
    workspace_id VARCHAR(255) NOT NULL, -- Reference to workspace this member belongs to
    account_id   VARCHAR(255) NULL,     -- Reference to associated account (can be null for system members)
    name         VARCHAR(255) NOT NULL, -- Display name of the member
    role         VARCHAR(255) NOT NULL, -- Member's role/permissions in the workspace
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT member_member_id_pkey PRIMARY KEY (member_id),
    CONSTRAINT member_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT member_account_id_fkey FOREIGN KEY (account_id) REFERENCES account (account_id),
    CONSTRAINT member_workspace_id_account_id_key UNIQUE (workspace_id, account_id)
);

-- Represents the customer table
-- There can be multiple customers per workspace
-- Each customer is uniquely identified by one of:
-- - workspace_id + external_id
-- - workspace_id + email
-- - workspace_id + phone
CREATE TABLE customer
(
    customer_id       VARCHAR(255) NOT NULL,               -- primary key
    workspace_id      VARCHAR(255) NOT NULL,               -- fk to workspace
    external_id       VARCHAR(255) NULL,                   -- external id of the customer (optional identifier)
    email             VARCHAR(255) NULL,                   -- email address of the customer (optional identifier)
    phone             VARCHAR(255) NULL,                   -- phone number of the customer (optional identifier)
    name              VARCHAR(255) NOT NULL,               -- display name of the customer
    role              VARCHAR(255) NOT NULL,               -- role/type of the customer
    is_email_verified BOOLEAN      NOT NULL DEFAULT FALSE, -- whether email has been verified
    created_at        TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT customer_customer_id_pkey PRIMARY KEY (customer_id),
    CONSTRAINT customer_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT customer_workspace_id_external_id_key UNIQUE (workspace_id, external_id),
    CONSTRAINT customer_workspace_id_email_key UNIQUE (workspace_id, email),
    CONSTRAINT customer_workspace_id_phone_key UNIQUE (workspace_id, phone)
);

CREATE TABLE customer_event
(
    event_id    VARCHAR(255) NOT NULL,               -- primary key
    customer_id VARCHAR(255) NOT NULL,               -- fk to customer
    title       VARCHAR(511) NOT NULL,               -- title of the event
    severity    VARCHAR(127) NOT NULL,               -- severity level of the event
    timestamp   TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- when the event occurred
    components  JSONB        NOT NULL,               -- custom JSON data structure for the event

    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- record creation timestamp
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- record last updated timestamp

    CONSTRAINT customer_event_event_id_pkey PRIMARY KEY (event_id),
    CONSTRAINT customer_event_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id)
);

CREATE TABLE postmark_setting
(
    workspace_id                VARCHAR(255) NOT NULL,
    server_id                   BIGINT       NOT NULL,
    server_token                VARCHAR(255) NOT NULL,
    is_enabled                  BOOLEAN      NOT NULL DEFAULT FALSE,
    email                       VARCHAR(255) NOT NULL,
    domain                      VARCHAR(255) NOT NULL,
    has_error                   BOOLEAN      NOT NULL DEFAULT FALSE,
    inbound_email               VARCHAR(255) NULL,
    has_forwarding_enabled      BOOLEAN      NOT NULL DEFAULT FALSE,
    has_dns                     BOOLEAN      NOT NULL DEFAULT FALSE,
    is_dns_verified             BOOLEAN      NOT NULL DEFAULT FALSE,
    dns_verified_at             TIMESTAMP    NULL,
    dns_domain_id               BIGINT       NULL,
    dkim_host                   VARCHAR(255) NULL,
    dkim_text_value             TEXT         NULL,
    dkim_update_status          VARCHAR(255) NULL,
    return_path_domain          VARCHAR(255) NULL,
    return_path_domain_cname    VARCHAR(255) NULL,
    return_path_domain_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at                  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT postmark_setting_workspace_id_pkey PRIMARY KEY (workspace_id),
    CONSTRAINT postmark_setting_inbound_email_key UNIQUE (inbound_email)
);

-- Represents a thread which is a conversation between a customer and members
-- Each thread belongs to a workspace and is associated with a customer
-- Threads can be assigned to members and have various statuses and priorities
CREATE TABLE thread
(
    thread_id            VARCHAR(255) NOT NULL,                           -- Unique identifier for the thread
    workspace_id         VARCHAR(255) NOT NULL,                           -- Workspace this thread belongs to
    customer_id          VARCHAR(255) NOT NULL,                           -- Customer associated with this thread
    assignee_id          VARCHAR(255) NULL,                               -- Member assigned to handle this thread
    assigned_at          TIMESTAMP    NULL,                               -- When the thread was assigned
    title                TEXT         NOT NULL,                           -- Thread title
    description          TEXT         NOT NULL,                           -- Thread description
    status               VARCHAR(127) NOT NULL,                           -- Current status of the thread
    status_changed_at    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP, -- When status was last changed
    status_changed_by_id VARCHAR(255) NOT NULL,                           -- Member who last changed the status
    stage                VARCHAR(127) NOT NULL,                           -- Current stage in the workflow
    replied              BOOLEAN      NOT NULL DEFAULT FALSE,             -- Whether a member has replied
    priority             VARCHAR(255) NOT NULL,                           -- Thread priority level
    channel              VARCHAR(127) NOT NULL,                           -- Communication channel used
    last_inbound_at      TIMESTAMP    NULL,
    last_outbound_at     TIMESTAMP    NULL,
    created_by_id        VARCHAR(255) NOT NULL,                           -- Member who created the thread
    updated_by_id        VARCHAR(255) NOT NULL,                           -- Member who last updated the thread
    created_at           TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT thread_thread_id_pkey PRIMARY KEY (thread_id),
    CONSTRAINT thread_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT thread_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
    CONSTRAINT thread_assignee_id_fkey FOREIGN KEY (assignee_id) REFERENCES member (member_id),
    CONSTRAINT thread_status_changed_by_id_fkey FOREIGN KEY (status_changed_by_id) REFERENCES member (member_id),
    CONSTRAINT thread_created_by_id_fkey FOREIGN KEY (created_by_id) REFERENCES member (member_id),
    CONSTRAINT thread_updated_by_id_fkey FOREIGN KEY (updated_by_id) REFERENCES member (member_id)
);

create table activity
(
    activity_id   varchar(255) not null,
    thread_id     varchar(255) not null,
    activity_type varchar(255) not null,
    body          jsonb        not null,
    customer_id   varchar(255) null,
    member_id     varchar(255) null,
    created_at    timestamp default current_timestamp,
    updated_at    timestamp default current_timestamp,

    CONSTRAINT activity_activity_id_pkey PRIMARY KEY (activity_id),
    CONSTRAINT activity_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES thread (thread_id),
    CONSTRAINT activity_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
    CONSTRAINT activity_member_id_fkey FOREIGN KEY (member_id) REFERENCES member (member_id),
    CONSTRAINT activity_participant_check CHECK (num_nonnulls(customer_id, member_id) = 1)
);

CREATE TABLE activity_attachment
(
    attachment_id VARCHAR(255) NOT NULL,
    activity_id    VARCHAR(255) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    content_type  VARCHAR(511) NOT NULL,
    content_key   VARCHAR(511) NOT NULL,
    content_url   VARCHAR(511) NOT NULL,
    spam          BOOLEAN      NOT NULL DEFAULT FALSE,
    has_error     BOOLEAN      NOT NULL DEFAULT FALSE,
    error         TEXT         NOT NULL,
    md5_hash      VARCHAR(511) NOT NULL,
    created_at    timestamp             DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT activity_attachment_id_pkey PRIMARY KEY (attachment_id),
    CONSTRAINT activity_attachment_activity_id_fkey FOREIGN KEY (activity_id) REFERENCES activity (activity_id)
);

CREATE TABLE postmark_message_log
(
    activity_id           VARCHAR(255) NOT NULL, -- References parent activity
    payload               JSONB        NOT NULL, -- Request payload
    postmark_message_id   VARCHAR(255) NOT NULL, -- Postmark's internal message ID
    mail_message_id       VARCHAR(255) NOT NULL, -- Email `Message-ID` header
    reply_mail_message_id VARCHAR(255) NULL,     -- Email `In-Reply-To` header
    has_error             BOOLEAN      NOT NULL DEFAULT FALSE,
    submitted_at          TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    error_code            bigint       NOT NULL,
    postmark_message      VARCHAR(255) NOT NULL,
    message_event         VARCHAR(255) NOT NULL,
    acknowledged          BOOLEAN      NOT NULL DEFAULT FALSE,
    message_type          VARCHAR(255) NOT NULL,
    created_at            TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT postmark_log_message_id_pkey PRIMARY KEY (activity_id),
    CONSTRAINT postmark_log_message_id_fkey FOREIGN KEY (activity_id) REFERENCES activity (activity_id),

    CONSTRAINT postmark_log_pm_message_id_key UNIQUE (postmark_message_id),
    CONSTRAINT postmark_log_mail_message_id_key UNIQUE (mail_message_id)
);
CREATE INDEX postmark_log_reply_mail_message_idx ON postmark_message_log (reply_mail_message_id);

-- Represents the label table
-- This table is used to store the labels linked to the workspace.
-- Each label is uniquely identified by the combination of `workspace_id` and `name`
CREATE TABLE label
(
    workspace_id VARCHAR(255) NOT NULL,
    label_id     VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    icon         VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT label_label_id_pkey PRIMARY KEY (label_id),
    CONSTRAINT label_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT label_workspace_id_name_key UNIQUE (workspace_id, name)
);

CREATE TABLE thread_label
(
    thread_label_id VARCHAR(255) NOT NULL,
    thread_id       VARCHAR(255) NOT NULL,
    label_id        VARCHAR(255) NOT NULL,
    addedby         VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT thread_label_thread_label_id_pkey PRIMARY KEY (thread_label_id),
    CONSTRAINT thread_label_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES thread (thread_id),
    CONSTRAINT thread_label_label_id_fkey FOREIGN KEY (label_id) REFERENCES label (label_id),
    CONSTRAINT thread_label_thread_label_id_key UNIQUE (thread_id, label_id)
);

-- Represents the secret key table
-- This table is used to store the secret key linked to the workspace.
CREATE TABLE workspace_secret
(
    workspace_id VARCHAR(255) NOT NULL,
    hmac         VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT workspace_secret_workspace_id_pkey PRIMARY KEY (workspace_id),
    CONSTRAINT workspace_secret_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT workspace_secret_hmac_key UNIQUE (hmac)
);

-- ************************************ --
-- tables below have been changed or deprecated.
-- ************************************ --

-- Represents the widget session table
-- This table is used to store the widget session linked to the widget.
-- CREATE TABLE widget_session
-- (
--     session_id VARCHAR(255) NOT NULL,
--     widget_id  VARCHAR(255) NOT NULL,
--     data       TEXT         NOT NULL,
--     expire_at  TIMESTAMP    NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--
--     CONSTRAINT widget_session_session_id_pkey PRIMARY KEY (session_id),
--     CONSTRAINT widget_session_widget_id_fkey FOREIGN KEY (widget_id) REFERENCES widget (widget_id)
-- );

-- Represents the widget table
-- This table is used to store the widgets linked to the workspace.
-- CREATE TABLE widget
-- (
--     workspace_id  VARCHAR(255) NOT NULL,
--     widget_id     VARCHAR(255) NOT NULL,
--     name          VARCHAR(255) NOT NULL,
--     configuration JSONB        NOT NULL,
--     created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--
--     CONSTRAINT widget_widget_id_pkey PRIMARY KEY (widget_id),
--     CONSTRAINT widget_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id)
-- );

-- Represents the multichannel thread message.
-- This table stores messages that are part of a thread, supporting multiple communication channels.
-- Messages can be from either a customer or a member (but not both).
-- Each message has both a text_body (plain text) and body (formatted/rich text).
-- CREATE TABLE message
-- (
--     message_id    VARCHAR(255) NOT NULL,               -- Unique identifier for the message
--     thread_id     VARCHAR(255) NOT NULL,               -- Thread this message belongs to
--     text_body     TEXT         NOT NULL,               -- Plain text content of the message
--     markdown_body TEXT         NOT NULL,               -- Rich text/formatted content of the message
--     html_body     TEXT         NOT NULL,               -- Rich text/formatted HTML content of the message
--     customer_id   VARCHAR(255) NULL,                   -- Customer who sent the message (if from customer)
--     member_id     VARCHAR(255) NULL,                   -- Member who sent the message (if from member)
--     channel       VARCHAR(255) NOT NULL,               -- Communication channel used (email, chat, etc)
--     created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the message was created
--     updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the message was last updated
--
--     -- Defining the primary key for the table
--     CONSTRAINT message_message_id_pkey PRIMARY KEY (message_id),
--
--     -- Foreign key constraints
--     CONSTRAINT message_thread_id_fkey FOREIGN KEY (thread_id) REFERENCES thread (thread_id),
--     CONSTRAINT message_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
--     CONSTRAINT message_member_id_fkey FOREIGN KEY (member_id) REFERENCES member (member_id),
--
--     -- Check constraint to enforce valid sender (only one of customer_id or member_id can be set)
--     CONSTRAINT message_sender_check CHECK (
--         (customer_id IS NULL AND member_id IS NOT NULL) OR
--         (customer_id IS NOT NULL AND member_id IS NULL)
--         )
-- );

-- @sanchitrk: changed usage?
-- Represents the workspace Thread QA table
-- This table is used to store the QA thread information linked to the workspace.
-- CREATE TABLE thread_qa (
--     workspace_id VARCHAR(255) NOT NULL,
--     customer_id VARCHAR(255) NOT NULL,
--     thread_id VARCHAR(255) NOT NULL,
--     parent_thread_id VARCHAR(255) NULL,
--     query TEXT NOT NULL,
--     title TEXT NOT NULL,
--     summary TEXT NOT NULL,
--     sequence BIGINT NOT NULL DEFAULT fn_next_id(),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--
--     CONSTRAINT thread_qa_thread_id_pkey PRIMARY KEY (thread_id),
--     CONSTRAINT thread_qa_parent_thread_id_fkey FOREIGN KEY (parent_thread_id) REFERENCES thread_qa (thread_id),
--     CONSTRAINT thread_qa_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
--     CONSTRAINT thread_qa_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
--     CONSTRAINT thread_qa_thread_id_parent_thread_id UNIQUE (thread_id, parent_thread_id)
-- );

-- @sanchitrk: changed usage?
-- CREATE TABLE thread_qa_answer (
--     workspace_id VARCHAR(255) NOT NULL,
--     thread_qa_id VARCHAR(255) NOT NULL,
--     answer_id VARCHAR(255) NOT NULL,
--     answer TEXT NOT NULL,
--     eval INT NULL DEFAULT NULL,
--     sequence BIGINT NOT NULL DEFAULT fn_next_id(),
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--
--     CONSTRAINT thread_qa_answer_answer_id_pkey PRIMARY KEY (answer_id),
--     CONSTRAINT thread_qa_answer_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
--     CONSTRAINT thread_qa_answer_thread_qa_id_fkey FOREIGN KEY (thread_qa_id) REFERENCES thread_qa (thread_id)
-- );

-- CREATE TABLE claimed_mail
-- (
--     claim_id      VARCHAR(255) NOT NULL,
--     workspace_id  VARCHAR(255) NOT NULL,
--     customer_id   VARCHAR(255) NOT NULL,
--     email         VARCHAR(255) NOT NULL,
--     has_conflict  BOOLEAN      NOT NULL DEFAULT TRUE,
--     expires_at    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
--     token         TEXT         NOT NULL,
--     is_mail_sent  BOOLEAN      NOT NULL DEFAULT FALSE,
--     platform      VARCHAR(255) NULL,
--     sender_id     VARCHAR(255) NULL,
--     sender_status VARCHAR(255) NULL,
--     sent_at       TIMESTAMP    NULL,
--     created_at    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
--     updated_at    TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
--
--     CONSTRAINT claimed_mail_id_pkey PRIMARY KEY (claim_id),
--     CONSTRAINT claimed_mail_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace,
--     CONSTRAINT claimed_mail_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer
-- );

-- sanchitrk: changed usage?
-- Represents the Event table
-- There can be multiple events per workspace
-- Each event is uniquely identified by the event_id
-- Event can be linked to customer or thread
-- Global events are not linked to customer or thread
CREATE TABLE event
(
    event_id     VARCHAR(255) NOT NULL,                      -- primary key
    workspace_id VARCHAR(255) NOT NULL,                      -- fk to workspace
    sequence     BIGINT       NOT NULL DEFAULT fn_next_id(), -- sequence number of the event
    ts           BIGINT       NOT NULL,                      -- epoch timestamp of the event
    severity     VARCHAR(127) NOT NULL,                      -- severity of the event e.g. info, warning, error, etc.
    category     VARCHAR(127) NOT NULL,                      -- category of the event e.g. auth, payments, general, etc.
    title        VARCHAR(255) NULL,                          -- title of the event
    body         TEXT         NOT NULL,                      -- body of the event
    customer_id  VARCHAR(255) NULL,                          -- fk to customer if event is customer
    created_at   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT event_event_id_pkey PRIMARY KEY (event_id),
    CONSTRAINT event_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT event_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id)
);

-- ************************************ --
-- Will work on Slack integration later --
-- ************************************ --

-- Represents the Slack workspace table
-- There is only one Slack workspace per Workspace
CREATE TABLE slack_workspace
(
    workspace_id VARCHAR(255) NOT NULL,          -- fk to workspace
    ref          VARCHAR(255) NOT NULL,          -- primary key and reference to Slack workspace or team id
    url          VARCHAR(255) NOT NULL,          -- Slack workspace url
    name         VARCHAR(255) NOT NULL,          -- Slack workspace name
    status       VARCHAR(127) NOT NULL,          -- current status of Slack workspace with respect to Workspace
    sync_status  VARCHAR(127) NOT NULL,          -- current sync status of Slack workspace
    synced_at    TIMESTAMP    NULL DEFAULT NULL, -- last time Slack workspace was synced defaults to NULL
    created_at   TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT slack_workspace_ref_pkey PRIMARY KEY (ref),
    CONSTRAINT slack_workspace_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT slack_workspace_workspace_id_key UNIQUE (workspace_id)
);

-- Represents the Slack bot table
-- There is only one Slack bot per Slack workspace, indirectly there is only one Slack bot per Workspace
CREATE TABLE slack_bot
(
    slack_workspace_ref VARCHAR(255) NOT NULL, -- fk to slack_workspace
    bot_id              VARCHAR(255) NOT NULL, -- primary key
    bot_user_ref        VARCHAR(255) NOT NULL, -- reference to Slack bot user id
    bot_ref             VARCHAR(255) NULL,     -- reference to Slack bot id
    app_ref             VARCHAR(255) NOT NULL, -- reference to Slack app id
    scope               TEXT         NOT NULL, -- comma separated list of scopes
    access_token        VARCHAR(255) NOT NULL, -- access token for the bot
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT slack_bot_bot_id_pkey PRIMARY KEY (bot_id),
    CONSTRAINT slack_bot_slack_workspace_ref_fkey FOREIGN KEY (slack_workspace_ref) REFERENCES slack_workspace (ref),
    CONSTRAINT slack_bot_slack_workspace_ref_key UNIQUE (slack_workspace_ref)
);

-- Represents the Slack channel table
-- There are many Slack channels per Slack workspace
CREATE TABLE slack_channel
(
    slack_workspace_ref   VARCHAR(255) NOT NULL,          -- fk to slack_workspace
    channel_id            VARCHAR(255) NOT NULL,          -- primary key
    channel_ref           VARCHAR(255) NOT NULL,          -- reference to Slack channel ID.
    is_channel            BOOLEAN      NOT NULL,
    is_ext_shared         BOOLEAN      NOT NULL,
    is_general            BOOLEAN      NOT NULL,
    is_group              BOOLEAN      NOT NULL,
    is_im                 BOOLEAN      NOT NULL,
    is_member             BOOLEAN      NOT NULL,
    is_mpim               BOOLEAN      NOT NULL,
    is_org_shared         BOOLEAN      NOT NULL,
    is_pending_ext_shared BOOLEAN      NOT NULL,
    is_private            BOOLEAN      NOT NULL,
    is_shared             BOOLEAN      NOT NULL,
    name                  VARCHAR(255) NOT NULL,          -- Slack channel name
    name_normalized       VARCHAR(255) NOT NULL,          -- Slack channel name normalized
    created               BIGINT       NOT NULL,
    updated               BIGINT       NOT NULL,
    status                VARCHAR(127) NOT NULL,          -- custom status of Slack Channel with respect to Slack workspace
    synced_at             TIMESTAMP    NULL DEFAULT NULL, -- custom timestamp Slack channel was synced defaults to NULL
    created_at            TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT slack_channel_channel_id_pkey PRIMARY KEY (channel_id),
    CONSTRAINT slack_channel_slack_workspace_ref_fkey FOREIGN KEY (slack_workspace_ref) REFERENCES slack_workspace (ref),
    CONSTRAINT slack_channel_slack_workspace_ref_channel_ref_key UNIQUE (slack_workspace_ref, channel_ref)
);

-- Stored procedure to generate next id
CREATE OR REPLACE FUNCTION fn_next_id(OUT result bigint) AS
$$
DECLARE
    start_epoch bigint := 1704047400000;
    now_millis  bigint;
BEGIN
    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - start_epoch);
END;
$$ LANGUAGE PLPGSQL;
