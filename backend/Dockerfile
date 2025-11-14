ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-bullseye AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/api

FROM debian:bullseye-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/app /app/app
COPY --from=builder /src/configs /app/configs

ENV APP_ENV=production
ENV PORT=8080
EXPOSE 8080

CMD ["/app/app"]