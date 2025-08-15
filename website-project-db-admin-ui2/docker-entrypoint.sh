#!/bin/bash
set -e

# Function to replace environment variables in JavaScript files
replace_env_vars() {
    # Find all JS and CSS files in the writable app directory
    find /app/html -type f \( -name "*.js" -o -name "*.css" -o -name "*.html" \) | while read file; do
        # Create temp file in /tmp which is writable
        tmpfile="/tmp/$(basename "$file").tmp"
        
        # Replace placeholder values with actual environment variables
        if [ -n "$VITE_API_BASE_URL" ]; then
            sed "s|http://localhost:3000|${VITE_API_BASE_URL}|g" "$file" > "$tmpfile" && mv "$tmpfile" "$file"
        fi
        if [ -n "$VITE_GOOGLE_CLIENT_ID" ]; then
            # Replace empty or placeholder Google Client ID
            sed "s|GOOGLE_CLIENT_ID_PLACEHOLDER|${VITE_GOOGLE_CLIENT_ID}|g" "$file" > "$tmpfile" && mv "$tmpfile" "$file"
        fi
    done
}

# Replace environment variables in built files
replace_env_vars

# Process nginx configuration template if it exists
if [ -f /app/nginx/templates/default.conf.template ]; then
    NGINX_PORT=${NGINX_PORT:-8080}
    envsubst '${NGINX_PORT}' < /app/nginx/templates/default.conf.template > /app/nginx/conf.d/default.conf
fi

# Execute the original command
exec "$@"