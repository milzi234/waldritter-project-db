"""Configuration settings for the URL extractor service."""

from pydantic_settings import BaseSettings
from functools import lru_cache


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    # API Keys
    anthropic_api_key: str = ""
    google_api_key: str = ""

    # LLM Configuration
    dspy_lm_model: str = "claude-3-5-haiku-latest"

    # Image Generation Configuration
    image_model: str = "imagen-4.0-generate-001"
    image_aspect_ratio: str = "1:1"

    # Explorer Safety Limits
    explorer_max_pages: int = 8
    explorer_max_depth: int = 3
    explorer_min_confidence: float = 0.7

    # Scraper Configuration
    scraper_timeout: int = 30
    scraper_max_retries: int = 3
    scraper_user_agent: str = (
        "Mozilla/5.0 (compatible; WaldritterBot/1.0; +https://waldritter.de)"
    )

    # Tag Matcher Configuration
    embedding_model: str = "paraphrase-multilingual-MiniLM-L12-v2"
    tag_match_threshold: float = 0.3

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


@lru_cache
def get_settings() -> Settings:
    """Get cached settings instance."""
    return Settings()
