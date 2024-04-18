-- Please follow the naming convention for consistency.

-- {tablename}_{columnname(s)}_{suffix}

-- where the suffix is one of the following:

-- pkey for a Primary Key constraint
-- key for a Unique constraint
-- excl for an Exclusion constraint
-- idx for any other kind of index
-- fkey for a Foreign key
-- check for a Check constraint

-- Standard suffix for sequences is
-- seq for all sequences

-- Thanks.
-- --------------------------------------------------

-- Represents the auth account table
-- This table is used to store the account information of the user pertaining to auth.
-- Attributes will depend on the auth provider.
-- Is Confirmed
CREATE TABLE account (
    account_id VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    provider VARCHAR(255) NOT NULL, 
    auth_user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT account_account_id_pkey PRIMARY KEY (account_id),
    CONSTRAINT account_email_key UNIQUE (email),
    CONSTRAINT account_auth_user_id_key UNIQUE (auth_user_id)
);


-- Represents the account PAT table
-- Personal Access Token
-- Is Confirmed
CREATE TABLE account_pat (
    account_id VARCHAR(255) NOT NULL, -- fk to account
    pat_id VARCHAR(255) NOT NULL, -- primary key
    token VARCHAR(255) NOT NULL, -- unique token across the system
    name VARCHAR(255) NOT NULL, -- name of the PAT
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT account_pat_pat_id_pkey PRIMARY KEY (pat_id),
    CONSTRAINT account_pat_account_id_fkey FOREIGN KEY (account_id) REFERENCES account (account_id),
    CONSTRAINT account_pat_token_key UNIQUE (token)
);

-- Represents the workspace table
-- This table is used to store the workspace information linked to the account.
-- Account can own multiple workspaces.
-- Is Confirmed
CREATE TABLE workspace (
    workspace_id VARCHAR(255) NOT NULL,
    account_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT workspace_workspace_id_pkey PRIMARY KEY (workspace_id),
    CONSTRAINT workspace_account_id_fkey FOREIGN KEY (account_id) REFERENCES account (account_id)
);

-- Represents the workspace Thread QA table
-- This table is used to store the QA thread information linked to the workspace.
CREATE TABLE thread_qa (
    workspace_id VARCHAR(255) NOT NULL,
    customer_id VARCHAR(255) NOT NULL,
    thread_id VARCHAR(255) NOT NULL,
    parent_thread_id VARCHAR(255) NULL,
    query TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    sequence BIGINT NOT NULL DEFAULT fn_next_id(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT thread_qa_thread_id_pkey PRIMARY KEY (thread_id),
    CONSTRAINT thread_qa_parent_thread_id_fkey FOREIGN KEY (parent_thread_id) REFERENCES thread_qa (thread_id),
    CONSTRAINT thread_qa_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT thread_qa_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
    CONSTRAINT thread_qa_thread_id_parent_thread_id UNIQUE (thread_id, parent_thread_id)
);

CREATE TABLE thread_qa_answer (
    workspace_id VARCHAR(255) NOT NULL,
    thread_qa_id VARCHAR(255) NOT NULL,
    answer_id VARCHAR(255) NOT NULL,
    answer TEXT NOT NULL,
    eval INT NULL DEFAULT NULL,
    sequence BIGINT NOT NULL DEFAULT fn_next_id(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT thread_qa_answer_answer_id_pkey PRIMARY KEY (answer_id),
    CONSTRAINT thread_qa_answer_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT thread_qa_answer_thread_qa_id_fkey FOREIGN KEY (thread_qa_id) REFERENCES thread_qa (thread_id)
);


-- Represents the chat thread table
-- This table is used to store the chat thread information linked to the workspace.
CREATE TABLE thread_chat (
    workspace_id VARCHAR(255) NOT NULL,
    customer_id VARCHAR(255) NOT NULL,
    assignee_id VARCHAR(255) NULL,
    thread_chat_id VARCHAR(255) NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    sequence BIGINT NOT NULL DEFAULT fn_next_id(),
    status VARCHAR(127) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT thread_chat_thread_chat_id_pkey PRIMARY KEY (thread_chat_id),
    CONSTRAINT thread_chat_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT thread_chat_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
    CONSTRAINT thread_chat_assignee_id_fkey FOREIGN KEY (assignee_id) REFERENCES member (member_id)
);


-- Represents the chat message table
-- This table is used to store the chat message information linked to the chat thread.
CREATE TABLE thread_chat_message (
    thread_chat_id VARCHAR(255) NOT NULL,
    thread_chat_message_id VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    sequence BIGINT NOT NULL DEFAULT fn_next_id(),
    customer_id VARCHAR(255) NULL,
    member_id VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT thread_chat_message_thread_chat_message_id_pkey PRIMARY KEY (thread_chat_message_id),
    CONSTRAINT thread_chat_message_thread_chat_id_fkey FOREIGN KEY (thread_chat_id) REFERENCES thread_chat (thread_chat_id),
    CONSTRAINT thread_chat_message_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
    CONSTRAINT thread_chat_message_member_id_fkey FOREIGN KEY (member_id) REFERENCES member (member_id)
);

-- Represents the member table
-- This table is used to store the member information linked to the workspace.
-- Each member is uniquely identified by the combination of workspace_id and account_id
-- Member has ability to authenticate to the workspace, hence has link to account
-- Done
CREATE TABLE member (
    member_id VARCHAR(255) NOT NULL, -- primary key
    workspace_id VARCHAR(255) NOT NULL, -- fk to workspace
    account_id VARCHAR(255) NOT NULL, -- fk to account
    name VARCHAR NULL, -- name of the member
    role VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT member_member_id_pkey PRIMARY KEY (member_id),
    CONSTRAINT member_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT member_account_id_fkey FOREIGN KEY (account_id) REFERENCES account (account_id),
    CONSTRAINT member_workspace_id_account_id_key UNIQUE (workspace_id, account_id),
    CONSTRAINT member_slug_key UNIQUE (slug)
);

-- Represents the Customer table
-- There can be multiple customers per workspace
-- Each customer is uniquely identified by one of external_id, email and phone, each unique to the workspace
-- Done
CREATE TABLE customer (
    customer_id VARCHAR(255) NOT NULL, -- primary key
    workspace_id VARCHAR(255) NOT NULL, -- fk to workspace
    external_id VARCHAR(255) NULL, -- external id of the customer
    email VARCHAR(255) NULL, -- email of the customer
    phone VARCHAR(255) NULL, -- phone of the customer
    name VARCHAR(255)  NULL, -- name of the customer
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT customer_customer_id_pkey PRIMARY KEY (customer_id),
    CONSTRAINT customer_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT customer_workspace_id_external_id_key UNIQUE (workspace_id, external_id),
    CONSTRAINT customer_workspace_id_email_key UNIQUE (workspace_id, email),
    CONSTRAINT customer_workspace_id_phone_key UNIQUE (workspace_id, phone)
);

-- Represents the Member Key table
-- There can be multiple keys per Workspace
-- Think of Member Key as alias for member account.
-- In the future we can have permissions on API keys
-- Deprecate This later
CREATE TABLE member_key (
    member_key_id VARCHAR(255) NOT NULL, -- primary key
    workspace_id VARCHAR(255) NOT NULL, -- fk to workspace
    member_id VARCHAR(255) NOT NULL, -- fk to member
    token VARCHAR(255) NOT NULL, -- unique API key across the system
    name VARCHAR(255) NOT NULL, -- name of the API key
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT member_key_member_key_id_pkey PRIMARY KEY (member_key_id),
    CONSTRAINT member_key_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT member_key_member_id_fkey FOREIGN KEY (member_id) REFERENCES member (member_id),
    CONSTRAINT member_key_token_key UNIQUE (token)
);

-- Represents the Event table
-- There can be multiple events per workspace
-- Each event is uniquely identified by the event_id
-- Event can be linked to customer or thread
-- Global events are not linked to customer or thread
-- TODO: add support for `threads`
CREATE TABLE event (
    event_id VARCHAR(255) NOT NULL, -- primary key
    workspace_id VARCHAR(255) NOT NULL, -- fk to workspace
    sequence BIGINT NOT NULL DEFAULT fn_next_id(), -- sequence number of the event
    ts BIGINT NOT NULL, -- epoch timestamp of the event
    severity VARCHAR(127) NOT NULL, -- severity of the event e.g. info, warning, error, etc.
    category VARCHAR(127) NOT NULL, -- category of the event e.g. auth, payments, general, etc.
    title VARCHAR(255) NULL, -- title of the event
    body TEXT NOT NULL, -- body of the event
    customer_id VARCHAR(255) NULL, -- fk to customer if event is customer specific
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT event_event_id_pkey PRIMARY KEY (event_id),
    CONSTRAINT event_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT event_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id)
);


-- Represents the Chat Thread table
-- There can be multiple threads per Workspace for a customer
CREATE TABLE chat_thread (
    workspace_id VARCHAR(255) NOT NULL, -- fk to workspace
    customer_id VARCHAR(255) NOT NULL, -- fk to customer
    thread_id VARCHAR(255) NOT NULL, -- primary key
    sequence BIGINT NOT NULL DEFAULT fn_next_id(), -- sequence number of the thread
    assignee_id VARCHAR(255) NULL, -- fk to member
    priority VARCHAR(127) NOT NULL, -- priority of the thread
    status VARCHAR(127) NOT NULL, -- status of the thread
    body TEXT NULL, -- body of the thread
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chat_thread_thread_id_pkey PRIMARY KEY (thread_id),
    CONSTRAINT chat_thread_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT chat_thread_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
    CONSTRAINT chat_thread_assignee_id_id_fkey FOREIGN KEY (assignee_id) REFERENCES member (member_id)
);

-- Represents the Workspace In App Chat Key
-- There can be multiple keys per Workspace
-- Chat Key is used to initiate chat with the customer for a workspace
CREATE TABLE chat_key (
    workspace_id VARCHAR(255) NOT NULL, -- fk to workspace
    chat_key_id VARCHAR(255) NOT NULL, -- primary key
    key VARCHAR(255) NOT NULL, -- unique key across the system
    name VARCHAR(255) NOT NULL, -- name of the key
    description TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chat_key_chat_key_id_pkey PRIMARY KEY (chat_key_id),
    CONSTRAINT chat_key_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT chat_key_key_key UNIQUE (key)
);

-- Represents the Customer Chat Session
-- There can be multiple sessions per Workspace
-- Each socket connection is unique to the customer
CREATE TABLE customer_chat_session (
    workspace_id VARCHAR(255) NOT NULL, -- fk to workspace
    customer_id VARCHAR(255) NOT NULL, -- fk to customer
    customer_chat_session_id VARCHAR(255) NOT NULL, -- primary key
    key VARCHAR(255) NOT NULL, -- unique key across the system
    chat_thread_id VARCHAR(255) NULL, -- fk to chat_thread
    socket_id VARCHAR(255) NOT NULL, -- socket id from provider
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chat_session_customer_chat_session_id_pkey PRIMARY KEY (customer_chat_session_id),
    CONSTRAINT chat_session_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT chat_session_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES customer (customer_id),
    CONSTRAINT chat_session_chat_thread_id_fkey FOREIGN KEY (chat_thread_id) REFERENCES chat_thread (thread_id),
    CONSTRAINT chat_session_key_key UNIQUE (key),
    CONSTRAINT chat_session_customer_id_socket_id_key UNIQUE (customer_id, socket_id)
);

-- ************************************ --
-- Will work on Slack integration later --
-- ************************************ --

-- Represents theh Slack workspace table
-- There is only one Slack workspace per Workspace
CREATE TABLE slack_workspace (
    workspace_id VARCHAR(255) NOT NULL, -- fk to workspace
    ref VARCHAR(255) NOT NULL, -- primary key and reference to Slack workspace or team id
    url VARCHAR(255) NOT NULL, -- Slack workspace url
    name VARCHAR(255) NOT NULL, -- Slack workspace name
    status VARCHAR(127) NOT NULL, -- current status of Slack workspace with respect to Workspace
    sync_status VARCHAR(127) NOT NULL, -- current sync status of Slack workspace
    synced_at TIMESTAMP NULL DEFAULT NULL, -- last time Slack workspace was synced defaults to NULL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT slack_workspace_ref_pkey PRIMARY KEY (ref),
    CONSTRAINT slack_workspace_workspace_id_fkey FOREIGN KEY (workspace_id) REFERENCES workspace (workspace_id),
    CONSTRAINT slack_workspace_workspace_id_key UNIQUE (workspace_id)
);

-- Represents the Slack bot table
-- There is only one Slack bot per Slack workspace, indirectly there is only one Slack bot per Workspace
CREATE TABLE slack_bot (
    slack_workspace_ref VARCHAR(255) NOT NULL, -- fk to slack_workspace
    bot_id VARCHAR(255) NOT NULL, -- primary key
    bot_user_ref VARCHAR(255) NOT NULL, -- reference to Slack bot user id
    bot_ref VARCHAR(255) NULL, -- reference to Slack bot id
    app_ref VARCHAR(255) NOT NULL, -- reference to Slack app id
    scope TEXT NOT NULL, -- comma separated list of scopes
    access_token VARCHAR(255) NOT NULL, -- access token for the bot
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT slack_bot_bot_id_pkey PRIMARY KEY (bot_id),
    CONSTRAINT slack_bot_slack_workspace_ref_fkey FOREIGN KEY (slack_workspace_ref) REFERENCES slack_workspace (ref),
    CONSTRAINT slack_bot_slack_workspace_ref_key UNIQUE (slack_workspace_ref)
);

-- Represents the Slack channel table
-- There are many Slack channels per Slack workspace
CREATE TABLE slack_channel (
    slack_workspace_ref VARCHAR(255) NOT NULL, -- fk to slack_workspace
    channel_id VARCHAR(255) NOT NULL, -- primary key
    channel_ref VARCHAR(255) NOT NULL, -- reference to Slack channel id
    is_channel BOOLEAN NOT NULL,
    is_ext_shared BOOLEAN NOT NULL,
    is_general BOOLEAN NOT NULL,
    is_group BOOLEAN NOT NULL,
    is_im BOOLEAN NOT NULL,
    is_member BOOLEAN NOT NULL,
    is_mpim BOOLEAN NOT NULL,
    is_org_shared BOOLEAN NOT NULL,
    is_pending_ext_shared BOOLEAN NOT NULL,
    is_private BOOLEAN NOT NULL,
    is_shared BOOLEAN NOT NULL,
    name VARCHAR(255) NOT NULL, -- Slack channel name
    name_normalized VARCHAR(255) NOT NULL, -- Slack channel name normalized
    created BIGINT NOT NULL,
    updated BIGINT NOT NULL,
    status VARCHAR(127) NOT NULL, -- custom status of Slack channel with respect to Slack workspace
    synced_at TIMESTAMP NULL DEFAULT NULL, -- custom timestamp Slack channel was synced defaults to NULL
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT slack_channel_channel_id_pkey PRIMARY KEY (channel_id),
    CONSTRAINT slack_channel_slack_workspace_ref_fkey FOREIGN KEY (slack_workspace_ref) REFERENCES slack_workspace (ref),
    CONSTRAINT slack_channel_slack_workspace_ref_channel_ref_key UNIQUE (slack_workspace_ref, channel_ref)
);

-- Stored procedure to generate next id
CREATE OR REPLACE FUNCTION fn_next_id(OUT result bigint) AS $$
DECLARE
    start_epoch bigint := 1704047400000;
    now_millis bigint;
BEGIN
    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - start_epoch);
END;
$$ LANGUAGE PLPGSQL;