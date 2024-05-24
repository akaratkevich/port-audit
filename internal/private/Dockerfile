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

