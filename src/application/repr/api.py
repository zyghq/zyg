from datetime import datetime
from typing import Dict, List

from pydantic import BaseModel

from src.domain.models import (
    InSyncSlackChannel,
    InSyncSlackUser,
    Issue,
    SlackChannel,
    SlackEvent,
    User,
)


class TenantRepr(BaseModel):
    name: str
    tenant_id: str
    slack_team_ref: str


class SlackCallBackEventRepr(BaseModel):
    event_id: str
    is_ack: bool


class InSyncSlackChannelRepr(BaseModel):
    id: str
    name: str
    created: int
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
    topic: dict
    updated: int
    updated_at: datetime
    created_at: datetime


class InSyncSlackUserRepr(BaseModel):
    id: str
    name: str
    updated_at: datetime
    created_at: datetime


class UpsertUserRepr(BaseModel):
    user_id: str
    name: str
    role: str


class InSyncSlackUserWithUpsertedUserRepr(InSyncSlackUserRepr):
    user: UpsertUserRepr | None = None


class TriageSlackChannelRepr(BaseModel):
    slack_channel_ref: str
    slack_channel_name: str


class SlackChannelRepr(BaseModel):
    channel_id: str
    slack_channel_ref: str
    slack_channel_name: str
    triage_channel: TriageSlackChannelRepr


class IssueRepr(BaseModel):
    tenant_id: str
    issue_id: str
    issue_number: int
    slack_channel_id: str
    slack_message_ts: str
    body: str
    status: str
    priority: int
    tags: List[str] = []


class UserRepr(BaseModel):
    user_id: str
    slack_user_ref: str
    name: str
    role: str


def slack_callback_event_repr(
    slack_event: SlackEvent,
) -> SlackCallBackEventRepr:
    return SlackCallBackEventRepr(
        event_id=slack_event.event_id,
        is_ack=slack_event.is_ack,
    )


def insync_slack_channel_repr(
    item: InSyncSlackChannel,
) -> InSyncSlackChannelRepr:
    topic = {
        "value": item.topic.get("value", ""),
    }
    return InSyncSlackChannelRepr(
        id=item.id,
        name=item.name,
        created=item.created,
        is_archived=item.is_archived,
        is_channel=item.is_channel,
        is_ext_shared=item.is_ext_shared,
        is_general=item.is_general,
        is_group=item.is_group,
        is_im=item.is_im,
        is_member=item.is_member,
        is_mpim=item.is_mpim,
        is_org_shared=item.is_org_shared,
        is_pending_ext_shared=item.is_pending_ext_shared,
        is_private=item.is_private,
        is_shared=item.is_shared,
        topic=topic,
        updated=item.updated,
        updated_at=item.updated_at,
        created_at=item.created_at,
    )


def insync_slack_user_repr(item: InSyncSlackUser) -> dict:
    return InSyncSlackUserRepr(
        id=item.id,
        name=item.real_name,
        updated_at=item.updated_at,
        created_at=item.created_at,
    )


def insync_slack_user_with_upsert(item: Dict[str, User | InSyncSlackUser]) -> dict:
    insync_user: InSyncSlackUser = item["insync_user"]
    user: User | None = item["user"]
    if user:
        user_repr = UpsertUserRepr(
            user_id=user.user_id,
            name=user.name,
            role=user.role,
        )
    else:
        user_repr = None
    return InSyncSlackUserWithUpsertedUserRepr(
        id=insync_user.id,
        name=insync_user.name,
        is_bot=insync_user.is_bot,
        updated_at=insync_user.updated_at,
        created_at=insync_user.created_at,
        user=user_repr,
    )


def slack_channel_repr(item: SlackChannel) -> SlackChannelRepr:
    triage_channel = item.triage_channel
    triage_slack_channel = TriageSlackChannelRepr(
        slack_channel_ref=triage_channel.slack_channel_ref,
        slack_channel_name=triage_channel.slack_channel_name,
    )
    return SlackChannelRepr(
        channel_id=item.slack_channel_id,
        slack_channel_ref=item.slack_channel_ref,
        slack_channel_name=item.slack_channel_name,
        triage_channel=triage_slack_channel,
    )


def issue_repr(item: Issue) -> IssueRepr:
    return IssueRepr(
        tenant_id=item.tenant_id,
        issue_id=item.issue_id,
        issue_number=item.issue_number,
        slack_channel_id=item.slack_channel_id,
        slack_message_ts=item.slack_message_ts,
        body=item.body,
        status=item.status,
        priority=item.priority,
        tags=item.tags,
    )


def user_repr(item: User) -> UpsertUserRepr:
    return UserRepr(
        user_id=item.user_id,
        slack_user_ref=item.slack_user_ref,
        name=item.display_name,
        role=item.role,
    )
