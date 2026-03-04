<?php
/**
 * Shortcode Registry Tests
 *
 * @package Waldritter\ProjectDB\Tests
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB\Tests;

use PHPUnit\Framework\TestCase;
use Waldritter\ProjectDB\ShortcodeRegistry;
use Waldritter\ProjectDB\GraphQLClient;
use Waldritter\ProjectDB\CacheManager;

class ShortcodeRegistryTest extends TestCase {
    private ShortcodeRegistry $registry;
    private GraphQLClient $client;
    private CacheManager $cache;

    protected function setUp(): void {
        parent::setUp();
        $this->cache = new CacheManager();
        $this->client = new GraphQLClient('https://api.waldritter.dev/graphql', $this->cache);
        $this->registry = new ShortcodeRegistry($this->client);
    }

    public function testRenderCarouselWithEmptyProjects(): void {
        $atts = [];

        $output = $this->registry->render_carousel($atts);

        $this->assertStringContainsString('waldritter-carousel', $output);
        // Should show placeholder when no projects
        $this->assertStringContainsString('Keine Projekte gefunden', $output);
    }

    public function testRenderCarouselWithAutoScrollSettings(): void {
        $atts = [
            'auto_scroll' => 'true',
            'interval' => '3000',
        ];

        $output = $this->registry->render_carousel($atts);

        $this->assertStringContainsString('data-auto-scroll="true"', $output);
        $this->assertStringContainsString('data-interval="3000"', $output);
    }

    public function testRenderSearchWithDefaultSettings(): void {
        $atts = [];

        $output = $this->registry->render_search($atts);

        $this->assertStringContainsString('waldritter-search', $output);
        $this->assertStringContainsString('data-per-page="5"', $output);
    }

    public function testRenderSearchWithCustomPerPage(): void {
        $atts = [
            'per_page' => '10',
        ];

        $output = $this->registry->render_search($atts);

        $this->assertStringContainsString('data-per-page="10"', $output);
    }

    public function testRenderSearchIncludesFilterButton(): void {
        $atts = [];

        $output = $this->registry->render_search($atts);

        $this->assertStringContainsString('waldritter-search__filter-button', $output);
        $this->assertStringContainsString('Filter', $output);
    }

    public function testRenderSearchIncludesDrawer(): void {
        $atts = [];

        $output = $this->registry->render_search($atts);

        $this->assertStringContainsString('waldritter-search__drawer', $output);
        $this->assertStringContainsString('aria-hidden="true"', $output);
    }
}
