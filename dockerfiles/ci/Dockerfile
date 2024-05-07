FROM registry.gitlab.com/gitlab-org/terraform-images/stable:latest

ENV PATH="${PATH}:/usr/local/go/bin"

RUN apk add --update --no-cache curl tar bash

ARG GO_VERSION
RUN curl -sL "https://golang.org/dl/${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    go version

# Clone and prepare the Go project
RUN git clone https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git /tmp/grit
WORKDIR /tmp/grit
RUN go mod download

# Clean up the image
RUN rm -rf /tmp/*
