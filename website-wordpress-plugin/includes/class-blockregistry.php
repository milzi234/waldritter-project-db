<?php
/**
 * Gutenberg Block Registry
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

/**
 * Registers and handles Gutenberg blocks
 */
class BlockRegistry {
    /**
     * @var GraphQLClient GraphQL client instance
     */
    private GraphQLClient $client;

    /**
     * Constructor
     *
     * @param GraphQLClient $client GraphQL client
     */
    public function __construct(GraphQLClient $client) {
        $this->client = $client;
    }

    /**
     * Register all Gutenberg blocks
     */
    public function register(): void {
        // Carousel block
        register_block_type('waldritter/carousel', [
            'editor_script' => 'waldritter-project-db-blocks',
            'editor_style' => 'waldritter-project-db-editor',
            'render_callback' => [$this, 'render_carousel'],
            'attributes' => [
                'id' => [
                    'type' => 'string',
                    'default' => '',
                ],
                'tags' => [
                    'type' => 'array',
                    'default' => [],
                    'items' => ['type' => 'string'],
                ],
                'projectIds' => [
                    'type' => 'array',
                    'default' => [],
                    'items' => ['type' => 'string'],
                ],
                'hiddenCategories' => [
                    'type' => 'array',
                    'default' => [],
                    'items' => ['type' => 'string'],
                ],
                'autoScroll' => [
                    'type' => 'boolean',
                    'default' => true,
                ],
                'interval' => [
                    'type' => 'number',
                    'default' => 5000,
                ],
                'highlights' => [
                    'type' => 'array',
                    'default' => [],
                ],
            ],
        ]);

        // Search block
        register_block_type('waldritter/search', [
            'editor_script' => 'waldritter-project-db-blocks',
            'editor_style' => 'waldritter-project-db-editor',
            'render_callback' => [$this, 'render_search'],
            'attributes' => [
                'id' => [
                    'type' => 'string',
                    'default' => '',
                ],
                'perPage' => [
                    'type' => 'number',
                    'default' => 5,
                ],
                'hiddenCategories' => [
                    'type' => 'array',
                    'default' => [],
                    'items' => ['type' => 'string'],
                ],
                'initialTags' => [
                    'type' => 'array',
                    'default' => [],
                    'items' => ['type' => 'string'],
                ],
                'basePatternTags' => [
                    'type' => 'array',
                    'default' => [],
                    'items' => ['type' => 'string'],
                ],
                'showOccurrences' => [
                    'type' => 'boolean',
                    'default' => false,
                ],
            ],
        ]);
    }

    /**
     * Render carousel block
     *
     * @param array $attributes Block attributes
     */
    public function render_carousel(array $attributes): string {
        $shortcode_registry = new ShortcodeRegistry($this->client);

        // Convert arrays to comma-separated strings
        $atts = [
            'id' => $attributes['id'] ?: 'carousel-' . wp_unique_id(),
            'tags' => implode(',', $attributes['tags'] ?? []),
            'project_ids' => implode(',', $attributes['projectIds'] ?? []),
            'hidden_categories' => implode(',', $attributes['hiddenCategories'] ?? []),
            'auto_scroll' => ($attributes['autoScroll'] ?? true) ? 'true' : 'false',
            'interval' => (string) ($attributes['interval'] ?? 5000),
            'highlights' => !empty($attributes['highlights']) ? wp_json_encode($attributes['highlights']) : '',
        ];

        return $shortcode_registry->render_carousel($atts);
    }

    /**
     * Render search block
     *
     * @param array $attributes Block attributes
     */
    public function render_search(array $attributes): string {
        $shortcode_registry = new ShortcodeRegistry($this->client);

        // Convert arrays to comma-separated strings
        $atts = [
            'id' => $attributes['id'] ?: 'search-' . wp_unique_id(),
            'per_page' => (string) ($attributes['perPage'] ?? 5),
            'hidden_categories' => implode(',', $attributes['hiddenCategories'] ?? []),
            'initial_tags' => implode(',', $attributes['initialTags'] ?? []),
            'base_pattern_tags' => implode(',', $attributes['basePatternTags'] ?? []),
            'show_occurrences' => ($attributes['showOccurrences'] ?? false) ? 'true' : 'false',
        ];

        return $shortcode_registry->render_search($atts);
    }
}
