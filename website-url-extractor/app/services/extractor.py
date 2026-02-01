"""DSPy-based extraction for project and event information."""

from dataclasses import dataclass, field
from typing import Optional

import dspy


class ExtractProjectInfo(dspy.Signature):
    """Extract project/organization information from website content."""

    page_content: str = dspy.InputField(desc="Combined markdown content from website pages")
    page_url: str = dspy.InputField(desc="The main URL of the website")

    title: str = dspy.OutputField(desc="Name of the project, group, or organization")
    description: str = dspy.OutputField(
        desc="4-6 sentence description of what this project/org does, its goals, "
        "target audience, and what makes it unique. Should be informative, engaging, "
        "and suitable for a project database entry."
    )
    homepage: str = dspy.OutputField(
        desc="Main website URL (use the provided URL if no better homepage found)"
    )
    keywords: list[str] = dspy.OutputField(
        desc="5-10 relevant keywords/topics for this project (in German if content is German)"
    )
    location: Optional[str] = dspy.OutputField(
        desc="City or region where this project is based, if mentioned"
    )
    contact_email: Optional[str] = dspy.OutputField(
        desc="Contact email address if found"
    )


class ExtractEvents(dspy.Signature):
    """Extract event information from page content, detecting recurrence patterns."""

    page_content: str = dspy.InputField(desc="Markdown content that may contain event information")

    events: list[dict] = dspy.OutputField(
        desc="List of events, each with: "
        "name (event name/title), "
        "start_date (ISO datetime like 2024-03-15T19:00:00), "
        "end_date (ISO datetime or null), "
        "recurrence_type (none/weekly/monthly-date/monthly-day), "
        "recurrence_day (day of week for weekly, e.g. 'wednesday'), "
        "recurrence_week (week of month for monthly-day, e.g. 'first', 'second', 'third', 'fourth', 'last'), "
        "description (brief event description)"
    )
    recurrence_reasoning: str = dspy.OutputField(
        desc="Explanation of detected patterns, e.g. 'jeden Mittwoch' → weekly on wednesday, "
        "'jeden ersten Montag im Monat' → monthly-day on first monday"
    )


@dataclass
class ExtractedProject:
    """Extracted project information."""

    title: str = ""
    description: str = ""
    homepage: str = ""
    keywords: list[str] = field(default_factory=list)
    location: Optional[str] = None
    contact_email: Optional[str] = None


@dataclass
class ExtractedEvent:
    """Extracted event information."""

    name: str = ""
    start_date: Optional[str] = None
    end_date: Optional[str] = None
    recurrence_type: str = "none"
    recurrence_day: Optional[str] = None
    recurrence_week: Optional[str] = None
    description: str = ""


@dataclass
class ExtractionResult:
    """Result of extracting project and event data."""

    project: ExtractedProject = field(default_factory=ExtractedProject)
    events: list[ExtractedEvent] = field(default_factory=list)
    recurrence_reasoning: str = ""


class ProjectExtractor:
    """Extracts structured project and event data using DSPy."""

    def __init__(self):
        self.project_extractor = dspy.Predict(ExtractProjectInfo)
        self.event_extractor = dspy.Predict(ExtractEvents)

    def extract(
        self,
        content: str,
        url: str,
    ) -> ExtractionResult:
        """Extract project and event information from content."""
        result = ExtractionResult()

        # Limit content size for LLM
        content_truncated = content[:15000]

        # Extract project info
        try:
            project_result = self.project_extractor(
                page_content=content_truncated,
                page_url=url,
            )

            result.project = ExtractedProject(
                title=project_result.title,
                description=project_result.description,
                homepage=project_result.homepage or url,
                keywords=project_result.keywords or [],
                location=project_result.location,
                contact_email=project_result.contact_email,
            )
        except Exception as e:
            # Fallback to minimal project info
            result.project = ExtractedProject(
                title="",
                description="",
                homepage=url,
            )

        # Extract events
        try:
            event_result = self.event_extractor(page_content=content_truncated)

            result.recurrence_reasoning = event_result.recurrence_reasoning

            for event_data in event_result.events or []:
                if isinstance(event_data, dict):
                    result.events.append(
                        ExtractedEvent(
                            name=event_data.get("name", ""),
                            start_date=event_data.get("start_date"),
                            end_date=event_data.get("end_date"),
                            recurrence_type=event_data.get("recurrence_type", "none"),
                            recurrence_day=event_data.get("recurrence_day"),
                            recurrence_week=event_data.get("recurrence_week"),
                            description=event_data.get("description", ""),
                        )
                    )
        except Exception:
            # Events are optional, continue without them
            pass

        return result
