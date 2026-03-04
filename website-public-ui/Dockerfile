# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build arguments for API URLs
ARG API_URL=http://localhost:3000
ARG PUBLIC_API_URL=http://localhost:3000

ENV API_URL=$API_URL
ENV PUBLIC_API_URL=$PUBLIC_API_URL

# Build the application
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy built output from builder
COPY --from=builder --chown=node:node /app/dist /app/dist
COPY --from=builder --chown=node:node /app/node_modules /app/node_modules
COPY --from=builder --chown=node:node /app/package.json /app/package.json

USER node

# Expose port
EXPOSE 4321

ENV HOST=0.0.0.0
ENV PORT=4321

# Start the Astro SSR server
CMD ["node", "./dist/server/entry.mjs"]
