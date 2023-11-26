FROM node:21.2-bookworm as builder
ENV NODE_ENV=production

WORKDIR /app

COPY ./frontend/app/package.json ./
RUN npm install

COPY ./frontend/app/ ./
RUN npx react-scripts build

FROM nginx:1.25.0
COPY --from=builder /app/build /static
COPY ./infra/nginx/nginx.conf /etc/nginx/nginx.conf
