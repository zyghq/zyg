from datetime import datetime
from typing import List

from pydantic import BaseModel


class DBEntity(BaseModel):
    created_at: datetime | None = None  # db timestamp
    updated_at: datetime | None = None  # db timestamp


class TenantDBEntity(DBEntity):
    tenant_id: str | None = None  # primary key
    name: str
    slack_team_ref: str | None = None


class UserDBEntity(DBEntity):
    user_id: str | None = None  # primary key
    tenant_id: str
    slack_user_ref: str
    name: str
    role: str


class SlackEventDBEntity(DBEntity):
    event_id: str | None = None  # primary key
    tenant_id: str
    slack_event_ref: str
    inner_event_type: str | None = None
    event_dispatched_ts: int
    api_app_id: str
    token: str
    payload: dict
    is_ack: bool = False


class InSyncSlackChannelDBEntity(DBEntity):
    tenant_id: str
    context_team_id: str
    created: int
    creator: str
    id: str
    is_archived: bool
    is_channel: bool
    is_ext_shared: bool
    is_general: bool
    is_group: bool
    is_im: bool
    is_member: bool
    is_mpim: bool
    is_org_shared: bool
    is_pending_ext_shared: bool
    is_private: bool
    is_shared: bool
    name: str
    name_normalized: str
    num_members: int
    parent_conversation: str | None = None
    pending_connected_team_ids: List[str] | None = None
    pending_shared: List[str] | None = None
    previous_names: List[str] | None = None
    purpose: dict | None = None
    shared_team_ids: List[str] | None = None
    topic: dict | None = None
    unlinked: int | None = None
    updated: int


class SlackChannelDBEntity(DBEntity):
    tenant_id: str
    slack_channel_id: str | None = None  # primary key
    slack_channel_ref: str
    slack_channel_name: str | None = None
    triage_slack_channel_ref: str
    triage_slack_channel_name: str | None = None


class IssueDBEntity(DBEntity):
    tenant_id: str
    issue_id: str | None = None  # primary key
    issue_number: int | None = None
    slack_channel_id: str
    slack_message_ts: str
    body: str
    status: str
    priority: int
    tags: List[str] | None = []


class InSyncSlackUserDBEntity(DBEntity):
    tenant_id: str
    id: str
    is_admin: bool
    is_app_user: bool
    is_bot: bool
    is_email_confirmed: bool
    is_owner: bool
    is_primary_owner: bool
    is_restricted: bool
    is_stranger: bool | None = None
    is_ultra_restricted: bool
    name: str
    profile: dict
    real_name: str
    team_id: str
    tz: str
    tz_label: str
    tz_offset: int
    updated: int
