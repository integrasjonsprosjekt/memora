FROM golang:1.24-alpine AS base

RUN apk add --no-cache git ca-certificates tzdata

RUN adduser -D -g '' app

FROM base AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags='-w -s -extldflags "-static"' \
  -a \
  -installsuffix cgo \
  -o memora \
  ./main.go

FROM base

COPY --from=builder /build/memora /memora

USER app

ENV PORT=8080
ENV APP_HOST=0.0.0.0
ENV LOG_LEVEL=warn
ENV GO_ENV=production

EXPOSE $PORT

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://0.0.0.0:${PORT}/api/v1/status/

ENTRYPOINT ["/memora"]
