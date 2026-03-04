<?php
/**
 * Plugin Name: Waldritter Project DB Integration
 * Plugin URI: https://waldritter.de
 * Description: Integrates with the Waldritter Project Database API to display projects via shortcodes and Gutenberg blocks.
 * Version: 1.0.0
 * Author: Waldritter e.V.
 * Author URI: https://waldritter.de
 * Text Domain: waldritter-project-db
 * Domain Path: /languages
 * Requires at least: 6.0
 * Requires PHP: 8.0
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

// Prevent direct access
if (!defined('ABSPATH')) {
    exit;
}

// Load Composer autoloader for dependencies (Parsedown, etc.)
$composer_autoload = __DIR__ . '/vendor/autoload.php';
if (file_exists($composer_autoload)) {
    require_once $composer_autoload;
}

// Plugin constants
define('WALDRITTER_PROJECT_DB_VERSION', '1.0.0');
define('WALDRITTER_PROJECT_DB_PLUGIN_DIR', plugin_dir_path(__FILE__));
define('WALDRITTER_PROJECT_DB_PLUGIN_URL', plugin_dir_url(__FILE__));
define('WALDRITTER_PROJECT_DB_API_URL', 'https://project-api.waldritter.dev/graphql');

// Autoload classes
spl_autoload_register(function (string $class): void {
    $prefix = 'Waldritter\\ProjectDB\\';
    $base_dir = WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'includes/';

    $len = strlen($prefix);
    if (strncmp($prefix, $class, $len) !== 0) {
        return;
    }

    $relative_class = substr($class, $len);
    $file = $base_dir . 'class-' . strtolower(str_replace(['\\', '_'], ['-', '-'], $relative_class)) . '.php';

    if (file_exists($file)) {
        require $file;
    }
});

// Initialize the plugin
function init(): void {
    // Load text domain for translations
    load_plugin_textdomain(
        'waldritter-project-db',
        false,
        dirname(plugin_basename(__FILE__)) . '/languages'
    );

    // Initialize main plugin class
    Plugin::get_instance()->init();
}
add_action('plugins_loaded', __NAMESPACE__ . '\\init');

// Activation hook
register_activation_hook(__FILE__, function (): void {
    // Clear any cached data
    CacheManager::flush_all();

    // Set default options
    if (get_option('waldritter_api_url') === false) {
        add_option('waldritter_api_url', WALDRITTER_PROJECT_DB_API_URL);
    }
    if (get_option('waldritter_cache_ttl') === false) {
        add_option('waldritter_cache_ttl', 3600); // 1 hour default
    }
});

// Deactivation hook
register_deactivation_hook(__FILE__, function (): void {
    CacheManager::flush_all();
});
