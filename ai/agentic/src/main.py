from restate import Service, Context
import restate
from pydantic import BaseModel
import uuid
import logfire

logfire.configure()

class ThreadPingReq(BaseModel):
    threadId: str


class ThreadPing(BaseModel):
    message: str


thread = Service("thread")


@thread.handler()
async def ping(ctx: Context, req: ThreadPingReq) -> ThreadPing:
    """Handle ping requests for a thread.

    Args:
        ctx: Restate context
        req: Thread ping request containing threadId

    Returns:
        ThreadPing response with confirmation message
    """
    logfire.info('thread service PING for {threadId}', threadId=req.threadId)
    
    ping_id = await ctx.run("generating ping UUID", lambda: str(uuid.uuid4()))
    print(ping_id)
    return ThreadPing(
        message=f"Got PING for threadId: {req.threadId} with PING ID: {ping_id}"
    )


app = restate.app(services=[thread])
