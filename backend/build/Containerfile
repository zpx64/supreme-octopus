FROM golang:1.20.4 as builder
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go build ./cmd/tzproj

FROM fedora:38
RUN update-ca-trust
WORKDIR /app
COPY ./db/migrations /app/db/migrations
COPY --from=builder /app/tzproj /app/app
ENTRYPOINT ["./app"]
