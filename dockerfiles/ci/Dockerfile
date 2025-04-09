ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine

RUN go install github.com/magefile/mage@latest

RUN apk add --update --no-cache curl tar bash make git gojq

# terraform 1.5.6 is the last version to be MPLv2 licensed
ENV TF_ZIP=terraform_1.5.6_linux_amd64.zip
RUN curl https://releases.hashicorp.com/terraform/1.5.6/${TF_ZIP} -o ${TF_ZIP} && \
    unzip ${TF_ZIP} -d /usr/local/bin && \
    rm ${TF_ZIP}
