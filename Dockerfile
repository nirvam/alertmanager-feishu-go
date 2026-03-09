# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o alertmanager-feishu-go main.go

# Final stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /app/alertmanager-feishu-go /alertmanager-feishu-go

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/alertmanager-feishu-go", "serve"]
