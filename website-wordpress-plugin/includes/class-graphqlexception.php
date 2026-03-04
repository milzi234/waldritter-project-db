<?php
/**
 * GraphQL Exception
 *
 * @package Waldritter\ProjectDB
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB;

use Exception;
use Throwable;

/**
 * Exception thrown for GraphQL API errors
 */
class GraphQLException extends Exception {
    /**
     * @var string|null Error code from the API
     */
    private ?string $error_code;

    /**
     * @var array GraphQL errors from the response
     */
    private array $graphql_errors;

    /**
     * Constructor
     *
     * @param string $message Error message
     * @param int $code HTTP status code
     * @param Throwable|null $previous Previous exception
     * @param string|null $error_code API error code
     * @param array $graphql_errors GraphQL errors from response
     */
    public function __construct(
        string $message = '',
        int $code = 0,
        ?Throwable $previous = null,
        ?string $error_code = null,
        array $graphql_errors = []
    ) {
        parent::__construct($message, $code, $previous);
        $this->error_code = $error_code;
        $this->graphql_errors = $graphql_errors;
    }

    /**
     * Get the API error code
     */
    public function get_error_code(): ?string {
        return $this->error_code;
    }

    /**
     * Get the GraphQL errors from the response
     */
    public function get_graphql_errors(): array {
        return $this->graphql_errors;
    }

    /**
     * Check if this is a network error
     */
    public function is_network_error(): bool {
        return $this->error_code === 'http_request_failed'
            || $this->error_code === 'NETWORK_ERROR';
    }

    /**
     * Check if this is a GraphQL-level error
     */
    public function is_graphql_error(): bool {
        return $this->error_code === 'GRAPHQL_ERROR';
    }

    /**
     * Get a user-friendly error message
     */
    public function get_user_message(): string {
        if ($this->is_network_error()) {
            return __('Unable to connect to the project database. Please try again later.', 'waldritter-project-db');
        }

        if ($this->is_graphql_error()) {
            return __('There was an error loading the project data.', 'waldritter-project-db');
        }

        return __('An unexpected error occurred.', 'waldritter-project-db');
    }
}
