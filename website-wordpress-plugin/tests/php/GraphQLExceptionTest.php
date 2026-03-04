<?php
/**
 * GraphQL Exception Tests
 *
 * @package Waldritter\ProjectDB\Tests
 */

declare(strict_types=1);

namespace Waldritter\ProjectDB\Tests;

use PHPUnit\Framework\TestCase;
use Waldritter\ProjectDB\GraphQLException;

class GraphQLExceptionTest extends TestCase {

    public function testExceptionMessage(): void {
        $exception = new GraphQLException('Test error message');

        $this->assertEquals('Test error message', $exception->getMessage());
    }

    public function testExceptionCode(): void {
        $exception = new GraphQLException('Test error', 500);

        $this->assertEquals(500, $exception->getCode());
    }

    public function testErrorCode(): void {
        $exception = new GraphQLException(
            'Test error',
            0,
            null,
            'CUSTOM_ERROR_CODE'
        );

        $this->assertEquals('CUSTOM_ERROR_CODE', $exception->get_error_code());
    }

    public function testGraphQLErrors(): void {
        $graphqlErrors = [
            ['message' => 'Error 1', 'path' => ['query', 'field']],
            ['message' => 'Error 2', 'path' => ['query', 'otherField']],
        ];

        $exception = new GraphQLException(
            'GraphQL errors occurred',
            0,
            null,
            'GRAPHQL_ERROR',
            $graphqlErrors
        );

        $this->assertEquals($graphqlErrors, $exception->get_graphql_errors());
    }

    public function testIsNetworkErrorForHttpRequestFailed(): void {
        $exception = new GraphQLException(
            'Connection failed',
            0,
            null,
            'http_request_failed'
        );

        $this->assertTrue($exception->is_network_error());
    }

    public function testIsNetworkErrorForNetworkError(): void {
        $exception = new GraphQLException(
            'Network error',
            0,
            null,
            'NETWORK_ERROR'
        );

        $this->assertTrue($exception->is_network_error());
    }

    public function testIsNotNetworkErrorForGraphQLError(): void {
        $exception = new GraphQLException(
            'GraphQL error',
            0,
            null,
            'GRAPHQL_ERROR'
        );

        $this->assertFalse($exception->is_network_error());
    }

    public function testIsGraphQLError(): void {
        $exception = new GraphQLException(
            'GraphQL error',
            0,
            null,
            'GRAPHQL_ERROR'
        );

        $this->assertTrue($exception->is_graphql_error());
    }

    public function testIsNotGraphQLErrorForNetworkError(): void {
        $exception = new GraphQLException(
            'Network error',
            0,
            null,
            'http_request_failed'
        );

        $this->assertFalse($exception->is_graphql_error());
    }

    public function testGetUserMessageForNetworkError(): void {
        $exception = new GraphQLException(
            'Connection failed',
            0,
            null,
            'http_request_failed'
        );

        $message = $exception->get_user_message();

        $this->assertStringContainsString('Unable to connect', $message);
    }

    public function testGetUserMessageForGraphQLError(): void {
        $exception = new GraphQLException(
            'GraphQL error',
            0,
            null,
            'GRAPHQL_ERROR'
        );

        $message = $exception->get_user_message();

        $this->assertStringContainsString('error loading', $message);
    }

    public function testGetUserMessageForUnknownError(): void {
        $exception = new GraphQLException('Unknown error');

        $message = $exception->get_user_message();

        $this->assertStringContainsString('unexpected error', $message);
    }
}
