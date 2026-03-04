<?php
/**
 * GraphQL API client
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

/**
 * GraphQL client for communicating with the Waldritter Project DB API
 */
class GraphQLClient {
    /**
     * @var string API endpoint URL
     */
    private string $api_url;

    /**
     * @var CacheManager Cache manager instance
     */
    private CacheManager $cache;

    /**
     * @var int Request timeout in seconds
     */
    private int $timeout = 30;

    /**
     * Constructor
     *
     * @param string $api_url The GraphQL API endpoint URL
     * @param CacheManager $cache Cache manager instance
     */
    public function __construct(string $api_url, CacheManager $cache) {
        $this->api_url = $api_url;
        $this->cache = $cache;
    }

    /**
     * Execute a GraphQL query
     *
     * @param string $query The GraphQL query
     * @param array $variables Query variables
     * @param bool $use_cache Whether to use caching
     * @return array The query result
     * @throws GraphQLException On API errors
     */
    public function query(string $query, array $variables = [], bool $use_cache = true): array {
        $cache_key = $this->cache->generate_query_key($query, $variables);

        if ($use_cache) {
            $cached = $this->cache->get($cache_key);
            if ($cached !== null) {
                return $cached;
            }
        }

        $result = $this->execute($query, $variables);

        if ($use_cache && !empty($result)) {
            $this->cache->set($cache_key, $result);
        }

        return $result;
    }

    /**
     * Execute a GraphQL mutation (never cached)
     *
     * @param string $mutation The GraphQL mutation
     * @param array $variables Mutation variables
     * @return array The mutation result
     * @throws GraphQLException On API errors
     */
    public function mutate(string $mutation, array $variables = []): array {
        return $this->execute($mutation, $variables);
    }

    /**
     * Execute a GraphQL operation
     *
     * @param string $query The GraphQL query or mutation
     * @param array $variables Variables
     * @return array The result data
     * @throws GraphQLException On API errors
     */
    private function execute(string $query, array $variables = []): array {
        $body = [
            'query' => $query,
        ];

        if (!empty($variables)) {
            $body['variables'] = $variables;
        }

        $response = wp_remote_post($this->api_url, [
            'timeout' => $this->timeout,
            'headers' => [
                'Content-Type' => 'application/json',
                'Accept' => 'application/json',
            ],
            'body' => wp_json_encode($body),
        ]);

        if (is_wp_error($response)) {
            throw new GraphQLException(
                sprintf(
                    /* translators: %s: error message */
                    __('API request failed: %s', 'waldritter-project-db'),
                    $response->get_error_message()
                ),
                0,
                null,
                $response->get_error_code()
            );
        }

        $status_code = wp_remote_retrieve_response_code($response);
        $body = wp_remote_retrieve_body($response);

        if ($status_code < 200 || $status_code >= 300) {
            throw new GraphQLException(
                sprintf(
                    /* translators: %d: HTTP status code */
                    __('API returned HTTP status %d', 'waldritter-project-db'),
                    $status_code
                ),
                $status_code
            );
        }

        $data = json_decode($body, true);

        if (json_last_error() !== JSON_ERROR_NONE) {
            throw new GraphQLException(
                __('Invalid JSON response from API', 'waldritter-project-db')
            );
        }

        if (isset($data['errors']) && !empty($data['errors'])) {
            $errors = array_map(
                fn($error) => $error['message'] ?? 'Unknown error',
                $data['errors']
            );
            throw new GraphQLException(
                implode(', ', $errors),
                0,
                null,
                'GRAPHQL_ERROR',
                $data['errors']
            );
        }

        return $data['data'] ?? [];
    }

    /**
     * Fetch all projects
     *
     * @param array $options Filter options (tag_ids, limit, offset)
     */
    public function get_projects(array $options = []): array {
        $query = <<<'GRAPHQL'
        query GetProjects($tagIds: [ID!], $limit: Int, $offset: Int) {
            projects(tagIds: $tagIds, limit: $limit, offset: $offset) {
                id
                title
                description
                homepage
                imageUrl
                tags {
                    id
                    title
                    category {
                        id
                        title
                    }
                }
                nextOccurrence {
                    id
                    startDate
                    endDate
                }
                occurrences {
                    id
                    startDate
                    endDate
                }
            }
        }
        GRAPHQL;

        $variables = [];
        if (isset($options['tag_ids'])) {
            $variables['tagIds'] = $options['tag_ids'];
        }
        if (isset($options['limit'])) {
            $variables['limit'] = $options['limit'];
        }
        if (isset($options['offset'])) {
            $variables['offset'] = $options['offset'];
        }

        $result = $this->query($query, $variables);
        return $result['projects'] ?? [];
    }

    /**
     * Fetch a single project by ID
     *
     * @param string|int $id Project ID
     */
    public function get_project($id): ?array {
        $query = <<<'GRAPHQL'
        query GetProject($id: ID!) {
            project(id: $id) {
                id
                title
                description
                homepage
                imageUrl
                tags {
                    id
                    title
                    category {
                        id
                        title
                    }
                }
                nextOccurrence {
                    id
                    startDate
                    endDate
                }
                occurrences {
                    id
                    startDate
                    endDate
                }
            }
        }
        GRAPHQL;

        $result = $this->query($query, ['id' => (string) $id]);
        return $result['project'] ?? null;
    }

    /**
     * Fetch all categories with their tags
     */
    public function get_categories(): array {
        $query = <<<'GRAPHQL'
        query GetCategories {
            categories {
                id
                title
                tags {
                    id
                    title
                }
            }
        }
        GRAPHQL;

        $result = $this->query($query);
        return $result['categories'] ?? [];
    }

    /**
     * Fetch all tags
     */
    public function get_tags(): array {
        $query = <<<'GRAPHQL'
        query GetTags {
            tags {
                id
                title
            }
        }
        GRAPHQL;

        $result = $this->query($query);
        return $result['tags'] ?? [];
    }

    /**
     * Search projects by date range and tags
     *
     * @param array $options Search options (start_date, end_date, tags)
     */
    public function search(array $options): array {
        $query = <<<'GRAPHQL'
        query Search($startDate: ISO8601DateTime!, $endDate: ISO8601DateTime!, $tags: [ID!]) {
            search(startDate: $startDate, endDate: $endDate, tags: $tags) {
                projects {
                    id
                    title
                    description
                    homepage
                    imageUrl
                    tags {
                        id
                        title
                        category {
                            id
                            title
                        }
                    }
                }
                events {
                    id
                    title
                    description
                }
                occurrences {
                    id
                    startDate
                    endDate
                }
            }
        }
        GRAPHQL;

        $variables = [
            'startDate' => $options['start_date'] ?? '1900-01-01T00:00:00Z',
            'endDate' => $options['end_date'] ?? '2100-12-31T23:59:59Z',
        ];

        if (isset($options['tags']) && !empty($options['tags'])) {
            $variables['tags'] = $options['tags'];
        }

        $result = $this->query($query, $variables);
        return $result['search'] ?? ['projects' => [], 'events' => [], 'occurrences' => []];
    }

    /**
     * Get projects filtered by tag titles
     *
     * @param array $tag_titles Array of tag titles to filter by
     * @param array $hidden_categories Category titles to exclude
     */
    public function get_projects_by_tags(array $tag_titles = [], array $hidden_categories = []): array {
        // Use the projects query with tag titles filter
        $query = <<<'GRAPHQL'
        query GetProjectsByTags($tags: [String!]) {
            projects(tags: $tags) {
                id
                title
                description
                homepage
                imageUrl
                tags {
                    id
                    title
                    category {
                        id
                        title
                    }
                }
                nextOccurrence {
                    id
                    startDate
                    endDate
                }
                occurrences {
                    id
                    startDate
                    endDate
                }
            }
        }
        GRAPHQL;

        $variables = [];
        if (!empty($tag_titles)) {
            $variables['tags'] = $tag_titles;
        }

        $result = $this->query($query, $variables);
        $projects = $result['projects'] ?? [];

        // Get categories for filtering
        $categories = $this->get_categories();

        // Build hidden category ID list
        $hidden_category_ids = [];
        foreach ($categories as $category) {
            if (in_array($category['title'], $hidden_categories, true)) {
                $hidden_category_ids[] = $category['id'];
            }
        }

        // Filter out projects in hidden categories
        if (!empty($hidden_category_ids)) {
            $projects = array_filter($projects, function ($project) use ($hidden_category_ids) {
                foreach ($project['tags'] ?? [] as $tag) {
                    if (isset($tag['category']['id']) && in_array($tag['category']['id'], $hidden_category_ids, true)) {
                        return false;
                    }
                }
                return true;
            });
            $projects = array_values($projects);
        }

        // Sort projects by next occurrence date (projects with upcoming dates first)
        usort($projects, function ($a, $b) {
            $a_next = $a['nextOccurrence']['startDate'] ?? null;
            $b_next = $b['nextOccurrence']['startDate'] ?? null;

            // Projects with upcoming occurrences come first
            if ($a_next && !$b_next) {
                return -1;
            }
            if (!$a_next && $b_next) {
                return 1;
            }
            if ($a_next && $b_next) {
                return strcmp($a_next, $b_next);
            }

            // Both have no occurrences - sort by title
            return strcmp($a['title'] ?? '', $b['title'] ?? '');
        });

        return $projects;
    }

    /**
     * Set request timeout
     *
     * @param int $seconds Timeout in seconds
     */
    public function set_timeout(int $seconds): void {
        $this->timeout = max(1, $seconds);
    }
}
