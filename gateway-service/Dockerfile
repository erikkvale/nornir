FROM cgr.dev/chainguard/go:latest-dev AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o gateway ./gateway-service

FROM cgr.dev/chainguard/static
COPY --from=builder /app/gateway /gateway

EXPOSE 8080
ENTRYPOINT ["/gateway"]
