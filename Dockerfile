# Build the Go manager binary
FROM golang:1.20 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the Go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/ internal/

# Build the manager binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go

# Use a multi-stage build to include the Python environment
FROM python:3.11-slim AS python-builder

WORKDIR /app

# Copy and install Python requirements
COPY python/climatik_operator/requirements.txt /app/
RUN pip install --no-cache-dir -r requirements.txt

# Copy the Python application code
COPY python/climatik_operator /app/

# Final stage: Use distroless as minimal base image to package the Go binary and Python environment
FROM gcr.io/distroless/static:nonroot
WORKDIR /

# Copy the built Go binary
COPY --from=builder /workspace/manager .

# Copy the installed Python environment
COPY --from=python-builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=python-builder /usr/local/bin /usr/local/bin

# Set the user
USER 65532:65532

# Command to run the Go manager binary
ENTRYPOINT ["/manager"]