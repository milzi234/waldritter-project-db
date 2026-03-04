<?php
/**
 * Markdown Renderer
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

use Parsedown;

/**
 * Helper class for rendering markdown content
 */
class MarkdownRenderer {
    /**
     * @var Parsedown Parsedown instance
     */
    private static ?Parsedown $parsedown = null;

    /**
     * Get the Parsedown instance
     */
    private static function getParsedown(): Parsedown {
        if (self::$parsedown === null) {
            self::$parsedown = new Parsedown();
            self::$parsedown->setSafeMode(true);
        }
        return self::$parsedown;
    }

    /**
     * Render markdown to HTML
     *
     * @param string $markdown The markdown content
     * @return string Sanitized HTML output
     */
    public static function render(string $markdown): string {
        if (empty($markdown)) {
            return '';
        }

        $html = self::getParsedown()->text($markdown);

        // Sanitize output using WordPress's wp_kses_post
        return wp_kses_post($html);
    }

    /**
     * Render markdown to HTML and trim to a word limit
     *
     * Useful for excerpts/previews where you need truncated plain text.
     *
     * @param string $markdown The markdown content
     * @param int $word_limit Maximum number of words
     * @return string Plain text, trimmed
     */
    public static function renderExcerpt(string $markdown, int $word_limit = 20): string {
        if (empty($markdown)) {
            return '';
        }

        // First render markdown to HTML
        $html = self::getParsedown()->text($markdown);

        // Strip HTML tags to get plain text
        $text = wp_strip_all_tags($html);

        // Trim to word limit
        return wp_trim_words($text, $word_limit, '...');
    }
}
