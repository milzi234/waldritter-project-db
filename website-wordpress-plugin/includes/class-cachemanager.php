<?php
/**
 * Cache manager using WordPress transients
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

/**
 * Manages caching of API responses using WordPress transients
 */
class CacheManager {
    /**
     * @var string Cache key prefix
     */
    private const CACHE_PREFIX = 'waldritter_pdb_';

    /**
     * @var int Default TTL in seconds (1 hour)
     */
    private const DEFAULT_TTL = 3600;

    /**
     * Get a cached value
     *
     * @param string $key The cache key
     * @return mixed|null The cached value or null if not found
     */
    public function get(string $key): mixed {
        $cache_key = $this->prefix_key($key);
        $value = get_transient($cache_key);

        if ($value === false) {
            return null;
        }

        return $value;
    }

    /**
     * Set a cached value
     *
     * @param string $key The cache key
     * @param mixed $value The value to cache
     * @param int|null $ttl Time to live in seconds (null uses default)
     */
    public function set(string $key, mixed $value, ?int $ttl = null): bool {
        $cache_key = $this->prefix_key($key);
        $ttl = $ttl ?? $this->get_ttl();

        return set_transient($cache_key, $value, $ttl);
    }

    /**
     * Delete a cached value
     *
     * @param string $key The cache key
     */
    public function delete(string $key): bool {
        $cache_key = $this->prefix_key($key);
        return delete_transient($cache_key);
    }

    /**
     * Check if a key exists in cache
     *
     * @param string $key The cache key
     */
    public function has(string $key): bool {
        return $this->get($key) !== null;
    }

    /**
     * Get or set a cached value
     *
     * @param string $key The cache key
     * @param callable $callback Function to generate value if not cached
     * @param int|null $ttl Time to live in seconds
     * @return mixed The cached or generated value
     */
    public function remember(string $key, callable $callback, ?int $ttl = null): mixed {
        $value = $this->get($key);

        if ($value !== null) {
            return $value;
        }

        $value = $callback();

        if ($value !== null) {
            $this->set($key, $value, $ttl);
        }

        return $value;
    }

    /**
     * Flush all plugin cache entries
     */
    public static function flush_all(): void {
        global $wpdb;

        // Delete all transients with our prefix
        $prefix = '_transient_' . self::CACHE_PREFIX;
        $timeout_prefix = '_transient_timeout_' . self::CACHE_PREFIX;

        $wpdb->query(
            $wpdb->prepare(
                "DELETE FROM {$wpdb->options} WHERE option_name LIKE %s OR option_name LIKE %s",
                $prefix . '%',
                $timeout_prefix . '%'
            )
        );

        // Clear object cache if available
        if (function_exists('wp_cache_flush_group')) {
            wp_cache_flush_group('waldritter_project_db');
        }
    }

    /**
     * Generate a cache key from query and variables
     *
     * @param string $query The GraphQL query
     * @param array $variables Query variables
     */
    public function generate_query_key(string $query, array $variables = []): string {
        $data = $query . '|' . json_encode($variables);
        return 'query_' . md5($data);
    }

    /**
     * Prefix a cache key
     */
    private function prefix_key(string $key): string {
        return self::CACHE_PREFIX . $key;
    }

    /**
     * Get the configured TTL
     */
    private function get_ttl(): int {
        $ttl = get_option('waldritter_cache_ttl', self::DEFAULT_TTL);
        return max(0, (int) $ttl);
    }
}
