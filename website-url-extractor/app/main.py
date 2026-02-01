"""FastAPI application for URL extraction service."""

from contextlib import asynccontextmanager

import dspy
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from .api import router
from .config import get_settings


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan handler for startup/shutdown."""
    settings = get_settings()

    # Configure DSPy with Anthropic
    if settings.anthropic_api_key:
        lm = dspy.LM(
            model=f"anthropic/{settings.dspy_lm_model}",
            api_key=settings.anthropic_api_key,
        )
        dspy.configure(lm=lm)

    yield


app = FastAPI(
    title="Waldritter URL Extractor",
    description="Extract project and event data from URLs using DSPy",
    version="1.0.0",
    lifespan=lifespan,
)

# CORS middleware for local development
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173", "http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include API routes
app.include_router(router, prefix="/api")


@app.get("/health")
async def health_check():
    """Health check endpoint."""
    return {"status": "healthy"}
