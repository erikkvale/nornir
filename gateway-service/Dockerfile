# ---- Builder ----
FROM cgr.dev/chainguard/go:latest-dev AS builder
WORKDIR /app

# Copy go mod/sum first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 go build -o gateway ./gateway-service

# ---- Runtime ----
FROM cgr.dev/chainguard/static
COPY --from=builder /app/gateway /gateway

EXPOSE 8080
ENTRYPOINT ["/gateway"]
