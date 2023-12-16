FROM golang:1.21.4 as builder
RUN apt update
RUN apt install -y libwebp-dev libwebp7

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build ./cmd/supreme-octopus/


FROM ubuntu:23.04
RUN apt update
RUN apt install -y ca-certificates
RUN apt install -y libwebp7
RUN update-ca-certificates
WORKDIR /app
COPY --from=builder /app/supreme-octopus /app/app
ENTRYPOINT ["./app"]
