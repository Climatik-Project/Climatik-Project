# Build the Go manager binary and webhook server
FROM golang:1.20 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the Go source
COPY cmd/controller/main.go cmd/controller/main.go
COPY api/ api/
COPY internal/ internal/

# Build the manager binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/controller/main.go

# Use a multi-stage build to include the Python environment
FROM python:3.11-slim AS python-builder

WORKDIR /app

# Copy and install Python requirements
COPY python/climatik_operator/requirements.txt /app/
RUN pip install --no-cache-dir -r requirements.txt

# Copy the Python application code
COPY python/climatik_operator /app/

# Final stage: Use alpine base image to package the Go binaries and Python environment
FROM alpine:latest
WORKDIR /

# Install dependencies
RUN apk add --no-cache bash

# Copy the built Go binaries
COPY --from=builder /workspace/manager .

# Copy the installed Python environment
COPY --from=python-builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=python-builder /usr/local/bin /usr/local/bin

# Copy the entire Python application code from the python-builder stage
COPY --from=python-builder /app /app

# Set the user
USER 65532:65532

# Command to run the Go manager binary, the Python operator, and the Slack webhook server
CMD ["/bin/sh", "-c", "/manager --metrics-bind-address=127.0.0.1:8080 --leader-elect"]