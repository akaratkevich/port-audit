FROM ubuntu:latest

# Set environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# Install necessary packages
RUN apt-get update && apt-get install -y \
    wget \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Download and install Go
RUN wget https://golang.org/dl/go1.21.3.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz \
    && rm go1.21.3.linux-amd64.tar.gz

# Create GOPATH directories
RUN mkdir -p $GOPATH/src $GOPATH/bin $GOPATH/pkg

# Set up a working directory
WORKDIR /workspace

# Clone the repository
#RUN git clone https://github.com/akaratkevich/port-audit.git

# Change working directory to the cloned repository
WORKDIR /workspace/port-audit/cmd

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags=-static" -o /workspace/port-audit/bin/linuxX64/port-audit

# Verify the binary is built
RUN ls -l /workspace/port-audit/bin/linuxX64

# Expose any ports if needed
# EXPOSE 8080

#########
# docker build -t port-audit-builder .
# docker run -it port-audit-builder /bin/bash


