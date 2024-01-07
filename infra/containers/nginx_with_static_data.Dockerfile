FROM imbios/bun-node as builder
ENV NODE_ENV=production

WORKDIR /app

COPY ./frontend/app/package.json ./
RUN bun install

COPY ./frontend/app/ ./
RUN bun run build

FROM nginx:1.25.0
COPY --from=builder /app/build /static
COPY ./infra/nginx/nginx.conf /etc/nginx/nginx.conf
