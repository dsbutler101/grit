FROM registry.gitlab.com/gitlab-org/terraform-images/stable:latest

ENV PATH="${PATH}:/usr/local/go/bin"

RUN apk add --update --no-cache curl tar bash make

ARG GO_VERSION
RUN curl -sL "https://golang.org/dl/${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    go version

RUN go install github.com/magefile/mage@latest && \
    mv $HOME/go/bin/mage /usr/local/go/bin

# Clean up the image
RUN rm -rf /tmp/* ${HOME}/go
