FROM golang:1.18 as builder
# Add a work directory
WORKDIR /workspace
# Cache and install dependencies
COPY go.mod go.sum ./
# Copy source files
COPY . .
# Build app
RUN go build -o go-rod-example

FROM debian:buster-slim
# Install chromium
RUN apt-get update && \
    apt-get install chromium -y

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /workspace/go-rod-example ./go-rod-example

RUN adduser app
RUN chown -R app:app /app

ENTRYPOINT ["/app/go-rod-example"]