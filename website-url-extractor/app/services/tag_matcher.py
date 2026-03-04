"""Tag matching using keyword matching and LLM suggestions."""

from dataclasses import dataclass

import dspy

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
    match_type: str  # 'keyword', 'llm'


@dataclass
class TagMatchResult:
    """Result of tag matching."""

    matches: list[TagMatch]
    reasoning: str = ""


class TagMatcher:
    """Matches project content to available tags using keyword + LLM approach."""

    def __init__(self):
        self.suggester = dspy.Predict(SuggestTags)

    def match(
        self,
        title: str,
        description: str,
        keywords: list[str],
        available_tags: list[dict],
    ) -> TagMatchResult:
        """Match project to available tags using keyword and LLM scoring."""
        if not available_tags:
            return TagMatchResult(matches=[])

        tag_names = [t["name"] for t in available_tags]
        tag_lookup = {t["name"]: t for t in available_tags}

        matches_dict: dict[str, TagMatch] = {}

        # Keyword matches
        keywords_lower = [k.lower() for k in keywords]
        for tag_name in tag_names:
            tag_lower = tag_name.lower()
            for keyword in keywords_lower:
                if keyword in tag_lower or tag_lower in keyword:
                    tag_info = tag_lookup[tag_name]
                    matches_dict[tag_name] = TagMatch(
                        tag_id=tag_info["id"],
                        tag_name=tag_name,
                        category_name=tag_info.get("category_name", ""),
                        score=0.5,
                        match_type="keyword",
                    )
                    break

        # LLM suggestions
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
                        matches_dict[suggested_tag].score = min(
                            1.0, matches_dict[suggested_tag].score + 0.3
                        )
                        matches_dict[suggested_tag].match_type = "keyword+llm"
                    else:
                        matches_dict[suggested_tag] = TagMatch(
                            tag_id=tag_info["id"],
                            tag_name=suggested_tag,
                            category_name=tag_info.get("category_name", ""),
                            score=0.7,
                            match_type="llm",
                        )
        except Exception:
            pass

        matches = sorted(matches_dict.values(), key=lambda m: m.score, reverse=True)
        return TagMatchResult(matches=matches, reasoning=reasoning)
