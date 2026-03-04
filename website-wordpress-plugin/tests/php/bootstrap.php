<?php
/**
 * PHPUnit bootstrap file
 *
 * @package Waldritter\ProjectDB\Tests
 */

declare(strict_types=1);

// Composer autoload
require_once dirname(__DIR__, 2) . '/vendor/autoload.php';

// Define WordPress constants for testing
if (!defined('ABSPATH')) {
    define('ABSPATH', '/var/www/html/');
}

// Yoast PHPUnit Polyfills
use Yoast\PHPUnitPolyfills\Autoload;

// Brain Monkey for WordPress function mocking
use Brain\Monkey;

// Set up mocking
Monkey\setUp();

// Mock WordPress functions commonly used
function add_action(...$args) {}
function add_filter(...$args) {}
function add_shortcode(...$args) {}
function do_action(...$args) {}
function apply_filters($tag, $value, ...$args) { return $value; }
function wp_enqueue_style(...$args) {}
function wp_enqueue_script(...$args) {}
function wp_localize_script(...$args) {}
function wp_remote_post(...$args) { return []; }
function wp_remote_retrieve_response_code($response) { return 200; }
function wp_remote_retrieve_body($response) { return '{}'; }
function is_wp_error($thing) { return false; }
function wp_json_encode($data) { return json_encode($data); }
function esc_html($text) { return htmlspecialchars($text, ENT_QUOTES, 'UTF-8'); }
function esc_attr($text) { return htmlspecialchars($text, ENT_QUOTES, 'UTF-8'); }
function esc_url($url) { return filter_var($url, FILTER_SANITIZE_URL); }
function wp_kses_post($text) { return $text; }
function wpautop($text) { return "<p>$text</p>"; }
function wp_unique_id($prefix = '') { return $prefix . uniqid(); }
function __($text, $domain = 'default') { return $text; }
function _e($text, $domain = 'default') { echo $text; }
function _n($single, $plural, $number, $domain = 'default') { return $number === 1 ? $single : $plural; }
function esc_html__($text, $domain = 'default') { return esc_html($text); }
function esc_attr__($text, $domain = 'default') { return esc_attr($text); }
function esc_html_e($text, $domain = 'default') { echo esc_html($text); }
function esc_attr_e($text, $domain = 'default') { echo esc_attr($text); }

function get_option($option, $default = false) {
    static $options = [];
    return $options[$option] ?? $default;
}

function update_option($option, $value) {
    return true;
}

function add_option($option, $value) {
    return true;
}

function get_transient($transient) {
    return false;
}

function set_transient($transient, $value, $expiration = 0) {
    return true;
}

function delete_transient($transient) {
    return true;
}

function plugin_dir_path($file) {
    return dirname($file) . '/';
}

function plugin_dir_url($file) {
    return 'http://example.com/wp-content/plugins/waldritter-project-db/';
}

function load_plugin_textdomain($domain, $deprecated = false, $plugin_rel_path = false) {
    return true;
}

function register_activation_hook($file, $callback) {}
function register_deactivation_hook($file, $callback) {}

function has_shortcode($content, $tag) { return false; }
function has_block($block_name, $post = null) { return false; }

function shortcode_atts($pairs, $atts, $shortcode = '') {
    $atts = (array) $atts;
    $out = [];
    foreach ($pairs as $name => $default) {
        if (array_key_exists($name, $atts)) {
            $out[$name] = $atts[$name];
        } else {
            $out[$name] = $default;
        }
    }
    return $out;
}

function wp_trim_words($text, $num_words = 55, $more = null) {
    $words = explode(' ', $text);
    if (count($words) <= $num_words) {
        return $text;
    }
    return implode(' ', array_slice($words, 0, $num_words)) . ($more ?? '...');
}

function checked($checked, $current = true, $display = true) {
    if ($checked === $current) {
        $result = " checked='checked'";
    } else {
        $result = '';
    }
    if ($display) {
        echo $result;
    }
    return $result;
}

function disabled($disabled, $current = true, $display = true) {
    if ($disabled === $current) {
        $result = " disabled='disabled'";
    } else {
        $result = '';
    }
    if ($display) {
        echo $result;
    }
    return $result;
}

// Define plugin constants
define('WALDRITTER_PROJECT_DB_VERSION', '1.0.0');
define('WALDRITTER_PROJECT_DB_PLUGIN_DIR', dirname(__DIR__, 2) . '/');
define('WALDRITTER_PROJECT_DB_PLUGIN_URL', 'http://example.com/wp-content/plugins/waldritter-project-db/');
define('WALDRITTER_PROJECT_DB_API_URL', 'https://api.waldritter.dev/graphql');

// Require plugin classes
require_once WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'includes/class-cachemanager.php';
require_once WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'includes/class-graphqlexception.php';
require_once WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'includes/class-graphqlclient.php';
require_once WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'includes/class-shortcoderegistry.php';
require_once WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'includes/class-blockregistry.php';
require_once WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'includes/class-restapi.php';
require_once WALDRITTER_PROJECT_DB_PLUGIN_DIR . 'includes/class-plugin.php';
