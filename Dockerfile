# Use golang alpine image as the builder stage for multi-arch support
FROM --platform=$BUILDPLATFORM golang:1.22.4-alpine3.20 AS builder

# Install git and other necessary tools
RUN apk update && apk add --no-cache git bash

# Set the Current Working Directory inside the container
WORKDIR /src

# Copy everything from the current directory to the working directory inside the container
COPY . .

# Fetch dependencies using go mod if your project uses Go modules
RUN go mod download

# Version and Git Commit build arguments
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE
ARG TARGETOS
ARG TARGETARCH

# Compile the Go application and embed version information
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -installsuffix cgo \
    -ldflags "-X github.com/mattmattox/go-smtp-relay/pkg/version.Version=$VERSION \
    -X github.com/mattmattox/go-smtp-relay/pkg/version.GitCommit=$GIT_COMMIT \
    -X github.com/mattmattox/go-smtp-relay/pkg/version.BuildTime=$BUILD_DATE" \
    -o /bin/go-smtp-relay

# Use minimal Ubuntu image as the final image for multi-arch support
FROM --platform=$TARGETPLATFORM ubuntu:latest

# Install common dependencies
RUN apt-get update && \
    apt-get install -y \
    ca-certificates \
    curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /bin/go-smtp-relay /root/go-smtp-relay

# Ensure the binary is executable
RUN chmod +x /root/go-smtp-relay

EXPOSE 8080

# Command to run the executable
CMD ["./go-smtp-relay"]
