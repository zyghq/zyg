from fastapi import FastAPI
from sqlalchemy.sql import text

from src.adapters.db import engine
from src.logger import logger

from .routers import events, interactions, issues, onboardings, tenants

app = FastAPI()


app.include_router(
    events.router,
    prefix="/events",
)

app.include_router(
    interactions.router,
    prefix="/interactions",
)

app.include_router(
    onboardings.router,
    prefix="/onboardings",
)

app.include_router(
    tenants.router,
    prefix="/tenants",
)

app.include_router(
    issues.router,
    prefix="/issues",
)


@app.get("/")
async def root():
    logger.info("Hey there! I am zyg.")
    return {"message": "Hey there! I am zyg."}


@app.on_event("startup")
async def startup():
    async with engine.begin() as conn:
        query = text("SELECT NOW()::timestamp AS now")
        rows = await conn.execute(query)
        result = rows.mappings().first()
        logger.info(f"db connected at: {result['now']}")


@app.on_event("shutdown")
async def shutdown():
    logger.warning("cleaning up...")
    await engine.dispose()
