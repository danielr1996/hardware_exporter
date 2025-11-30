# --- Build stage ---
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build only the main application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /hardware_exporter ./cmd/hardware_exporter

# --- Runtime stage ---
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /hardware_exporter /hardware_exporter
EXPOSE 9105

ENTRYPOINT ["/hardware_exporter"]
