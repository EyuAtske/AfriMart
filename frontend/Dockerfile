# ---- Stage 1: Build ----
FROM node:22-alpine AS build

RUN corepack enable && corepack prepare pnpm@9 --activate

WORKDIR /app

# Install dependencies first (cached layer)
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Copy source and build
COPY . .

ARG API_PROXY_TARGET=http://localhost:8080
ENV API_PROXY_TARGET=${API_PROXY_TARGET}

RUN pnpm run build --ignore-workspace

# ---- Stage 2: Run ----
FROM node:22-alpine AS runtime

WORKDIR /app

# Copy only the built output
COPY --from=build /app/.output .output

ENV HOST=0.0.0.0
ENV PORT=3000

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
