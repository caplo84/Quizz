# Use ARG to allow dynamic Go version (defaults to 1.23)
ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-bullseye

# Install common tools
RUN apt-get update && apt-get install -y \
    curl \
    bash \
    netcat \
    postgresql-client \
    zsh \
    lsof \
 && rm -rf /var/lib/apt/lists/*

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /usr/local/bin/wait-for-it.sh
RUN chmod +x /usr/local/bin/wait-for-it.sh

WORKDIR /workspace