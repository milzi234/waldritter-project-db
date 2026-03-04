<?php
/**
 * Shortcode Registry
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

/**
 * Registers and handles all plugin shortcodes
 */
class ShortcodeRegistry {
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
     * Register all shortcodes
     */
    public function register(): void {
        add_shortcode('waldritter_carousel', [$this, 'render_carousel']);
        add_shortcode('waldritter_search', [$this, 'render_search']);
    }

    /**
     * Render the carousel widget
     *
     * @param array|string $atts Shortcode attributes
     */
    public function render_carousel($atts): string {
        $atts = shortcode_atts([
            'id' => 'carousel-' . wp_unique_id(),
            'tags' => '',
            'project_ids' => '',
            'hidden_categories' => '',
            'auto_scroll' => 'true',
            'interval' => '5000',
            'highlights' => '',
        ], $atts, 'waldritter_carousel');

        // Parse comma-separated values
        $tag_titles = array_filter(array_map('trim', explode(',', $atts['tags'])));
        $project_ids = array_filter(array_map('trim', explode(',', $atts['project_ids'])));
        $hidden_categories = array_filter(array_map('trim', explode(',', $atts['hidden_categories'])));

        // Convert string values
        $atts['auto_scroll'] = filter_var($atts['auto_scroll'], FILTER_VALIDATE_BOOLEAN);
        $atts['interval'] = (int) $atts['interval'];

        // Parse highlights JSON if provided
        $manual_highlights = [];
        if (!empty($atts['highlights'])) {
            $decoded = json_decode($atts['highlights'], true);
            if (is_array($decoded)) {
                $manual_highlights = $decoded;
            }
        }

        // Fetch projects
        $projects = [];
        try {
            if (!empty($project_ids)) {
                // Fetch specific projects by ID
                foreach ($project_ids as $id) {
                    $project = $this->client->get_project($id);
                    if ($project) {
                        $projects[] = $project;
                    }
                }
            } else {
                // Fetch projects by tag filter
                $projects = $this->client->get_projects_by_tags($tag_titles, $hidden_categories);
            }
        } catch (GraphQLException $e) {
            error_log('Waldritter Project DB: Failed to fetch carousel projects: ' . $e->getMessage());
        }

        // Limit to 10 items max
        $projects = array_slice($projects, 0, 10);

        return $this->render_template('carousel', [
            'attributes' => $atts,
            'projects' => $projects,
            'highlights' => $manual_highlights,
        ]);
    }

    /**
     * Render the search widget
     *
     * @param array|string $atts Shortcode attributes
     */
    public function render_search($atts): string {
        $atts = shortcode_atts([
            'id' => 'search-' . wp_unique_id(),
            'per_page' => '5',
            'hidden_categories' => '',
            'initial_tags' => '',
            'base_pattern_tags' => '',
            'show_occurrences' => 'false',
        ], $atts, 'waldritter_search');

        // Parse comma-separated values
        $hidden_categories = array_filter(array_map('trim', explode(',', $atts['hidden_categories'])));
        $initial_tags = array_filter(array_map('trim', explode(',', $atts['initial_tags'])));
        $base_pattern_tags = array_filter(array_map('trim', explode(',', $atts['base_pattern_tags'])));

        // Convert string values
        $atts['per_page'] = max(1, min(20, (int) $atts['per_page']));
        $atts['show_occurrences'] = filter_var($atts['show_occurrences'], FILTER_VALIDATE_BOOLEAN);

        // Merge base pattern with initial tags for server-side render (AND logic)
        $effective_tags = array_unique(array_merge($base_pattern_tags, $initial_tags));

        // Fetch categories for filter UI
        $categories = [];
        $projects = [];
        try {
            $categories = $this->client->get_categories();
            $projects = $this->client->get_projects_by_tags($effective_tags, $hidden_categories);
        } catch (GraphQLException $e) {
            error_log('Waldritter Project DB: Failed to fetch search data: ' . $e->getMessage());
        }

        // Filter out hidden categories from display
        if (!empty($hidden_categories)) {
            $categories = array_filter($categories, function ($category) use ($hidden_categories) {
                return !in_array($category['title'], $hidden_categories, true);
            });
            $categories = array_values($categories);
        }

        return $this->render_template('search', [
            'attributes' => $atts,
            'categories' => $categories,
            'projects' => $projects,
            'hidden_categories' => $hidden_categories,
            'initial_tags' => $initial_tags,
            'base_pattern_tags' => $base_pattern_tags,
        ]);
    }

    /**
     * Render a template file
     *
     * @param string $template Template name
     * @param array $data Data to pass to template
     */
    private function render_template(string $template, array $data): string {
        $template_file = WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'src/templates/' . $template . '.php';

        if (!file_exists($template_file)) {
            return sprintf(
                '<!-- Waldritter Project DB: Template not found: %s -->',
                esc_html($template)
            );
        }

        // Extract data for template
        extract($data, EXTR_SKIP);

        ob_start();
        include $template_file;
        return ob_get_clean();
    }
}
