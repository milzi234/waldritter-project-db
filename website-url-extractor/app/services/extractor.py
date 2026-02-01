"""DSPy-based extraction for project and event information."""

import json
import logging
from dataclasses import dataclass, field
from typing import Optional

import dspy

logger = logging.getLogger(__name__)


class ExtractProjectInfo(dspy.Signature):
    """Extract project/organization information from website content. Output in German."""

    page_content: str = dspy.InputField(desc="Combined markdown content from website pages")
    page_url: str = dspy.InputField(desc="The main URL of the website")

    title: str = dspy.OutputField(desc="Name des Projekts, der Gruppe oder Organisation auf Deutsch")
    description: str = dspy.OutputField(
        desc="4-6 Sätze auf Deutsch, die beschreiben was dieses Projekt/Organisation macht, "
        "ihre Ziele, Zielgruppe und was sie einzigartig macht. Sollte informativ, ansprechend "
        "und geeignet für einen Projektdatenbank-Eintrag sein. IMMER auf Deutsch schreiben."
    )
    homepage: str = dspy.OutputField(
        desc="Haupt-Website-URL (verwende die angegebene URL falls keine bessere Homepage gefunden wird)"
    )
    keywords: list[str] = dspy.OutputField(
        desc="5-10 relevante Schlüsselwörter/Themen für dieses Projekt, auf Deutsch"
    )
    location: Optional[str] = dspy.OutputField(
        desc="Stadt oder Region wo dieses Projekt ansässig ist, falls erwähnt"
    )
    contact_email: Optional[str] = dspy.OutputField(
        desc="Kontakt-E-Mail-Adresse falls gefunden"
    )


class ExtractEvents(dspy.Signature):
    """Extract event/date information from page content. ALWAYS return events if any dates are mentioned."""

    page_content: str = dspy.InputField(desc="Markdown content that may contain event information")

    events: list[dict] = dspy.OutputField(
        desc="WICHTIG: Wenn Daten/Termine im Text erwähnt werden, MUSS mindestens ein Event zurückgegeben werden. "
        "Liste von Terminen als JSON array, jeweils mit: "
        "name (string: Terminname/-titel auf Deutsch), "
        "start_date (string: ISO datetime wie '2025-10-02T18:00:00'), "
        "end_date (string: ISO datetime wie '2025-10-05T13:00:00' oder null), "
        "recurrence_type (string: 'none' für einmalig, 'weekly', 'monthly-date', 'monthly-day'), "
        "recurrence_day (string oder null: Wochentag für weekly), "
        "recurrence_week (string oder null: 'first'/'second'/'third'/'fourth'/'last' für monthly-day), "
        "description (string: kurze Terminbeschreibung auf Deutsch). "
        "Beispiel: [{'name': 'LARP Event', 'start_date': '2025-10-02T18:00:00', 'end_date': '2025-10-05T13:00:00', "
        "'recurrence_type': 'none', 'recurrence_day': null, 'recurrence_week': null, 'description': 'Hauptevent'}]"
    )
    recurrence_reasoning: str = dspy.OutputField(
        desc="Erklärung der erkannten Termine und Muster auf Deutsch"
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

            # Handle keywords - DSPy sometimes returns a string instead of list
            keywords = project_result.keywords or []
            if isinstance(keywords, str):
                # Split comma-separated string into list
                keywords = [k.strip() for k in keywords.split(",") if k.strip()]

            result.project = ExtractedProject(
                title=project_result.title,
                description=project_result.description,
                homepage=project_result.homepage or url,
                keywords=keywords,
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

            result.recurrence_reasoning = event_result.recurrence_reasoning or ""

            # Handle events - DSPy might return string, list, or other formats
            events_data = event_result.events
            logger.info(f"Raw events data type: {type(events_data)}, value: {events_data}")

            # If it's a string, try to parse as JSON
            if isinstance(events_data, str):
                try:
                    events_data = json.loads(events_data)
                except json.JSONDecodeError:
                    logger.warning(f"Could not parse events as JSON: {events_data}")
                    events_data = []

            for event_data in events_data or []:
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
                    logger.info(f"Added event: {event_data.get('name')}")
        except Exception as e:
            logger.exception(f"Event extraction failed: {e}")
            # Events are optional, continue without them

        return result
