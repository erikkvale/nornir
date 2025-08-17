# Builder stage
# ---- Builder ----
FROM cgr.dev/chainguard/go:latest-dev AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o worker ./worker-service

# ---- Runtime ----
FROM cgr.dev/chainguard/static
COPY --from=builder /app/worker /worker

EXPOSE 50051
ENTRYPOINT ["/worker"]
