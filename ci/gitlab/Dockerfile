# syntax=docker/dockerfile:1

FROM ghcr.io/zoomoid/assignments/cli:alpine

ADD https://release-cli-downloads.s3.amazonaws.com/latest/release-cli-linux-amd64 /usr/local/bin/release-cli

RUN chmod +x /usr/local/bin/release-cli && release-cli -v

