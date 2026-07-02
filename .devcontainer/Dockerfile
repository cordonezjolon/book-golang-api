# Use an official Go 1.26 image as the base for development

FROM mcr.microsoft.com/devcontainers/go:1.26

# Install common utilities and PostgreSQL client
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        postgresql-client \
        git \
        curl \
        unzip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Ensure workspace directory exists
RUN mkdir -p /workspace
WORKDIR /workspace
