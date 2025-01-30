import logging
import os
import random
from enum import IntEnum
from typing import Dict, List, Optional

from pydantic import BaseModel, Field
from pydantic_ai import Agent

from .errors import PromptError, SummarizerError
from .prompts import (
    BULLET_POINT_PROMPT,
    ENTITY_PROMPT,
    EXECUTIVE_PROMPT,
    PROMPT_FUNCTIONS,
    QA_PROMPT,
    BulletPointSummary,
    Entity,
    ExecutiveSummary,
    QuestionAnswerSummary,
)

API_KEY = os.getenv("GEMINI_API_KEY")
if not API_KEY:
    raise ValueError("GEMINI_API_KEY environment variable is not set")


class PromptVersion(IntEnum):
    EXECUTIVE = EXECUTIVE_PROMPT
    ENTITY = ENTITY_PROMPT
    QA = QA_PROMPT
    BULLET_POINT = BULLET_POINT_PROMPT

    @classmethod
    def get_description(cls) -> str:
        return ", ".join(f"{member.value} ({member.name})" for member in cls)


class SummarizeRequest(BaseModel):
    text: str = Field(
        ...,
        min_length=10,
        max_length=50000,  # Add reasonable maximum length
        description="The text to summarize",
    )
    max_words: Optional[int] = Field(
        default=100,
        gt=0,
        le=1000,  # Add reasonable maximum
        description="Maximum words in the summary",
    )
    context: Optional[str] = Field(
        default=None,
        max_length=1000,  # Add reasonable maximum
        description="Additional context about the text",
    )
    prompt_version: Optional[PromptVersion] = Field(
        default=None,
        description=f"Specify prompt version: {PromptVersion.get_description()}",
    )


class SummarizeResponse(BaseModel):
    summary: str
    original_length: int
    summary_length: int
    inferred_subject: str
    prompt_version_used: int
    metadata: Optional[dict] = Field(
        default=None, description="Metadata about the summary"
    )


class CustomerInfo(BaseModel):
    customerId: str
    name: str


class MemberInfo(BaseModel):
    memberId: str
    name: str


class Message(BaseModel):
    messageId: str
    textBody: str
    customer: Optional[CustomerInfo] = None
    member: Optional[MemberInfo] = None


class Label(BaseModel):
    labelId: str
    name: str


class SupportThread(BaseModel):
    threadId: str
    title: str
    customer: CustomerInfo
    channel: str
    messages: List[Message]
    labels: List[Label]
    priority: str
    createdAt: str
    updatedAt: str


class SummarizeThreadRequest(BaseModel):
    thread: SupportThread
    prompt_version: Optional[int] = Field(
        default=BULLET_POINT_PROMPT,
        description=f"Specify prompt version: {', '.join(str(k) + ' (' + v.__name__.split('_')[1] + ')' for k, v in PROMPT_FUNCTIONS.items())}",
    )


# Create agents for each summary type
executive_agent = Agent("gemini-1.5-pro-latest", result_type=ExecutiveSummary)
entity_agent = Agent("gemini-1.5-pro-latest", result_type=Entity)
qa_agent = Agent("gemini-1.5-pro-latest", result_type=QuestionAnswerSummary)
bullet_point_agent = Agent("gemini-1.5-pro-latest", result_type=BulletPointSummary)

# Type hint for the agent mapping
AGENTS: Dict[int, Agent] = {
    PromptVersion.EXECUTIVE.value: executive_agent,
    PromptVersion.ENTITY.value: entity_agent,
    PromptVersion.QA.value: qa_agent,
    PromptVersion.BULLET_POINT.value: bullet_point_agent,
}


async def summarize_text_fn(request: SummarizeRequest) -> SummarizeResponse:
    logger = logging.getLogger(__name__)

    try:
        # Determine which prompt version to use
        prompt_version = (
            request.prompt_version.value
            if request.prompt_version
            else random.choice([v.value for v in PromptVersion])
        )

        # Get the prompt function and agent
        prompt_function = PROMPT_FUNCTIONS.get(prompt_version)
        agent = AGENTS.get(prompt_version)

        if not prompt_function or not agent:
            raise PromptError(
                detail=f"Invalid prompt version: {prompt_version}. Available versions: {PromptVersion.get_description()}"
            )

        logger.info(
            f"Generating summary using prompt version: {PromptVersion(prompt_version).name}"
        )

        # Create structured prompt
        prompt = prompt_function(
            text=request.text, max_words=request.max_words, context=request.context
        )

        # Generate summary using the appropriate agent
        result = await agent.run(prompt)
        response_data = result.data

        # Extract common fields
        title = response_data.title
        summary = response_data.summary.strip()

        # Calculate lengths
        original_length = len(request.text.split())
        summary_length = len(summary.split())

        if summary_length > request.max_words:
            logger.warning(
                f"Summary exceeds max_words limit: {summary_length} > {request.max_words}"
            )

        # Construct metadata based on prompt version
        metadata = {
            PromptVersion.EXECUTIVE.value: lambda: {"audience": response_data.audience},
            PromptVersion.ENTITY.value: lambda: {"entities": response_data.entities},
            PromptVersion.QA.value: lambda: {"qa_pairs": response_data.qa_pairs},
            PromptVersion.BULLET_POINT.value: lambda: {
                "bullet_points": response_data.bullet_points
            },
        }.get(prompt_version, lambda: {})()

        return SummarizeResponse(
            summary=summary,
            original_length=original_length,
            summary_length=summary_length,
            inferred_subject=title,
            prompt_version_used=prompt_version,
            metadata=metadata,
        )
    except PromptError:
        raise
    except Exception as e:
        logger.error(f"Error generating summary: {str(e)}", exc_info=True)
        raise SummarizerError(detail=f"Error generating summary: {str(e)}")
