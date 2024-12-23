# Use Node.js LTS version
FROM node:20-alpine

# Install OpenSSL and Prisma dependencies
RUN apk add --no-cache openssl libc6-compat libstdc++ libgcc libssl3


WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies including TypeScript and ts-node
RUN npm ci && \
    npm install typescript ts-node @types/node @types/redis --save-dev

# Copy source code and config
COPY . .
COPY .env.example .env

# Expose port
EXPOSE 8080

# Start application with ts-node
CMD ["npx", "ts-node", "main.ts"]
