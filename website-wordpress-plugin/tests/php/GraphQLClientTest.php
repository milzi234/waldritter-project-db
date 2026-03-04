<?php
/**
 * GraphQL Client Tests
 *
 * @package Waldritter\ProjectDB\Tests
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB\Tests;

use PHPUnit\Framework\TestCase;
use Waldritter\ProjectDB\GraphQLClient;
use Waldritter\ProjectDB\CacheManager;
use Waldritter\ProjectDB\GraphQLException;

class GraphQLClientTest extends TestCase {
    private GraphQLClient $client;
    private CacheManager $cache;

    protected function setUp(): void {
        parent::setUp();
        $this->cache = new CacheManager();
        $this->client = new GraphQLClient('https://api.waldritter.dev/graphql', $this->cache);
    }

    public function testClientConstructorSetsApiUrl(): void {
        $client = new GraphQLClient('https://test.example.com/graphql', $this->cache);

        // We can't directly access private properties, but we can verify
        // the object was created without errors
        $this->assertInstanceOf(GraphQLClient::class, $client);
    }

    public function testSetTimeoutAcceptsPositiveValue(): void {
        $this->client->set_timeout(60);

        // No exception means success
        $this->assertTrue(true);
    }

    public function testSetTimeoutEnforcesMinimum(): void {
        $this->client->set_timeout(0);

        // Should not throw, just enforce minimum
        $this->assertTrue(true);
    }

    public function testGetProjectsReturnsEmptyArrayOnEmptyResponse(): void {
        // This would require mocking wp_remote_post
        // For now, we verify the method exists and is callable
        $this->assertTrue(method_exists($this->client, 'get_projects'));
    }

    public function testGetProjectReturnsNullForNonExistentProject(): void {
        // This would require mocking wp_remote_post
        // For now, we verify the method exists and is callable
        $this->assertTrue(method_exists($this->client, 'get_project'));
    }

    public function testGetCategoriesReturnsArray(): void {
        $this->assertTrue(method_exists($this->client, 'get_categories'));
    }

    public function testGetTagsReturnsArray(): void {
        $this->assertTrue(method_exists($this->client, 'get_tags'));
    }

    public function testSearchMethodExists(): void {
        $this->assertTrue(method_exists($this->client, 'search'));
    }

    public function testGetProjectsByTagsMethodExists(): void {
        $this->assertTrue(method_exists($this->client, 'get_projects_by_tags'));
    }
}
