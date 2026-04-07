# Build and run jolt-go example
# Match the Ubuntu 24.04 base used to build the libraries (GCC 13)
FROM ubuntu:24.04

# Install Go and build dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.25.3
RUN wget -q https://go.dev/dl/go1.25.3.linux-arm64.tar.gz \
    && tar -C /usr/local -xzf go1.25.3.linux-arm64.tar.gz \
    && rm go1.25.3.linux-arm64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Set working directory
WORKDIR /app

# Copy go module files
COPY go.mod go.sum* ./

# Download dependencies (if any)
RUN go mod download

# Copy entire project including pre-built binaries
COPY . .

# Enable CGO (required for jolt-go)
ENV CGO_ENABLED=1

# Build the example from the app directory
RUN go build -o example/jolt-example ./example

# Run the example
CMD ["./example/jolt-example"]
