"""Image generation service using Google Gemini Imagen."""

import asyncio
import base64
import io
import logging
from typing import Optional

import dspy
from google import genai
from google.genai import types

from ..config import get_settings

logger = logging.getLogger(__name__)


class GenerateImagePrompt(dspy.Signature):
    """Generate an image prompt for a project thumbnail based on extracted project info."""

    title: str = dspy.InputField(desc="Projektname/-titel")
    description: str = dspy.InputField(desc="Projektbeschreibung")
    keywords: list[str] = dspy.InputField(desc="Extrahierte Schlüsselwörter/Themen")
    variation_hint: str = dspy.InputField(
        desc="Hinweis für Variation: welcher Stil oder Fokus soll gewählt werden"
    )

    image_prompt: str = dspy.OutputField(
        desc="A detailed, visual prompt in ENGLISH for generating a square thumbnail image. "
        "Describe concrete visual elements: symbols, objects, colors, composition. "
        "Focus on imagery that represents the project's theme and purpose. "
        "NEVER include text, letters, or words in the image. "
        "NEVER depict people, faces, children, or human figures — use abstract symbols instead. "
        "NEVER reference weapons, violence, conflict, or military themes — even for LARP/adventure projects, "
        "focus on nature, landscapes, fantasy symbols, or abstract representations. "
        "Use the variation_hint to choose a DIFFERENT style or focus than previous attempts. "
        "Style: modern, clean, professional illustration. Prefer abstract, symbolic, geometric designs."
    )
    reasoning: str = dspy.OutputField(
        desc="Kurze Erklärung auf Deutsch, warum diese visuellen Elemente gewählt wurden"
    )


class ImageGenerator:
    """Generate project thumbnail images using Google Imagen."""

    def __init__(self):
        settings = get_settings()
        self.google_api_key = settings.google_api_key
        self.image_model = settings.image_model
        self.image_aspect_ratio = settings.image_aspect_ratio
        self.prompt_generator = dspy.Predict(GenerateImagePrompt)

        # Initialize Google GenAI client
        if self.google_api_key:
            self.client = genai.Client(api_key=self.google_api_key)
        else:
            self.client = None

    # Variation hints for generating different image styles
    VARIATION_HINTS = [
        "Fokus auf Symbole und abstrakte Darstellung",
        "Fokus auf Menschen und Aktivitäten",
        "Fokus auf Atmosphäre und Stimmung mit dramatischer Beleuchtung",
        "Fokus auf Objekte und Requisiten im Detail",
        "Fokus auf Landschaft oder Umgebung",
        "Minimalistischer, ikonischer Stil",
        "Cinematischer, filmischer Stil",
        "Illustrativer, künstlerischer Stil",
    ]

    async def generate(
        self, title: str, description: str, keywords: list[str], variation: int = 0
    ) -> tuple[bytes, str, str]:
        """
        Generate a project thumbnail image.

        Args:
            title: Project title
            description: Project description
            keywords: List of keywords/topics
            variation: Variation number for different styles (0-7)

        Returns:
            Tuple of (image_bytes, prompt_used, reasoning)
        """
        if not self.google_api_key or not self.client:
            raise ValueError("GOOGLE_API_KEY is not configured")

        # Select variation hint based on variation number
        variation_hint = self.VARIATION_HINTS[variation % len(self.VARIATION_HINTS)]

        # Step 1: Use LLM to generate optimized image prompt
        logger.info(f"Generating image prompt for project: {title} (variation: {variation})")
        prompt_result = self.prompt_generator(
            title=title,
            description=description[:1000],  # Limit description length
            keywords=keywords[:10] if keywords else [],  # Limit keywords
            variation_hint=variation_hint,
        )

        image_prompt = prompt_result.image_prompt
        reasoning = prompt_result.reasoning
        logger.info(f"Generated prompt: {image_prompt}")

        # Step 2: Call Imagen API
        image_bytes = await self._call_imagen(image_prompt)

        return image_bytes, image_prompt, reasoning

    async def _call_imagen(self, prompt: str) -> bytes:
        """
        Call Google Imagen API to generate an image.

        Args:
            prompt: The image generation prompt

        Returns:
            Image bytes (PNG format)
        """
        logger.info(f"Calling Imagen with model: {self.image_model}")

        # Run the synchronous API in a thread pool
        def generate_sync():
            response = self.client.models.generate_images(
                model=self.image_model,
                prompt=prompt,
                config=types.GenerateImagesConfig(
                    number_of_images=1,
                    aspect_ratio=self.image_aspect_ratio,
                    output_mime_type="image/png",
                ),
            )
            return response

        response = await asyncio.to_thread(generate_sync)

        if not response.generated_images:
            # Retry with a simplified, safe prompt
            logger.warning(f"Imagen returned no image for prompt: {prompt}")
            safe_prompt = (
                "A colorful abstract geometric illustration with soft gradients, "
                "representing community and nature. Clean, modern, professional style. "
                "No text, no people, no faces."
            )

            def generate_safe():
                return self.client.models.generate_images(
                    model=self.image_model,
                    prompt=safe_prompt,
                    config=types.GenerateImagesConfig(
                        number_of_images=1,
                        aspect_ratio=self.image_aspect_ratio,
                        output_mime_type="image/png",
                    ),
                )

            response = await asyncio.to_thread(generate_safe)
            if not response.generated_images:
                raise ValueError(
                    f"Bildgenerierung fehlgeschlagen (vermutlich Content-Filter). "
                    f"Ursprünglicher Prompt: {prompt}"
                )

        # Get the first image
        generated_image = response.generated_images[0]

        # Convert to bytes
        return self._image_to_bytes(generated_image.image)

    def _image_to_bytes(self, image) -> bytes:
        """Convert Imagen image response to PNG bytes."""
        # google.genai.types.Image has image_bytes attribute
        if hasattr(image, 'image_bytes') and image.image_bytes:
            return image.image_bytes

        # If it has data attribute (base64)
        if hasattr(image, 'data'):
            if isinstance(image.data, bytes):
                return image.data
            return base64.b64decode(image.data)

        # If it's a PIL Image
        if hasattr(image, 'save'):
            buffer = io.BytesIO()
            image.save(buffer, "PNG")
            return buffer.getvalue()

        raise ValueError(f"Unknown image format: {type(image)}, attrs: {dir(image)}")
