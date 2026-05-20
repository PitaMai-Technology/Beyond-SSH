# syntax=docker/dockerfile:1

# ---------- Build Stage ----------
FROM node:20-alpine AS builder
WORKDIR /app

# Install bun (if not already in node) - use curl script
RUN curl -fsSL https://bun.sh/install | bash && \
    export BUN_INSTALL=\"$HOME/.bun\" && \ 
    export PATH=\"$BUN_INSTALL/bin:$PATH\"

# Copy package files and install dependencies
COPY package.json bun.lock ./
RUN /root/.bun/bin/bun install --frozen-lockfile

# Copy source files
COPY . .

# Build the Vite app (produces dist folder)
RUN /root/.bun/bin/bun run build

# ---------- Runtime Stage ----------
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# Expose default HTTP port
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
