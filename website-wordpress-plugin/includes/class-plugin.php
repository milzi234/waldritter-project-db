<?php
/**
 * Main plugin class
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

/**
 * Main singleton plugin class that coordinates all components
 */
class Plugin {
    /**
     * @var Plugin|null Singleton instance
     */
    private static ?Plugin $instance = null;

    /**
     * @var GraphQLClient GraphQL API client
     */
    private GraphQLClient $graphql_client;

    /**
     * @var CacheManager Cache manager instance
     */
    private CacheManager $cache_manager;

    /**
     * @var ShortcodeRegistry Shortcode registry
     */
    private ShortcodeRegistry $shortcode_registry;

    /**
     * @var BlockRegistry Gutenberg block registry
     */
    private BlockRegistry $block_registry;

    /**
     * @var RestAPI REST API endpoints
     */
    private RestAPI $rest_api;

    /**
     * Private constructor for singleton pattern
     */
    private function __construct() {
        $this->cache_manager = new CacheManager();
        $this->graphql_client = new GraphQLClient(
            $this->get_api_url(),
            $this->cache_manager
        );
        $this->shortcode_registry = new ShortcodeRegistry($this->graphql_client);
        $this->block_registry = new BlockRegistry($this->graphql_client);
        $this->rest_api = new RestAPI($this->graphql_client);
    }

    /**
     * Get singleton instance
     */
    public static function get_instance(): Plugin {
        if (self::$instance === null) {
            self::$instance = new Plugin();
        }
        return self::$instance;
    }

    /**
     * Initialize the plugin components
     */
    public function init(): void {
        // Register shortcodes
        $this->shortcode_registry->register();

        // Register Gutenberg blocks
        add_action('init', [$this->block_registry, 'register']);

        // Register REST API endpoints
        add_action('rest_api_init', [$this->rest_api, 'register_routes']);

        // Enqueue frontend assets
        add_action('wp_enqueue_scripts', [$this, 'enqueue_frontend_assets']);

        // Enqueue editor assets
        add_action('enqueue_block_editor_assets', [$this, 'enqueue_editor_assets']);

        // Register admin settings
        add_action('admin_menu', [$this, 'register_admin_menu']);
        add_action('admin_init', [$this, 'register_settings']);
    }

    /**
     * Enqueue frontend CSS and JS
     */
    public function enqueue_frontend_assets(): void {
        // Only enqueue if a widget is on the page
        global $post;
        if (!$post) {
            return;
        }

        $has_widget = has_shortcode($post->post_content, 'waldritter_carousel')
            || has_shortcode($post->post_content, 'waldritter_search')
            || has_block('waldritter/carousel', $post)
            || has_block('waldritter/search', $post);

        if (!$has_widget) {
            return;
        }

        wp_enqueue_style(
            'waldritter-project-db',
            WALDRITTER_PROJECT_DB_PLUGIN_URL . 'assets/css/frontend.css',
            [],
            WALDRITTER_PROJECT_DB_VERSION
        );

        // Enqueue marked.js for markdown rendering
        wp_enqueue_script(
            'marked',
            'https://cdn.jsdelivr.net/npm/marked@9.1.6/marked.min.js',
            [],
            '9.1.6',
            true
        );

        wp_enqueue_script(
            'waldritter-project-db',
            WALDRITTER_PROJECT_DB_PLUGIN_URL . 'assets/js/frontend.js',
            ['marked'],
            WALDRITTER_PROJECT_DB_VERSION,
            true
        );

        // Localize script with REST API URL
        wp_localize_script('waldritter-project-db', 'waldritterProjectDB', [
            'restUrl' => rest_url('waldritter/v1/'),
            'nonce' => wp_create_nonce('wp_rest'),
        ]);
    }

    /**
     * Enqueue block editor assets
     */
    public function enqueue_editor_assets(): void {
        wp_enqueue_script(
            'waldritter-project-db-blocks',
            WALDRITTER_PROJECT_DB_PLUGIN_URL . 'assets/js/blocks.js',
            ['wp-blocks', 'wp-element', 'wp-editor', 'wp-components', 'wp-i18n', 'wp-api-fetch'],
            WALDRITTER_PROJECT_DB_VERSION,
            true
        );

        wp_enqueue_style(
            'waldritter-project-db-editor',
            WALDRITTER_PROJECT_DB_PLUGIN_URL . 'assets/css/editor.css',
            ['wp-edit-blocks'],
            WALDRITTER_PROJECT_DB_VERSION
        );

        wp_localize_script('waldritter-project-db-blocks', 'waldritterProjectDB', [
            'restUrl' => rest_url('waldritter/v1/'),
            'nonce' => wp_create_nonce('wp_rest'),
        ]);
    }

    /**
     * Register admin menu
     */
    public function register_admin_menu(): void {
        add_options_page(
            __('Waldritter Project DB', 'waldritter-project-db'),
            __('Waldritter DB', 'waldritter-project-db'),
            'manage_options',
            'waldritter-project-db',
            [$this, 'render_settings_page']
        );
    }

    /**
     * Register plugin settings
     */
    public function register_settings(): void {
        register_setting('waldritter_project_db_settings', 'waldritter_api_url', [
            'type' => 'string',
            'sanitize_callback' => 'esc_url_raw',
            'default' => WALDRITTER_PROJECT_DB_API_URL,
        ]);

        register_setting('waldritter_project_db_settings', 'waldritter_cache_ttl', [
            'type' => 'integer',
            'sanitize_callback' => 'absint',
            'default' => 3600,
        ]);

        add_settings_section(
            'waldritter_project_db_main',
            __('API Settings', 'waldritter-project-db'),
            function (): void {
                echo '<p>' . esc_html__('Configure the connection to the Waldritter Project Database API.', 'waldritter-project-db') . '</p>';
            },
            'waldritter-project-db'
        );

        add_settings_field(
            'waldritter_api_url',
            __('API URL', 'waldritter-project-db'),
            function (): void {
                $value = get_option('waldritter_api_url', WALDRITTER_PROJECT_DB_API_URL);
                echo '<input type="url" name="waldritter_api_url" value="' . esc_attr($value) . '" class="regular-text" />';
                echo '<p class="description">' . esc_html__('The GraphQL endpoint URL.', 'waldritter-project-db') . '</p>';
            },
            'waldritter-project-db',
            'waldritter_project_db_main'
        );

        add_settings_field(
            'waldritter_cache_ttl',
            __('Cache Duration (seconds)', 'waldritter-project-db'),
            function (): void {
                $value = get_option('waldritter_cache_ttl', 3600);
                echo '<input type="number" name="waldritter_cache_ttl" value="' . esc_attr($value) . '" min="0" max="86400" />';
                echo '<p class="description">' . esc_html__('How long to cache API responses. Set to 0 to disable caching.', 'waldritter-project-db') . '</p>';
            },
            'waldritter-project-db',
            'waldritter_project_db_main'
        );
    }

    /**
     * Render admin settings page
     */
    public function render_settings_page(): void {
        if (!current_user_can('manage_options')) {
            return;
        }

        // Check for cache flush action
        if (isset($_POST['waldritter_flush_cache']) && check_admin_referer('waldritter_flush_cache')) {
            $this->cache_manager->flush_all();
            echo '<div class="notice notice-success"><p>' . esc_html__('Cache cleared successfully.', 'waldritter-project-db') . '</p></div>';
        }

        ?>
        <div class="wrap">
            <h1><?php echo esc_html(get_admin_page_title()); ?></h1>

            <form action="options.php" method="post">
                <?php
                settings_fields('waldritter_project_db_settings');
                do_settings_sections('waldritter-project-db');
                submit_button();
                ?>
            </form>

            <hr />

            <h2><?php esc_html_e('Cache Management', 'waldritter-project-db'); ?></h2>
            <form method="post">
                <?php wp_nonce_field('waldritter_flush_cache'); ?>
                <p>
                    <button type="submit" name="waldritter_flush_cache" class="button button-secondary">
                        <?php esc_html_e('Clear Cache', 'waldritter-project-db'); ?>
                    </button>
                </p>
            </form>

            <hr />

            <h2><?php esc_html_e('Available Shortcodes', 'waldritter-project-db'); ?></h2>
            <table class="widefat">
                <thead>
                    <tr>
                        <th><?php esc_html_e('Shortcode', 'waldritter-project-db'); ?></th>
                        <th><?php esc_html_e('Description', 'waldritter-project-db'); ?></th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td><code>[waldritter_carousel]</code></td>
                        <td><?php esc_html_e('Displays a carousel of featured projects with optional tag filtering.', 'waldritter-project-db'); ?></td>
                    </tr>
                    <tr>
                        <td><code>[waldritter_search]</code></td>
                        <td><?php esc_html_e('Displays an interactive project search with tag filtering and pagination.', 'waldritter-project-db'); ?></td>
                    </tr>
                </tbody>
            </table>
        </div>
        <?php
    }

    /**
     * Get the API URL from settings
     */
    public function get_api_url(): string {
        return get_option('waldritter_api_url', WALDRITTER_PROJECT_DB_API_URL);
    }

    /**
     * Get the cache TTL from settings
     */
    public function get_cache_ttl(): int {
        return (int) get_option('waldritter_cache_ttl', 3600);
    }

    /**
     * Get the GraphQL client instance
     */
    public function get_graphql_client(): GraphQLClient {
        return $this->graphql_client;
    }

    /**
     * Get the cache manager instance
     */
    public function get_cache_manager(): CacheManager {
        return $this->cache_manager;
    }
}
