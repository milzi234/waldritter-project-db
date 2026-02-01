"""Agentic site exploration using LLM-driven decisions."""

from dataclasses import dataclass, field
from typing import Optional

import dspy

from ..config import get_settings
from .scraper import ProjectScraper, ScrapedPage


class AnalyzePageContent(dspy.Signature):
    """Analyze what project information is available on this page."""

    page_content: str = dspy.InputField(desc="Markdown content of the current page")
    gathered_so_far: str = dspy.InputField(
        desc="Summary of info collected from previous pages"
    )

    has_title: bool = dspy.OutputField(desc="Page has a clear project/organization title")
    has_description: bool = dspy.OutputField(
        desc="Page has a meaningful description of what this project/org does"
    )
    has_events: bool = dspy.OutputField(
        desc="Page mentions events, meetings, or recurring activities"
    )
    has_contact: bool = dspy.OutputField(
        desc="Page has contact information (email, address, phone)"
    )
    extracted_info: str = dspy.OutputField(
        desc="Key information extracted from this page (title, description snippets, event mentions)"
    )
    missing_info: list[str] = dspy.OutputField(
        desc="What key info is still missing for a complete project entry"
    )
    confidence: float = dspy.OutputField(
        desc="0-1 confidence we have enough to create a good project entry"
    )


class PrioritizeLinks(dspy.Signature):
    """Decide which links are most likely to contain missing information."""

    available_links: list[str] = dspy.InputField(desc="List of internal URLs to explore")
    missing_info: list[str] = dspy.InputField(desc="What information we're looking for")
    page_context: str = dspy.InputField(desc="Brief context about the site/project")

    ranked_links: list[str] = dspy.OutputField(
        desc="Top 3 links most likely to have missing info, ordered by priority"
    )
    reasoning: str = dspy.OutputField(desc="Why these links were prioritized")


class DecideExploration(dspy.Signature):
    """Decide whether to continue exploring or stop."""

    confidence: float = dspy.InputField(desc="Current confidence level 0-1")
    pages_visited: int = dspy.InputField(desc="Number of pages already scraped")
    max_pages: int = dspy.InputField(desc="Maximum pages allowed")
    missing_info: list[str] = dspy.InputField(desc="What info is still missing")
    has_promising_links: bool = dspy.InputField(desc="Whether there are unvisited promising links")

    should_continue: bool = dspy.OutputField(
        desc="True if we should explore more pages"
    )
    reason: str = dspy.OutputField(desc="Explanation of the decision")


@dataclass
class ExplorationResult:
    """Result of exploring a site."""

    pages: list[ScrapedPage] = field(default_factory=list)
    combined_content: str = ""
    extracted_info: str = ""
    confidence: float = 0.0
    exploration_log: list[str] = field(default_factory=list)


class SiteExplorer:
    """Explores a site using LLM-driven decisions."""

    def __init__(self, scraper: Optional[ProjectScraper] = None):
        settings = get_settings()
        self.max_pages = settings.explorer_max_pages
        self.max_depth = settings.explorer_max_depth
        self.min_confidence = settings.explorer_min_confidence
        self.scraper = scraper or ProjectScraper()

        # Initialize DSPy predictors
        self.analyzer = dspy.Predict(AnalyzePageContent)
        self.link_prioritizer = dspy.Predict(PrioritizeLinks)
        self.decision_maker = dspy.Predict(DecideExploration)

    async def explore(self, start_url: str) -> ExplorationResult:
        """Explore a site starting from the given URL."""
        result = ExplorationResult()
        visited_urls: set[str] = set()
        pending_links: list[tuple[str, int]] = [(start_url, 0)]  # (url, depth)
        gathered_info = ""

        while pending_links and len(result.pages) < self.max_pages:
            url, depth = pending_links.pop(0)

            # Skip if already visited or too deep
            if url in visited_urls or depth > self.max_depth:
                continue

            visited_urls.add(url)
            result.exploration_log.append(f"Scraping: {url} (depth {depth})")

            # Scrape the page
            page = await self.scraper.scrape(url)
            if not page.success:
                result.exploration_log.append(f"  Failed: {page.error}")
                continue

            result.pages.append(page)

            # Truncate content for LLM analysis
            content_for_analysis = page.markdown_content[:8000]

            # Analyze what we found
            try:
                analysis = self.analyzer(
                    page_content=content_for_analysis,
                    gathered_so_far=gathered_info,
                )

                result.exploration_log.append(
                    f"  Confidence: {analysis.confidence:.2f}, "
                    f"Missing: {', '.join(analysis.missing_info) or 'nothing'}"
                )

                # Update gathered info
                gathered_info = analysis.extracted_info
                result.extracted_info = gathered_info
                result.confidence = analysis.confidence

                # Check if we should continue
                decision = self.decision_maker(
                    confidence=analysis.confidence,
                    pages_visited=len(result.pages),
                    max_pages=self.max_pages,
                    missing_info=analysis.missing_info,
                    has_promising_links=bool(page.internal_links),
                )

                if not decision.should_continue:
                    result.exploration_log.append(f"  Stopping: {decision.reason}")
                    break

                # Prioritize links if we need more info
                if page.internal_links and analysis.missing_info:
                    prioritized = self.link_prioritizer(
                        available_links=page.internal_links[:20],  # Limit for LLM
                        missing_info=analysis.missing_info,
                        page_context=page.title or start_url,
                    )

                    # Add prioritized links to pending
                    for link in prioritized.ranked_links[:3]:
                        if link not in visited_urls:
                            pending_links.append((link, depth + 1))

                    result.exploration_log.append(
                        f"  Prioritized links: {prioritized.ranked_links[:3]}"
                    )

            except Exception as e:
                result.exploration_log.append(f"  Analysis error: {e}")
                # Continue with remaining pages even if analysis fails

        # Combine all page content
        result.combined_content = self._combine_pages(result.pages)

        return result

    def _combine_pages(self, pages: list[ScrapedPage]) -> str:
        """Combine content from all scraped pages."""
        sections = []

        for page in pages:
            section = f"# {page.title or page.url}\n\n{page.markdown_content}"
            sections.append(section)

        return "\n\n---\n\n".join(sections)
