"""API routes for the extraction service."""

import base64
from typing import Optional

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel, HttpUrl

from ..services import ProjectScraper, SiteExplorer, ProjectExtractor, TagMatcher, ImageGenerator

router = APIRouter()


class TagInput(BaseModel):
    """Tag data from Rails."""

    id: int
    name: str
    category_name: Optional[str] = None


class ExtractRequest(BaseModel):
    """Request to extract project data from a URL."""

    url: HttpUrl
    available_tags: list[TagInput] = []


class ExtractedEventResponse(BaseModel):
    """Extracted event in response."""

    name: str
    start_date: Optional[str] = None
    end_date: Optional[str] = None
    recurrence_type: str = "none"
    recurrence_day: Optional[str] = None
    recurrence_week: Optional[str] = None
    description: str = ""
    enabled: bool = True


class TagMatchResponse(BaseModel):
    """Tag match in response."""

    tag_id: int
    tag_name: str
    category_name: str
    score: float
    match_type: str


class ExtractResponse(BaseModel):
    """Response with extracted project data."""

    # Project info
    title: str
    description: str
    homepage: str
    keywords: list[str] = []
    location: Optional[str] = None
    contact_email: Optional[str] = None

    # Events
    events: list[ExtractedEventResponse] = []
    recurrence_reasoning: str = ""

    # Tag suggestions
    suggested_tags: list[TagMatchResponse] = []
    tag_reasoning: str = ""

    # Exploration metadata
    pages_explored: int = 0
    exploration_confidence: float = 0.0
    exploration_log: list[str] = []


# Image generation models
class GenerateImageRequest(BaseModel):
    """Request to generate a project thumbnail image."""

    title: str
    description: str
    keywords: list[str] = []
    variation: int = 0  # Variation number for different styles


class GenerateImageResponse(BaseModel):
    """Response with generated image."""

    image_base64: str
    prompt_used: str
    reasoning: str


# Initialize services (lazy loading)
_scraper: Optional[ProjectScraper] = None
_explorer: Optional[SiteExplorer] = None
_extractor: Optional[ProjectExtractor] = None
_tag_matcher: Optional[TagMatcher] = None
_image_generator: Optional[ImageGenerator] = None


def get_scraper() -> ProjectScraper:
    global _scraper
    if _scraper is None:
        _scraper = ProjectScraper()
    return _scraper


def get_explorer() -> SiteExplorer:
    global _explorer
    if _explorer is None:
        _explorer = SiteExplorer(get_scraper())
    return _explorer


def get_extractor() -> ProjectExtractor:
    global _extractor
    if _extractor is None:
        _extractor = ProjectExtractor()
    return _extractor


def get_tag_matcher() -> TagMatcher:
    global _tag_matcher
    if _tag_matcher is None:
        _tag_matcher = TagMatcher()
    return _tag_matcher


def get_image_generator() -> ImageGenerator:
    global _image_generator
    if _image_generator is None:
        _image_generator = ImageGenerator()
    return _image_generator


@router.post("/generate-image", response_model=GenerateImageResponse)
async def generate_image(request: GenerateImageRequest):
    """
    Generate a project thumbnail image using AI.

    This endpoint:
    1. Uses LLM to generate an optimized image prompt from project data
    2. Calls Gemini image generation to create the image
    3. Returns the image as base64-encoded PNG
    """
    generator = get_image_generator()

    try:
        image_bytes, prompt_used, reasoning = await generator.generate(
            title=request.title,
            description=request.description,
            keywords=request.keywords,
            variation=request.variation,
        )

        return GenerateImageResponse(
            image_base64=base64.b64encode(image_bytes).decode("utf-8"),
            prompt_used=prompt_used,
            reasoning=reasoning,
        )
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Image generation failed: {e}")


@router.post("/extract", response_model=ExtractResponse)
async def extract_from_url(request: ExtractRequest):
    """
    Extract project and event data from a URL.

    This endpoint:
    1. Explores the site using LLM-driven decisions
    2. Extracts project info (title, description, etc.)
    3. Extracts events with recurrence detection
    4. Matches content to available tags
    """
    url_str = str(request.url)

    # Step 1: Explore the site
    explorer = get_explorer()
    try:
        exploration = await explorer.explore(url_str)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Exploration failed: {e}")

    if not exploration.pages:
        raise HTTPException(
            status_code=400,
            detail="Could not scrape any pages from the provided URL",
        )

    # Step 2: Extract project and event info
    extractor = get_extractor()
    extraction = extractor.extract(
        content=exploration.combined_content,
        url=url_str,
    )

    # Step 3: Match tags
    tag_matches = []
    tag_reasoning = ""

    if request.available_tags:
        tag_matcher = get_tag_matcher()
        tags_for_matching = [
            {
                "id": t.id,
                "name": t.name,
                "category_name": t.category_name or "",
            }
            for t in request.available_tags
        ]

        match_result = tag_matcher.match(
            title=extraction.project.title,
            description=extraction.project.description,
            keywords=extraction.project.keywords,
            available_tags=tags_for_matching,
        )

        tag_matches = [
            TagMatchResponse(
                tag_id=m.tag_id,
                tag_name=m.tag_name,
                category_name=m.category_name,
                score=m.score,
                match_type=m.match_type,
            )
            for m in match_result.matches
        ]
        tag_reasoning = match_result.reasoning

    # Build response
    return ExtractResponse(
        title=extraction.project.title,
        description=extraction.project.description,
        homepage=extraction.project.homepage,
        keywords=extraction.project.keywords,
        location=extraction.project.location,
        contact_email=extraction.project.contact_email,
        events=[
            ExtractedEventResponse(
                name=e.name,
                start_date=e.start_date,
                end_date=e.end_date,
                recurrence_type=e.recurrence_type,
                recurrence_day=e.recurrence_day,
                recurrence_week=e.recurrence_week,
                description=e.description,
            )
            for e in extraction.events
        ],
        recurrence_reasoning=extraction.recurrence_reasoning,
        suggested_tags=tag_matches,
        tag_reasoning=tag_reasoning,
        pages_explored=len(exploration.pages),
        exploration_confidence=exploration.confidence,
        exploration_log=exploration.exploration_log,
    )
