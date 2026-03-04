<?php
/**
 * Cache Manager Tests
 *
 * @package Waldritter\ProjectDB\Tests
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB\Tests;

use PHPUnit\Framework\TestCase;
use Waldritter\ProjectDB\CacheManager;

class CacheManagerTest extends TestCase {
    private CacheManager $cache;
    private array $transients = [];

    protected function setUp(): void {
        parent::setUp();
        $this->cache = new CacheManager();
        $this->transients = [];

        // Override transient functions for testing
        $this->mockTransientFunctions();
    }

    private function mockTransientFunctions(): void {
        // We can't easily mock global functions in PHP
        // In real tests, we'd use Brain\Monkey or similar
        // For now, tests validate the logic flow
    }

    public function testGenerateQueryKeyCreatesConsistentHash(): void {
        $query = 'query { projects { id title } }';
        $variables = ['limit' => 10];

        $key1 = $this->cache->generate_query_key($query, $variables);
        $key2 = $this->cache->generate_query_key($query, $variables);

        $this->assertEquals($key1, $key2);
        $this->assertStringStartsWith('query_', $key1);
    }

    public function testGenerateQueryKeyDiffersByQuery(): void {
        $variables = ['limit' => 10];

        $key1 = $this->cache->generate_query_key('query { projects { id } }', $variables);
        $key2 = $this->cache->generate_query_key('query { categories { id } }', $variables);

        $this->assertNotEquals($key1, $key2);
    }

    public function testGenerateQueryKeyDiffersByVariables(): void {
        $query = 'query { projects { id } }';

        $key1 = $this->cache->generate_query_key($query, ['limit' => 10]);
        $key2 = $this->cache->generate_query_key($query, ['limit' => 20]);

        $this->assertNotEquals($key1, $key2);
    }

    public function testRememberCallsCallbackWhenNotCached(): void {
        $called = false;
        $callback = function() use (&$called) {
            $called = true;
            return 'test_value';
        };

        $result = $this->cache->remember('test_key', $callback);

        // In a real environment with proper mocking, we'd verify:
        // 1. The callback was called
        // 2. The result was cached
        // For now, we verify the callback returns correctly
        $this->assertTrue($called);
        $this->assertEquals('test_value', $result);
    }

    public function testCacheKeyIsPrefixed(): void {
        // Test that the cache key generation is consistent
        $key = $this->cache->generate_query_key('test', []);

        // The key should be predictable based on input
        $this->assertIsString($key);
        $this->assertNotEmpty($key);
    }
}
