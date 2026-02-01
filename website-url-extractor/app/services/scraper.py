"""Web scraping service with retry logic and link extraction."""

import re
from dataclasses import dataclass, field
from urllib.parse import urljoin, urlparse

import html2text
import httpx
from bs4 import BeautifulSoup

from ..config import get_settings


@dataclass
class ScrapedPage:
    """Result of scraping a single page."""

    url: str
    title: str
    markdown_content: str
    internal_links: list[str] = field(default_factory=list)
    meta_description: str = ""
    success: bool = True
    error: str = ""


class ProjectScraper:
    """Scrapes web pages and converts to structured content."""

    def __init__(self):
        settings = get_settings()
        self.timeout = settings.scraper_timeout
        self.max_retries = settings.scraper_max_retries
        self.user_agent = settings.scraper_user_agent

        # Configure html2text
        self.html_converter = html2text.HTML2Text()
        self.html_converter.ignore_links = False
        self.html_converter.ignore_images = False
        self.html_converter.ignore_emphasis = False
        self.html_converter.body_width = 0  # No wrapping

    async def scrape(self, url: str) -> ScrapedPage:
        """Scrape a URL with retry logic."""
        last_error = ""

        for attempt in range(self.max_retries):
            try:
                return await self._do_scrape(url)
            except httpx.TimeoutException:
                last_error = f"Timeout after {self.timeout}s"
            except httpx.HTTPStatusError as e:
                last_error = f"HTTP {e.response.status_code}"
                if e.response.status_code in (404, 403, 410):
                    break  # Don't retry for these
            except Exception as e:
                last_error = str(e)

        return ScrapedPage(
            url=url,
            title="",
            markdown_content="",
            success=False,
            error=last_error,
        )

    async def _do_scrape(self, url: str) -> ScrapedPage:
        """Perform the actual scraping."""
        async with httpx.AsyncClient(
            timeout=self.timeout,
            follow_redirects=True,
            headers={"User-Agent": self.user_agent},
        ) as client:
            response = await client.get(url)
            response.raise_for_status()

            # Parse HTML
            soup = BeautifulSoup(response.text, "lxml")

            # Extract title
            title = ""
            if soup.title:
                title = soup.title.string or ""
            if not title:
                h1 = soup.find("h1")
                if h1:
                    title = h1.get_text(strip=True)

            # Extract meta description
            meta_desc = ""
            meta_tag = soup.find("meta", attrs={"name": "description"})
            if meta_tag and meta_tag.get("content"):
                meta_desc = meta_tag["content"]

            # Remove unwanted elements before conversion
            for element in soup.find_all(
                ["script", "style", "nav", "footer", "header", "aside"]
            ):
                element.decompose()

            # Convert to markdown
            markdown = self.html_converter.handle(str(soup))
            markdown = self._clean_markdown(markdown)

            # Extract internal links
            base_url = f"{urlparse(url).scheme}://{urlparse(url).netloc}"
            internal_links = self._extract_internal_links(soup, url, base_url)

            return ScrapedPage(
                url=url,
                title=title.strip(),
                markdown_content=markdown,
                internal_links=internal_links,
                meta_description=meta_desc,
            )

    def _clean_markdown(self, markdown: str) -> str:
        """Clean up the converted markdown."""
        # Remove excessive newlines
        markdown = re.sub(r"\n{3,}", "\n\n", markdown)
        # Remove empty links
        markdown = re.sub(r"\[]\([^)]*\)", "", markdown)
        # Clean up whitespace
        markdown = markdown.strip()
        return markdown

    def _extract_internal_links(
        self, soup: BeautifulSoup, current_url: str, base_url: str
    ) -> list[str]:
        """Extract internal links from the page."""
        links = set()
        current_domain = urlparse(current_url).netloc

        for a_tag in soup.find_all("a", href=True):
            href = a_tag["href"]

            # Skip anchors, javascript, mailto, tel
            if href.startswith(("#", "javascript:", "mailto:", "tel:")):
                continue

            # Resolve relative URLs
            full_url = urljoin(current_url, href)
            parsed = urlparse(full_url)

            # Only include internal links
            if parsed.netloc == current_domain:
                # Normalize URL (remove fragments, trailing slashes)
                normalized = f"{parsed.scheme}://{parsed.netloc}{parsed.path}"
                if parsed.query:
                    normalized += f"?{parsed.query}"
                normalized = normalized.rstrip("/")

                # Skip if same as current page
                current_normalized = current_url.rstrip("/")
                if normalized != current_normalized:
                    links.add(normalized)

        return list(links)

