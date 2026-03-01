# Stage 1: Install dependencies
FROM node:20-alpine AS deps

WORKDIR /app

# Copy dependency files first for caching
COPY package.json yarn.lock ./

# Install only production dependencies
RUN yarn install --frozen-lockfile --production

# Stage 2: Final image
FROM node:20-alpine

WORKDIR /app

# Copy node_modules from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy the rest of the app
COPY . .

# Expose port
EXPOSE 8080

# Start the app from src folder
CMD ["node", "src/index.js"]