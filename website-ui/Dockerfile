# Build stage
FROM node:20-alpine AS builder

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Set working directory
WORKDIR /app

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source code
COPY . .

# Build arguments for environment variables
ARG VITE_GRAPHQL_URI=http://localhost:1337/graphql
ARG VITE_RAILS_API_URL=http://localhost:3000
ARG VITE_APP_ENV=production

# Set environment variables for build
ENV VITE_GRAPHQL_URI=$VITE_GRAPHQL_URI
ENV VITE_RAILS_API_URL=$VITE_RAILS_API_URL
ENV VITE_APP_ENV=$VITE_APP_ENV

# Build the application
RUN pnpm build

# Production stage
FROM nginx:alpine

# Install gettext for envsubst
RUN apk add --no-cache gettext bash

# Create app user to match Kubernetes security context
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# Create writable directories for nginx
RUN mkdir -p /app/html /app/nginx/cache /app/nginx/log /app/nginx/conf.d /app/nginx/run && \
    chown -R appuser:appuser /app

# Copy built files from builder to writable location
COPY --from=builder --chown=appuser:appuser /app/dist /app/html

# Copy nginx configuration template
COPY --chown=appuser:appuser nginx.conf.template /app/nginx/templates/default.conf.template

# Copy entrypoint script
COPY --chown=appuser:appuser docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

# Create custom nginx.conf that works with non-root user
RUN echo 'worker_processes auto;' > /app/nginx/nginx.conf && \
    echo 'error_log /app/nginx/log/error.log notice;' >> /app/nginx/nginx.conf && \
    echo 'pid /app/nginx/run/nginx.pid;' >> /app/nginx/nginx.conf && \
    echo '' >> /app/nginx/nginx.conf && \
    echo 'events {' >> /app/nginx/nginx.conf && \
    echo '    worker_connections 1024;' >> /app/nginx/nginx.conf && \
    echo '}' >> /app/nginx/nginx.conf && \
    echo '' >> /app/nginx/nginx.conf && \
    echo 'http {' >> /app/nginx/nginx.conf && \
    echo '    include /etc/nginx/mime.types;' >> /app/nginx/nginx.conf && \
    echo '    default_type application/octet-stream;' >> /app/nginx/nginx.conf && \
    echo '' >> /app/nginx/nginx.conf && \
    echo '    log_format main' "'" '$remote_addr - $remote_user [$time_local] "$request" ' >> /app/nginx/nginx.conf && \
    echo '                    ' "'" '$status $body_bytes_sent "$http_referer" ' >> /app/nginx/nginx.conf && \
    echo '                    ' "'" '"$http_user_agent" "$http_x_forwarded_for"' "'" ';' >> /app/nginx/nginx.conf && \
    echo '' >> /app/nginx/nginx.conf && \
    echo '    access_log /app/nginx/log/access.log main;' >> /app/nginx/nginx.conf && \
    echo '    sendfile on;' >> /app/nginx/nginx.conf && \
    echo '    keepalive_timeout 65;' >> /app/nginx/nginx.conf && \
    echo '' >> /app/nginx/nginx.conf && \
    echo '    # Temp paths' >> /app/nginx/nginx.conf && \
    echo '    client_body_temp_path /app/nginx/cache/client_temp;' >> /app/nginx/nginx.conf && \
    echo '    proxy_temp_path /app/nginx/cache/proxy_temp;' >> /app/nginx/nginx.conf && \
    echo '    fastcgi_temp_path /app/nginx/cache/fastcgi_temp;' >> /app/nginx/nginx.conf && \
    echo '    uwsgi_temp_path /app/nginx/cache/uwsgi_temp;' >> /app/nginx/nginx.conf && \
    echo '    scgi_temp_path /app/nginx/cache/scgi_temp;' >> /app/nginx/nginx.conf && \
    echo '' >> /app/nginx/nginx.conf && \
    echo '    include /app/nginx/conf.d/*.conf;' >> /app/nginx/nginx.conf && \
    echo '}' >> /app/nginx/nginx.conf && \
    chown appuser:appuser /app/nginx/nginx.conf

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Use custom entrypoint
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["nginx", "-c", "/app/nginx/nginx.conf", "-g", "daemon off;"]