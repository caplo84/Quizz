ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}-bullseye AS dev
WORKDIR /workspace

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    curl \
    netcat-openbsd \
    postgresql-client \
 && rm -rf /var/lib/apt/lists/*

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080
CMD ["go", "run", "./cmd/api/"]

FROM golang:${GO_VERSION}-bullseye AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -o /out/app ./cmd/api

FROM debian:bullseye-slim AS runtime
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/app /app/app
COPY --from=builder /src/configs /app/configs

ENV APP_ENV=production
ENV PORT=8080
EXPOSE 8080

CMD ["/app/app"]