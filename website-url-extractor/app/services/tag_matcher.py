"""Semantic tag matching using embeddings and LLM suggestions."""

from dataclasses import dataclass
from typing import Optional

import dspy
from sentence_transformers import SentenceTransformer

from ..config import get_settings


class SuggestTags(dspy.Signature):
    """Suggest relevant tags from available options based on project content."""

    project_title: str = dspy.InputField(desc="Title of the project")
    project_description: str = dspy.InputField(desc="Description of the project")
    project_keywords: list[str] = dspy.InputField(desc="Keywords extracted from project")
    available_tags: list[str] = dspy.InputField(desc="List of available tag names to choose from")

    suggested_tags: list[str] = dspy.OutputField(
        desc="List of tag names from available_tags that best match this project. "
        "Only include tags that are clearly relevant."
    )
    reasoning: str = dspy.OutputField(desc="Brief explanation of why these tags were chosen")


@dataclass
class TagMatch:
    """A matched tag with confidence score."""

    tag_id: int
    tag_name: str
    category_name: str
    score: float
    match_type: str  # 'semantic', 'keyword', 'llm'


@dataclass
class TagMatchResult:
    """Result of tag matching."""

    matches: list[TagMatch]
    reasoning: str = ""


class TagMatcher:
    """Matches project content to available tags using hybrid approach."""

    def __init__(self):
        settings = get_settings()
        self.threshold = settings.tag_match_threshold

        # Load embedding model (multilingual for German content)
        self.model = SentenceTransformer(settings.embedding_model)

        # DSPy predictor for LLM suggestions
        self.suggester = dspy.Predict(SuggestTags)

        # Cache for tag embeddings
        self._tag_embeddings = None
        self._tag_data = None

    def match(
        self,
        title: str,
        description: str,
        keywords: list[str],
        available_tags: list[dict],
    ) -> TagMatchResult:
        """Match project to available tags using hybrid scoring."""
        if not available_tags:
            return TagMatchResult(matches=[])

        # Prepare tag data
        tag_names = [t["name"] for t in available_tags]
        tag_lookup = {t["name"]: t for t in available_tags}

        # Compute embeddings for project content
        project_text = f"{title}. {description}. {', '.join(keywords)}"
        project_embedding = self.model.encode(project_text)

        # Compute embeddings for tags (with caching)
        tag_embeddings = self._get_tag_embeddings(tag_names)

        # Compute semantic similarity scores
        similarities = self.model.similarity(project_embedding, tag_embeddings)[0]

        # Build initial matches from semantic similarity
        matches_dict: dict[str, TagMatch] = {}

        for i, tag_name in enumerate(tag_names):
            score = float(similarities[i])
            if score >= self.threshold:
                tag_info = tag_lookup[tag_name]
                matches_dict[tag_name] = TagMatch(
                    tag_id=tag_info["id"],
                    tag_name=tag_name,
                    category_name=tag_info.get("category_name", ""),
                    score=score,
                    match_type="semantic",
                )

        # Boost scores for keyword matches
        keywords_lower = [k.lower() for k in keywords]
        for tag_name in tag_names:
            tag_lower = tag_name.lower()
            for keyword in keywords_lower:
                if keyword in tag_lower or tag_lower in keyword:
                    tag_info = tag_lookup[tag_name]
                    if tag_name in matches_dict:
                        # Boost existing match
                        matches_dict[tag_name].score = min(
                            1.0, matches_dict[tag_name].score + 0.2
                        )
                        matches_dict[tag_name].match_type = "keyword+semantic"
                    else:
                        # Add new keyword match
                        matches_dict[tag_name] = TagMatch(
                            tag_id=tag_info["id"],
                            tag_name=tag_name,
                            category_name=tag_info.get("category_name", ""),
                            score=0.5,
                            match_type="keyword",
                        )
                    break

        # Get LLM suggestions
        reasoning = ""
        try:
            llm_result = self.suggester(
                project_title=title,
                project_description=description,
                project_keywords=keywords,
                available_tags=tag_names,
            )

            reasoning = llm_result.reasoning

            for suggested_tag in llm_result.suggested_tags or []:
                if suggested_tag in tag_lookup:
                    tag_info = tag_lookup[suggested_tag]
                    if suggested_tag in matches_dict:
                        # Boost existing match
                        matches_dict[suggested_tag].score = min(
                            1.0, matches_dict[suggested_tag].score + 0.15
                        )
                        if "llm" not in matches_dict[suggested_tag].match_type:
                            matches_dict[suggested_tag].match_type += "+llm"
                    else:
                        # Add LLM suggestion
                        matches_dict[suggested_tag] = TagMatch(
                            tag_id=tag_info["id"],
                            tag_name=suggested_tag,
                            category_name=tag_info.get("category_name", ""),
                            score=0.6,
                            match_type="llm",
                        )
        except Exception:
            # LLM suggestions are optional
            pass

        # Sort by score and return
        matches = sorted(matches_dict.values(), key=lambda m: m.score, reverse=True)

        return TagMatchResult(matches=matches, reasoning=reasoning)

    def _get_tag_embeddings(self, tag_names: list[str]):
        """Get embeddings for tags, using cache if available."""
        # Simple cache - in production, consider more sophisticated caching
        cache_key = tuple(sorted(tag_names))
        if self._tag_data != cache_key:
            self._tag_embeddings = self.model.encode(tag_names)
            self._tag_data = cache_key
        return self._tag_embeddings
