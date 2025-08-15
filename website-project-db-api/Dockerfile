# Build stage
FROM ruby:3.3.2-alpine AS builder

# Install build dependencies including gcompat for platform issues
RUN apk add --no-cache \
    build-base \
    sqlite-dev \
    nodejs \
    yarn \
    git \
    tzdata \
    imagemagick \
    vips-dev \
    gcompat

# Set working directory
WORKDIR /app

# Copy Gemfile and Gemfile.lock
COPY Gemfile Gemfile.lock ./

# Install gems with platform-specific settings
# Remove platform-specific gems and force compilation
RUN bundle lock --add-platform ruby && \
    bundle config set --local deployment 'false' && \
    bundle config set --local without 'development test' && \
    bundle config set force_ruby_platform true && \
    rm -rf vendor/bundle && \
    bundle install --jobs 4 --retry 3

# Copy application code
COPY . .

# Precompile assets (if any) and prepare database
RUN SECRET_KEY_BASE=dummy bundle exec rails assets:precompile 2>/dev/null || true

# Production stage
FROM ruby:3.3.2-alpine

# Install runtime dependencies including gcompat for nokogiri
RUN apk add --no-cache \
    sqlite-libs \
    tzdata \
    imagemagick \
    vips \
    curl \
    gcompat \
    bash

# Create non-root user
RUN addgroup -g 1000 -S rails && \
    adduser -u 1000 -S rails -G rails

# Set working directory
WORKDIR /app

# Copy application from builder with proper ownership
COPY --from=builder --chown=rails:rails /app /app
COPY --from=builder --chown=rails:rails /usr/local/bundle /usr/local/bundle

# Create necessary directories with proper permissions
RUN mkdir -p /app/tmp/pids /app/tmp/cache /app/log /app/storage /app/db && \
    chown -R rails:rails /app && \
    chmod -R 755 /app

# Copy entrypoint script with proper permissions
COPY --chown=rails:rails docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

# Switch to non-root user
USER rails

# Set Rails environment
ENV RAILS_ENV=production \
    RAILS_LOG_TO_STDOUT=true \
    RAILS_SERVE_STATIC_FILES=true

# Expose port
EXPOSE 3000

# Entrypoint
ENTRYPOINT ["/app/docker-entrypoint.sh"]

# Start server
CMD ["bundle", "exec", "puma", "-C", "config/puma.rb"]