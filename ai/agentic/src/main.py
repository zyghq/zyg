import uuid

import logfire
import restate
from pydantic import BaseModel
from restate import Context, Service
from restate.exceptions import TerminalError

from agents.summarizer.errors import PromptError, SummarizerError
from agents.summarizer.summarize import (
    SummarizeRequest,
    SummarizeResponse,
    summarize_text_fn,
)

logfire.configure()


class ThreadPingReq(BaseModel):
    threadId: str


class ThreadPingResponse(BaseModel):
    message: str


defaults = Service("defaults")

thread = Service("thread")


@thread.handler(name="ping")
async def ping(ctx: Context, req: ThreadPingReq) -> ThreadPingResponse:
    """Handle ping requests for a thread.

    Args:
        ctx: Restate context
        req: Thread ping request containing threadId

    Returns:
        ThreadPing response with confirmation message
    """
    logfire.info("thread service PING for {threadId}", threadId=req.threadId)

    ping_id = await ctx.run("generating ping UUID", lambda: str(uuid.uuid4()))
    return ThreadPingResponse(
        message=f"Got PING for threadId: {req.threadId} with PING ID: {ping_id}"
    )


@defaults.handler(name="summarize")
async def summarize(ctx: Context, req: SummarizeRequest) -> SummarizeResponse:
    """Handle text summarization requests.

    This handler processes incoming summarization requests by calling the summarize_text function.
    It handles errors appropriately to work with Restate's retry mechanism.

    Args:
        ctx: Restate context for the request
        req: SummarizeRequest containing the text to summarize and optional parameters
            - text: The input text to summarize
            - max_words: Optional maximum word count for the summary
            - context: Optional additional context
            - prompt_version: Optional specific prompt version to use

    Returns:
        SummarizeResponse containing:
            - summary: The generated summary text
            - original_length: Word count of input text
            - summary_length: Word count of summary
            - inferred_subject: Detected subject/title
            - prompt_version_used: Which prompt version was used
            - metadata: Additional summary-specific data

    Raises:
        TerminalError: For unrecoverable errors (invalid prompts, malformed input).
            These errors are terminal and will not be retried by Restate.
            - PromptError: For invalid prompt configurations
            - SummarizerError: For summarization processing errors

    Note:
        Restate implements infinite retries with exponential backoff for non-terminal errors.
        We use TerminalError for cases where retrying would not help (invalid input, etc.)
        to immediately fail the request and propagate the error to the caller.
    """
    try:
        result = await summarize_text_fn(req)
        return result
    except PromptError as e:
        logfire.error(e)
        raise TerminalError(e.detail, status_code=400)
    except SummarizerError as e:
        logfire.error(e)
        raise TerminalError(e.detail, status_code=500)


app = restate.app(services=[thread, defaults])
