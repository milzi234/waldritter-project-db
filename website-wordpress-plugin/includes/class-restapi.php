<?php
/**
 * REST API endpoints
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

use WP_REST_Request;
use WP_REST_Response;
use WP_Error;

/**
 * Registers REST API endpoints for Gutenberg block previews and AJAX operations
 */
class RestAPI {
    /**
     * @var string REST API namespace
     */
    private const NAMESPACE = 'waldritter/v1';

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
     * Register REST API routes
     */
    public function register_routes(): void {
        // Get all projects
        register_rest_route(self::NAMESPACE, '/projects', [
            'methods' => 'GET',
            'callback' => [$this, 'get_projects'],
            'permission_callback' => '__return_true',
            'args' => [
                'tags' => [
                    'type' => 'array',
                    'items' => ['type' => 'string'],
                    'default' => [],
                ],
                'hidden_categories' => [
                    'type' => 'array',
                    'items' => ['type' => 'string'],
                    'default' => [],
                ],
                'limit' => [
                    'type' => 'integer',
                    'default' => 20,
                    'minimum' => 1,
                    'maximum' => 100,
                ],
                'offset' => [
                    'type' => 'integer',
                    'default' => 0,
                    'minimum' => 0,
                ],
            ],
        ]);

        // Get single project
        register_rest_route(self::NAMESPACE, '/projects/(?P<id>\d+)', [
            'methods' => 'GET',
            'callback' => [$this, 'get_project'],
            'permission_callback' => '__return_true',
            'args' => [
                'id' => [
                    'type' => 'string',
                    'required' => true,
                ],
            ],
        ]);

        // Get all categories
        register_rest_route(self::NAMESPACE, '/categories', [
            'methods' => 'GET',
            'callback' => [$this, 'get_categories'],
            'permission_callback' => '__return_true',
        ]);

        // Get all tags
        register_rest_route(self::NAMESPACE, '/tags', [
            'methods' => 'GET',
            'callback' => [$this, 'get_tags'],
            'permission_callback' => '__return_true',
        ]);

        // Search projects
        register_rest_route(self::NAMESPACE, '/search', [
            'methods' => 'GET',
            'callback' => [$this, 'search'],
            'permission_callback' => '__return_true',
            'args' => [
                'tags' => [
                    'type' => 'array',
                    'items' => ['type' => 'string'],
                    'default' => [],
                ],
                'hidden_categories' => [
                    'type' => 'array',
                    'items' => ['type' => 'string'],
                    'default' => [],
                ],
                'page' => [
                    'type' => 'integer',
                    'default' => 1,
                    'minimum' => 1,
                ],
                'per_page' => [
                    'type' => 'integer',
                    'default' => 5,
                    'minimum' => 1,
                    'maximum' => 20,
                ],
            ],
        ]);

        // Clear cache (admin only)
        register_rest_route(self::NAMESPACE, '/cache/clear', [
            'methods' => 'POST',
            'callback' => [$this, 'clear_cache'],
            'permission_callback' => [$this, 'check_admin_permission'],
        ]);
    }

    /**
     * Get projects
     *
     * @param WP_REST_Request $request Request object
     */
    public function get_projects(WP_REST_Request $request): WP_REST_Response|WP_Error {
        try {
            $tags = $request->get_param('tags') ?? [];
            $hidden_categories = $request->get_param('hidden_categories') ?? [];

            $projects = $this->client->get_projects_by_tags($tags, $hidden_categories);

            // Apply pagination
            $limit = (int) $request->get_param('limit');
            $offset = (int) $request->get_param('offset');
            $total = count($projects);
            $projects = array_slice($projects, $offset, $limit);

            return new WP_REST_Response([
                'projects' => $projects,
                'total' => $total,
                'limit' => $limit,
                'offset' => $offset,
            ], 200);
        } catch (GraphQLException $e) {
            return new WP_Error(
                'api_error',
                $e->get_user_message(),
                ['status' => 500]
            );
        }
    }

    /**
     * Get a single project
     *
     * @param WP_REST_Request $request Request object
     */
    public function get_project(WP_REST_Request $request): WP_REST_Response|WP_Error {
        try {
            $id = $request->get_param('id');
            $project = $this->client->get_project($id);

            if ($project === null) {
                return new WP_Error(
                    'not_found',
                    __('Project not found', 'waldritter-project-db'),
                    ['status' => 404]
                );
            }

            return new WP_REST_Response($project, 200);
        } catch (GraphQLException $e) {
            return new WP_Error(
                'api_error',
                $e->get_user_message(),
                ['status' => 500]
            );
        }
    }

    /**
     * Get all categories
     */
    public function get_categories(): WP_REST_Response|WP_Error {
        try {
            $categories = $this->client->get_categories();
            return new WP_REST_Response($categories, 200);
        } catch (GraphQLException $e) {
            return new WP_Error(
                'api_error',
                $e->get_user_message(),
                ['status' => 500]
            );
        }
    }

    /**
     * Get all tags
     */
    public function get_tags(): WP_REST_Response|WP_Error {
        try {
            $tags = $this->client->get_tags();
            return new WP_REST_Response($tags, 200);
        } catch (GraphQLException $e) {
            return new WP_Error(
                'api_error',
                $e->get_user_message(),
                ['status' => 500]
            );
        }
    }

    /**
     * Search projects with pagination
     *
     * @param WP_REST_Request $request Request object
     */
    public function search(WP_REST_Request $request): WP_REST_Response|WP_Error {
        try {
            $tags = $request->get_param('tags') ?? [];
            $hidden_categories = $request->get_param('hidden_categories') ?? [];
            $page = (int) $request->get_param('page');
            $per_page = (int) $request->get_param('per_page');

            $projects = $this->client->get_projects_by_tags($tags, $hidden_categories);

            // Calculate pagination
            $total = count($projects);
            $total_pages = (int) ceil($total / $per_page);
            $offset = ($page - 1) * $per_page;
            $projects = array_slice($projects, $offset, $per_page);

            return new WP_REST_Response([
                'projects' => $projects,
                'pagination' => [
                    'page' => $page,
                    'per_page' => $per_page,
                    'total' => $total,
                    'total_pages' => $total_pages,
                ],
            ], 200);
        } catch (GraphQLException $e) {
            return new WP_Error(
                'api_error',
                $e->get_user_message(),
                ['status' => 500]
            );
        }
    }

    /**
     * Clear the cache
     */
    public function clear_cache(): WP_REST_Response {
        CacheManager::flush_all();

        return new WP_REST_Response([
            'success' => true,
            'message' => __('Cache cleared successfully', 'waldritter-project-db'),
        ], 200);
    }

    /**
     * Check if current user has admin permission
     */
    public function check_admin_permission(): bool {
        return current_user_can('manage_options');
    }
}
