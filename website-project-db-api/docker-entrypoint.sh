#!/bin/bash
set -e

# Remove any existing server.pid
rm -f /app/tmp/pids/server.pid

# Create database if it doesn't exist
if [ ! -f /app/db/production.sqlite3 ]; then
    echo "Creating database..."
    bundle exec rails db:create
fi

# Run migrations
echo "Running database migrations..."
bundle exec rails db:migrate

# Seed database if requested
if [ "$SEED_DATABASE" = "true" ]; then
    echo "Seeding database..."
    bundle exec rails db:seed
fi

# Execute the main command
exec "$@"